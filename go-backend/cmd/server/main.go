package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"biliup/internal/api"
	"biliup/internal/config"
	"biliup/internal/database"
	"biliup/internal/service"
)

func main() {
	port := flag.Int("port", 19159, "服务器端口")
	dbHost := flag.String("db-host", "127.0.0.1", "MySQL 主机")
	dbPort := flag.Int("db-port", 3306, "MySQL 端口")
	dbUser := flag.String("db-user", "root", "MySQL 用户名")
	dbPass := flag.String("db-pass", "", "MySQL 密码")
	dbName := flag.String("db-name", "biliup", "MySQL 数据库名")
	flag.Parse()

	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		*dbUser, *dbPass, *dbHost, *dbPort, *dbName,
	)

	// 初始化数据库
	db, err := database.Init(dsn)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 加载全局配置
	cfg, err := config.Load(db)
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 初始化服务
	twitchService := service.NewTwitchService(cfg)
	uploadService := service.NewUploadService(cfg)
	manager := service.NewDownloadManager(db, cfg, twitchService, uploadService)

	// 启动任务调度器
	manager.Start()

	// 启动 HTTP 服务
	router := api.SetupRouter(db, cfg, manager)
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("🚀 Biliup 服务启动在 http://localhost%s", addr)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	<-quit
	log.Println("正在关闭服务...")
	manager.Stop()
	log.Println("服务已关闭")
}
