package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"biliup/internal/config"
	"biliup/internal/model"

	"github.com/jmoiron/sqlx"
)

// DownloadManager 下载管理器
type DownloadManager struct {
	db             *sqlx.DB
	cfg            *config.Config
	twitch         *TwitchService
	upload         *UploadService
	tasks          map[int64]*RecordTask
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	startTime      time.Time
	logBuffer      *LogBuffer
}

// RecordTask 录制任务
type RecordTask struct {
	Streamer  *model.LiveStreamer
	Status    string // idle, checking, recording, uploading, error
	Error     string
	StartTime time.Time
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// LogBuffer 日志缓冲
type LogBuffer struct {
	mu       sync.RWMutex
	messages []LogMessage
	maxSize  int
	subscribers map[string]chan LogMessage
	subMu    sync.Mutex
}

type LogMessage struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

func NewLogBuffer(maxSize int) *LogBuffer {
	return &LogBuffer{
		messages:    make([]LogMessage, 0, maxSize),
		maxSize:     maxSize,
		subscribers: make(map[string]chan LogMessage),
	}
}

func (lb *LogBuffer) Add(level, msg string) {
	lm := LogMessage{
		Time:    time.Now(),
		Level:   level,
		Message: msg,
	}

	lb.mu.Lock()
	if len(lb.messages) >= lb.maxSize {
		lb.messages = lb.messages[1:]
	}
	lb.messages = append(lb.messages, lm)
	lb.mu.Unlock()

	// 通知所有订阅者
	lb.subMu.Lock()
	for _, ch := range lb.subscribers {
		select {
		case ch <- lm:
		default:
		}
	}
	lb.subMu.Unlock()
}

func (lb *LogBuffer) Subscribe(id string) chan LogMessage {
	ch := make(chan LogMessage, 100)
	lb.subMu.Lock()
	lb.subscribers[id] = ch
	lb.subMu.Unlock()
	return ch
}

func (lb *LogBuffer) Unsubscribe(id string) {
	lb.subMu.Lock()
	if ch, ok := lb.subscribers[id]; ok {
		close(ch)
		delete(lb.subscribers, id)
	}
	lb.subMu.Unlock()
}

func (lb *LogBuffer) GetRecent(n int) []LogMessage {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if n > len(lb.messages) {
		n = len(lb.messages)
	}
	result := make([]LogMessage, n)
	copy(result, lb.messages[len(lb.messages)-n:])
	return result
}

func NewDownloadManager(db *sqlx.DB, cfg *config.Config, twitch *TwitchService, upload *UploadService) *DownloadManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &DownloadManager{
		db:        db,
		cfg:       cfg,
		twitch:    twitch,
		upload:    upload,
		tasks:     make(map[int64]*RecordTask),
		ctx:       ctx,
		cancel:    cancel,
		startTime: time.Now(),
		logBuffer: NewLogBuffer(1000),
	}
}

// Start 启动调度循环
func (dm *DownloadManager) Start() {
	dm.logBuffer.Add("info", "下载管理器启动")

	go dm.schedulerLoop()
}

// Stop 停止所有任务
func (dm *DownloadManager) Stop() {
	dm.logBuffer.Add("info", "下载管理器停止中...")
	dm.cancel()

	dm.mu.Lock()
	for _, task := range dm.tasks {
		if task.cancel != nil {
			task.cancel()
		}
	}
	dm.mu.Unlock()
}

// GetLogBuffer 获取日志缓冲
func (dm *DownloadManager) GetLogBuffer() *LogBuffer {
	return dm.logBuffer
}

// schedulerLoop 定期检测直播状态
func (dm *DownloadManager) schedulerLoop() {
	// 首次检测
	dm.checkAllStreamers()

	interval := dm.cfg.GetInt("check_interval", 60)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-dm.ctx.Done():
			return
		case <-ticker.C:
			dm.checkAllStreamers()
		}
	}
}

// checkAllStreamers 检测所有主播
func (dm *DownloadManager) checkAllStreamers() {
	var streamers []model.LiveStreamer
	err := dm.db.Select(&streamers, "SELECT * FROM live_streamers WHERE paused = 0")
	if err != nil {
		dm.logBuffer.Add("error", fmt.Sprintf("查询主播列表失败: %v", err))
		return
	}

	if len(streamers) == 0 {
		return
	}

	// 提取频道名
	channels := make([]string, 0, len(streamers))
	channelMap := make(map[string]*model.LiveStreamer)
	for i := range streamers {
		channel := ExtractChannelName(streamers[i].URL)
		channels = append(channels, channel)
		channelMap[channel] = &streamers[i]
	}

	// 批量检测
	results, err := dm.twitch.BatchCheckStreams(channels)
	if err != nil {
		dm.logBuffer.Add("error", fmt.Sprintf("批量检测失败: %v", err))
		return
	}

	for channel, info := range results {
		streamer := channelMap[channel]
		if streamer == nil {
			continue
		}

		dm.mu.RLock()
		task, exists := dm.tasks[streamer.ID]
		dm.mu.RUnlock()

		if info.IsLive {
			if !exists || task.Status == "idle" || task.Status == "error" {
				dm.logBuffer.Add("info", fmt.Sprintf("检测到 %s 正在直播: %s", channel, info.Title))
				dm.startRecording(streamer, info)
			}
		} else {
			if exists && task.Status == "recording" {
				dm.logBuffer.Add("info", fmt.Sprintf("%s 已下播", channel))
			}
		}
	}
}

// startRecording 开始录制
func (dm *DownloadManager) startRecording(streamer *model.LiveStreamer, info *model.TwitchStreamInfo) {
	ctx, cancel := context.WithCancel(dm.ctx)

	task := &RecordTask{
		Streamer:  streamer,
		Status:    "recording",
		StartTime: time.Now(),
		cancel:    cancel,
	}

	dm.mu.Lock()
	dm.tasks[streamer.ID] = task
	dm.mu.Unlock()

	go dm.recordStream(ctx, task, info)
}

// recordStream 实际录制逻辑
func (dm *DownloadManager) recordStream(ctx context.Context, task *RecordTask, info *model.TwitchStreamInfo) {
	channel := ExtractChannelName(task.Streamer.URL)
	dm.logBuffer.Add("info", fmt.Sprintf("开始录制 %s", channel))

	defer func() {
		task.mu.Lock()
		if task.Status == "recording" {
			task.Status = "idle"
		}
		task.mu.Unlock()
	}()

	// 获取流 URL
	streamURL, err := dm.twitch.GetStreamURL(channel)
	if err != nil {
		dm.logBuffer.Add("error", fmt.Sprintf("获取 %s 流地址失败: %v", channel, err))
		task.mu.Lock()
		task.Status = "error"
		task.Error = err.Error()
		task.mu.Unlock()
		return
	}

	// 准备输出目录
	downloadDir := dm.cfg.Get("download_dir")
	if downloadDir == "" {
		downloadDir = "./downloads"
	}
	os.MkdirAll(downloadDir, 0755)

	// 生成文件名
	now := time.Now()
	prefix := task.Streamer.FilenamePrefix
	if prefix == "" {
		prefix = "{streamer}_%Y-%m-%d_%H_%M_%S"
	}
	filename := strings.ReplaceAll(prefix, "{streamer}", channel)
	filename = now.Format(strings.NewReplacer(
		"%Y", "2006",
		"%m", "01",
		"%d", "02",
		"%H", "15",
		"%M", "04",
		"%S", "05",
	).Replace(filename))

	format := task.Streamer.Format
	if format == "" {
		format = "flv"
	}
	outputPath := filepath.Join(downloadDir, fmt.Sprintf("%s.%s", filename, format))

	// 保存录制信息到数据库
	result, err := dm.db.Exec(
		"INSERT INTO streamer_info (name, url, title, status) VALUES (?, ?, ?, 'recording')",
		channel, task.Streamer.URL, info.Title,
	)
	var streamerInfoID int64
	if err == nil {
		streamerInfoID, _ = result.LastInsertId()
	}

	// 使用 FFmpeg 录制
	segmentTime := task.Streamer.SegmentTime
	if segmentTime == 0 {
		segmentTime = dm.cfg.GetInt("segment_time", 3600)
	}

	args := []string{
		"-y",
		"-headers", fmt.Sprintf("Client-ID: %s\r\n", twitchClientID),
		"-i", streamURL,
		"-c", "copy",
		"-f", format,
	}

	// 分段录制
	if segmentTime > 0 {
		segmentPattern := filepath.Join(downloadDir, fmt.Sprintf("%s_%%03d.%s", filename, format))
		args = append(args,
			"-f", "segment",
			"-segment_time", fmt.Sprintf("%d", segmentTime),
			"-reset_timestamps", "1",
			segmentPattern,
		)
	} else {
		args = append(args, outputPath)
	}

	dm.logBuffer.Add("info", fmt.Sprintf("FFmpeg 命令: ffmpeg %s", strings.Join(args, " ")))

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			dm.logBuffer.Add("info", fmt.Sprintf("录制 %s 被取消", channel))
		} else {
			dm.logBuffer.Add("error", fmt.Sprintf("录制 %s 出错: %v", channel, err))
			task.mu.Lock()
			task.Status = "error"
			task.Error = err.Error()
			task.mu.Unlock()
		}
	}

	dm.logBuffer.Add("info", fmt.Sprintf("录制 %s 完成", channel))

	// 更新数据库
	dm.db.Exec("UPDATE streamer_info SET status = 'completed' WHERE id = ?", streamerInfoID)

	// 查找录制好的文件
	var recordedFiles []string
	if segmentTime > 0 {
		pattern := filepath.Join(downloadDir, fmt.Sprintf("%s_*.%s", filename, format))
		matches, _ := filepath.Glob(pattern)
		recordedFiles = matches
	} else {
		if _, err := os.Stat(outputPath); err == nil {
			recordedFiles = []string{outputPath}
		}
	}

	// 保存文件记录
	for _, f := range recordedFiles {
		fi, _ := os.Stat(f)
		var size int64
		if fi != nil {
			size = fi.Size()
		}
		dm.db.Exec(
			"INSERT INTO file_list (file_path, file_size, streamer_info_id) VALUES (?, ?, ?)",
			f, size, streamerInfoID,
		)
	}

	// 自动上传
	if task.Streamer.UploadTemplateID > 0 && len(recordedFiles) > 0 {
		go dm.uploadFiles(task, recordedFiles, streamerInfoID)
	}
}

// uploadFiles 上传录制的文件
func (dm *DownloadManager) uploadFiles(task *RecordTask, files []string, streamerInfoID int64) {
	task.mu.Lock()
	task.Status = "uploading"
	task.mu.Unlock()

	defer func() {
		task.mu.Lock()
		task.Status = "idle"
		task.mu.Unlock()
	}()

	// 获取上传模板
	var template model.UploadTemplate
	err := dm.db.Get(&template, "SELECT * FROM upload_templates WHERE id = ?", task.Streamer.UploadTemplateID)
	if err != nil {
		dm.logBuffer.Add("error", fmt.Sprintf("获取上传模板失败: %v", err))
		return
	}

	// 获取 cookie
	cookie := template.UserCookie
	if cookie == "" {
		var user model.User
		if err := dm.db.Get(&user, "SELECT * FROM users LIMIT 1"); err != nil {
			dm.logBuffer.Add("error", "未找到B站登录信息")
			return
		}
		cookie = user.CookieData
	}

	for _, filePath := range files {
		dm.logBuffer.Add("info", fmt.Sprintf("开始上传: %s", filepath.Base(filePath)))

		// 过��小文件
		threshold := dm.cfg.GetInt("filtering_threshold", 10)
		fi, err := os.Stat(filePath)
		if err != nil {
			continue
		}
		if fi.Size() < int64(threshold)*1024*1024 {
			dm.logBuffer.Add("info", fmt.Sprintf("文件 %s 太小 (%.1fMB < %dMB)，跳过上传",
				filepath.Base(filePath), float64(fi.Size())/1024/1024, threshold))
			continue
		}

		if err := dm.upload.UploadVideo(filePath, &template, cookie); err != nil {
			dm.logBuffer.Add("error", fmt.Sprintf("上传失败: %v", err))
			continue
		}

		// 标记已上传
		dm.db.Exec("UPDATE file_list SET uploaded = 1 WHERE file_path = ?", filePath)
		dm.logBuffer.Add("info", fmt.Sprintf("上传完成: %s", filepath.Base(filePath)))
	}
}

// ============ 对外接口 ============

// GetStatus 获取系统状态
func (dm *DownloadManager) GetStatus() *model.StatusResponse {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	streamerStatus := make(map[string]string)
	activeTasks := 0
	for _, task := range dm.tasks {
		channel := ExtractChannelName(task.Streamer.URL)
		task.mu.Lock()
		streamerStatus[channel] = task.Status
		if task.Status == "recording" || task.Status == "uploading" {
			activeTasks++
		}
		task.mu.Unlock()
	}

	var totalRecorded, totalUploaded int64
	dm.db.Get(&totalRecorded, "SELECT COUNT(*) FROM file_list")
	dm.db.Get(&totalUploaded, "SELECT COUNT(*) FROM file_list WHERE uploaded = 1")

	uptime := time.Since(dm.startTime)

	return &model.StatusResponse{
		Version:        "2.0.0",
		Uptime:         formatDuration(uptime),
		ActiveTasks:    activeTasks,
		TotalRecorded:  totalRecorded,
		TotalUploaded:  totalUploaded,
		StreamerStatus: streamerStatus,
	}
}

// GetTaskStatus 获取指定主播的任务状态
func (dm *DownloadManager) GetTaskStatus(streamerID int64) string {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if task, ok := dm.tasks[streamerID]; ok {
		task.mu.Lock()
		defer task.mu.Unlock()
		return task.Status
	}
	return "idle"
}

// StopTask 停止指定录制任务
func (dm *DownloadManager) StopTask(streamerID int64) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if task, ok := dm.tasks[streamerID]; ok {
		if task.cancel != nil {
			task.cancel()
		}
		task.mu.Lock()
		task.Status = "idle"
		task.mu.Unlock()
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
