package api

import (
	"net/http"

	"biliup/internal/config"
	"biliup/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// SetupRouter 设置路由
func SetupRouter(db *sqlx.DB, cfg *config.Config, manager *service.DownloadManager) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// CORS 中间件
	r.Use(corsMiddleware())

	// 静态文件（前端）
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/", "./frontend/dist/index.html")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// API 路由
	h := NewHandler(db, cfg, manager)

	v1 := r.Group("/v1")
	{
		// 主播管理
		v1.GET("/streamers", h.GetStreamers)
		v1.POST("/streamers", h.AddStreamer)
		v1.PUT("/streamers", h.UpdateStreamer)
		v1.DELETE("/streamers/:id", h.DeleteStreamer)
		v1.PUT("/streamers/:id/pause", h.PauseStreamer)

		// 全局配置
		v1.GET("/configuration", h.GetConfiguration)
		v1.PUT("/configuration", h.UpdateConfiguration)

		// 上传模板
		v1.GET("/upload/streamers", h.GetUploadTemplates)
		v1.POST("/upload/streamers", h.AddUploadTemplate)
		v1.GET("/upload/streamers/:id", h.GetUploadTemplate)
		v1.DELETE("/upload/streamers/:id", h.DeleteUploadTemplate)

		// 用户管理 (B站账号)
		v1.GET("/users", h.GetUsers)
		v1.POST("/users", h.AddUser)
		v1.DELETE("/users/:id", h.DeleteUser)

		// B站相关
		v1.GET("/get_qrcode", h.GetQRCode)
		v1.POST("/login_by_qrcode", h.LoginByQRCode)

		// 视频文件
		v1.GET("/videos", h.GetVideos)

		// 系统状态
		v1.GET("/status", h.GetStatus)

		// WebSocket 日志
		v1.GET("/logs", h.WebSocketLogs)
	}

	// Bilibili 代理
	bili := r.Group("/bili")
	{
		bili.GET("/archive/pre", h.GetArchivePre)
		bili.GET("/space/myinfo", h.GetMyInfo)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
