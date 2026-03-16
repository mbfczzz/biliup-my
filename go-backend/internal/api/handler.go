package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"biliup/internal/config"
	"biliup/internal/model"
	"biliup/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

// Handler API 处理器
type Handler struct {
	db      *sqlx.DB
	cfg     *config.Config
	manager *service.DownloadManager
	upload  *service.UploadService
}

func NewHandler(db *sqlx.DB, cfg *config.Config, manager *service.DownloadManager) *Handler {
	return &Handler{
		db:      db,
		cfg:     cfg,
		manager: manager,
		upload:  service.NewUploadService(cfg),
	}
}

// ============ 主播管理 ============

func (h *Handler) GetStreamers(c *gin.Context) {
	var streamers []model.LiveStreamer
	if err := h.db.Select(&streamers, "SELECT * FROM live_streamers ORDER BY id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 附加运行时状态
	for i := range streamers {
		streamers[i].Status = h.manager.GetTaskStatus(streamers[i].ID)
	}

	c.JSON(http.StatusOK, streamers)
}

func (h *Handler) AddStreamer(c *gin.Context) {
	var input model.LiveStreamer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 验证 Twitch URL
	if _, err := service.ParseTwitchURL(input.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.db.Exec(
		`INSERT INTO live_streamers (url, remark, filename_prefix, time_range, upload_template_id, format, segment_time, file_size, opt_args, excluded_keywords)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.URL, input.Remark, input.FilenamePrefix, input.TimeRange,
		input.UploadTemplateID, input.Format, input.SegmentTime, input.FileSize,
		input.OptArgs, input.ExcludedKeywords,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	input.ID = id
	c.JSON(http.StatusOK, input)
}

func (h *Handler) UpdateStreamer(c *gin.Context) {
	var input model.LiveStreamer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	_, err := h.db.Exec(
		`UPDATE live_streamers SET url=?, remark=?, filename_prefix=?, time_range=?,
		 upload_template_id=?, format=?, segment_time=?, file_size=?,
		 opt_args=?, excluded_keywords=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		input.URL, input.Remark, input.FilenamePrefix, input.TimeRange,
		input.UploadTemplateID, input.Format, input.SegmentTime, input.FileSize,
		input.OptArgs, input.ExcludedKeywords, input.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func (h *Handler) DeleteStreamer(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err == nil {
		h.manager.StopTask(idInt)
	}

	_, err = h.db.Exec("DELETE FROM live_streamers WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *Handler) PauseStreamer(c *gin.Context) {
	id := c.Param("id")

	var current model.LiveStreamer
	if err := h.db.Get(&current, "SELECT * FROM live_streamers WHERE id = ?", id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "主播不存在"})
		return
	}

	newPaused := 1
	if current.Paused == 1 {
		newPaused = 0
	}

	h.db.Exec("UPDATE live_streamers SET paused = ? WHERE id = ?", newPaused, id)

	// 如果暂停，停止录制
	if newPaused == 1 {
		idInt, _ := strconv.ParseInt(id, 10, 64)
		h.manager.StopTask(idInt)
	}

	status := "已恢复"
	if newPaused == 1 {
		status = "已暂停"
	}

	c.JSON(http.StatusOK, gin.H{"message": status, "paused": newPaused})
}

// ============ 全局配置 ============

func (h *Handler) GetConfiguration(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.GetAll())
}

func (h *Handler) UpdateConfiguration(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.cfg.SetAll(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置已保存"})
}

// ============ 上传模板 ============

func (h *Handler) GetUploadTemplates(c *gin.Context) {
	var templates []model.UploadTemplate
	if err := h.db.Select(&templates, "SELECT * FROM upload_templates ORDER BY id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *Handler) AddUploadTemplate(c *gin.Context) {
	var input model.UploadTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	result, err := h.db.Exec(
		`INSERT INTO upload_templates (template_name, title, tid, copyright, copyright_source,
		 cover_path, description, dynamic, dtime, dolby, hires, no_reprint, tags, user_cookie)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		input.TemplateName, input.Title, input.TID, input.Copyright, input.CopyrightSource,
		input.CoverPath, input.Description, input.Dynamic, input.Dtime,
		input.Dolby, input.Hires, input.NoReprint, input.Tags, input.UserCookie,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	input.ID = id
	c.JSON(http.StatusOK, input)
}

func (h *Handler) GetUploadTemplate(c *gin.Context) {
	id := c.Param("id")
	var template model.UploadTemplate
	if err := h.db.Get(&template, "SELECT * FROM upload_templates WHERE id = ?", id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "模板不存在"})
		return
	}
	c.JSON(http.StatusOK, template)
}

func (h *Handler) DeleteUploadTemplate(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("DELETE FROM upload_templates WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ============ 用户管理 ============

func (h *Handler) GetUsers(c *gin.Context) {
	var users []model.User
	if err := h.db.Select(&users, "SELECT id, username, uid, nickname, created_at FROM users ORDER BY id"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) AddUser(c *gin.Context) {
	var input struct {
		Username   string `json:"username"`
		CookieData string `json:"cookie_data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	result, err := h.db.Exec(
		"INSERT INTO users (username, cookie_data) VALUES (?, ?) ON DUPLICATE KEY UPDATE cookie_data = VALUES(cookie_data)",
		input.Username, input.CookieData,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id, "message": "添加成功"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ============ B站登录 ============

func (h *Handler) GetQRCode(c *gin.Context) {
	qr, err := h.upload.GetQRCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qr)
}

func (h *Handler) LoginByQRCode(c *gin.Context) {
	var input struct {
		QRCodeKey string `json:"qrcode_key"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	result, err := h.upload.LoginByQRCode(input.QRCodeKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetMyInfo(c *gin.Context) {
	cookie := c.GetHeader("Cookie")
	if cookie == "" {
		// 从数据库获取第一个用户的 cookie
		var user model.User
		if err := h.db.Get(&user, "SELECT * FROM users LIMIT 1"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未登录"})
			return
		}
		cookie = user.CookieData
	}

	info, err := h.upload.GetMyInfo(cookie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) GetArchivePre(c *gin.Context) {
	var user model.User
	if err := h.db.Get(&user, "SELECT * FROM users LIMIT 1"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未登录"})
		return
	}

	info, err := h.upload.GetArchivePre(user.CookieData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// ============ 视频文件 ============

func (h *Handler) GetVideos(c *gin.Context) {
	downloadDir := h.cfg.Get("download_dir")
	if downloadDir == "" {
		downloadDir = "./downloads"
	}

	var videos []model.VideoFile

	filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".flv" || ext == ".mp4" || ext == ".ts" || ext == ".mkv" {
			videos = append(videos, model.VideoFile{
				Name:    info.Name(),
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}
		return nil
	})

	c.JSON(http.StatusOK, videos)
}

// ============ 系统状态 ============

func (h *Handler) GetStatus(c *gin.Context) {
	status := h.manager.GetStatus()
	c.JSON(http.StatusOK, status)
}

// ============ WebSocket 日志 ============

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) WebSocketLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	logBuffer := h.manager.GetLogBuffer()
	subID := uuid.New().String()
	ch := logBuffer.Subscribe(subID)
	defer logBuffer.Unsubscribe(subID)

	// 发送最近的日志
	recent := logBuffer.GetRecent(50)
	for _, msg := range recent {
		conn.WriteJSON(msg)
	}

	// 实时推送
	for msg := range ch {
		if err := conn.WriteJSON(msg); err != nil {
			break
		}
	}
}

// ============ 辅助函数 ============

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
