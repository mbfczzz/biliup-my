package model

import "time"

// LiveStreamer 直播录制配置
type LiveStreamer struct {
	ID               int64     `db:"id" json:"id"`
	URL              string    `db:"url" json:"url"`
	Remark           string    `db:"remark" json:"remark"`
	FilenamePrefix   string    `db:"filename_prefix" json:"filename_prefix"`
	TimeRange        string    `db:"time_range" json:"time_range"`
	UploadTemplateID int64     `db:"upload_template_id" json:"upload_template_id"`
	Format           string    `db:"format" json:"format"`
	SegmentTime      int       `db:"segment_time" json:"segment_time"`
	FileSize         int64     `db:"file_size" json:"file_size"`
	OptArgs          string    `db:"opt_args" json:"opt_args"`
	ExcludedKeywords string    `db:"excluded_keywords" json:"excluded_keywords"`
	Paused           int       `db:"paused" json:"paused"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
	// 运行时状态（不存数据库）
	Status string `db:"-" json:"status"`
}

// UploadTemplate B站上传模板
type UploadTemplate struct {
	ID              int64     `db:"id" json:"id"`
	TemplateName    string    `db:"template_name" json:"template_name"`
	Title           string    `db:"title" json:"title"`
	TID             int       `db:"tid" json:"tid"`
	Copyright       int       `db:"copyright" json:"copyright"`
	CopyrightSource string    `db:"copyright_source" json:"copyright_source"`
	CoverPath       string    `db:"cover_path" json:"cover_path"`
	Description     string    `db:"description" json:"description"`
	Dynamic         string    `db:"dynamic" json:"dynamic"`
	Dtime           int64     `db:"dtime" json:"dtime"`
	Dolby           int       `db:"dolby" json:"dolby"`
	Hires           int       `db:"hires" json:"hires"`
	NoReprint       int       `db:"no_reprint" json:"no_reprint"`
	Tags            string    `db:"tags" json:"tags"`
	UserCookie      string    `db:"user_cookie" json:"user_cookie"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// StreamerInfo 录制元数据
type StreamerInfo struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	URL           string    `db:"url" json:"url"`
	Title         string    `db:"title" json:"title"`
	LiveCoverPath string    `db:"live_cover_path" json:"live_cover_path"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// FileItem 已录文件
type FileItem struct {
	ID             int64     `db:"id" json:"id"`
	FilePath       string    `db:"file_path" json:"file_path"`
	FileSize       int64     `db:"file_size" json:"file_size"`
	StreamerInfoID int64     `db:"streamer_info_id" json:"streamer_info_id"`
	Uploaded       int       `db:"uploaded" json:"uploaded"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

// User B站用户（cookie）
type User struct {
	ID         int64     `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	CookieData string    `db:"cookie_data" json:"cookie_data"`
	UID        int64     `db:"uid" json:"uid"`
	Nickname   string    `db:"nickname" json:"nickname"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// Configuration 全局配置
type Configuration struct {
	ID    int64  `db:"id" json:"id"`
	Key   string `db:"key" json:"key"`
	Value string `db:"value" json:"value"`
}

// ============ 请求/响应结构 ============

// TwitchStreamInfo Twitch 流信息
type TwitchStreamInfo struct {
	UserLogin   string `json:"user_login"`
	UserName    string `json:"user_name"`
	Title       string `json:"title"`
	ViewerCount int    `json:"viewer_count"`
	GameName    string `json:"game_name"`
	ThumbnailURL string `json:"thumbnail_url"`
	StreamURL   string `json:"stream_url"`
	IsLive      bool   `json:"is_live"`
}

// QRCodeResponse 二维码登录响应
type QRCodeResponse struct {
	URL       string `json:"url"`
	QRCodeKey string `json:"qrcode_key"`
	Image     string `json:"image"` // base64
}

// LoginResponse 登录结果
type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cookie  string `json:"cookie,omitempty"`
}

// StatusResponse 系统状态
type StatusResponse struct {
	Version        string            `json:"version"`
	Uptime         string            `json:"uptime"`
	ActiveTasks    int               `json:"active_tasks"`
	TotalRecorded  int64             `json:"total_recorded"`
	TotalUploaded  int64             `json:"total_uploaded"`
	StreamerStatus map[string]string `json:"streamer_status"`
}

// VideoFile 视频文件信息
type VideoFile struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mod_time"`
}
