package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sm-medical/internal/api"
	"sm-medical/internal/crypto"
	"sm-medical/internal/middleware"
	"sm-medical/pkg/config"
	"sm-medical/pkg/database"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化国密算法
	if err := crypto.InitCrypto(cfg.Crypto.SM4Key); err != nil {
		log.Fatalf("Failed to init crypto: %v", err)
	}
	log.Println("Crypto initialized successfully")

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	
	// 确保上传目录存在
	if err := os.MkdirAll(cfg.Upload.UploadPath, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	log.Printf("Upload directory ready: %s", cfg.Upload.UploadPath)

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.CORS())
	
	// 配置静态文件服务（提供上传的文件访问）
	r.Static("/uploads", cfg.Upload.UploadPath)
	log.Printf("Static files served from: %s -> /uploads", cfg.Upload.UploadPath)

	// 注册路由
	api.RegisterRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server is running on %s", addr)
	log.Printf("Environment: %s", cfg.Server.Mode)
	log.Printf("Database: %s@%s:%d/%s", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
