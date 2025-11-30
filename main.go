package main

import (
	"html-manager/config"
	"html-manager/handlers"
	"html-manager/models"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 确保数据目录存在
	dbDir := filepath.Dir(config.AppConfig.DBPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal("创建数据目录失败:", err)
	}

	// 初始化数据库
	db, err := models.InitDB(config.AppConfig.DBPath)
	if err != nil {
		log.Fatal("数据库初始化失败:", err)
	}
	defer db.Close()

	// 设置路由
	router := handlers.SetupRouter(db)

	// 加载 HTML 模板
	router.LoadHTMLGlob("templates/*.html")

	// 启动服务器
	addr := ":" + config.AppConfig.Port
	log.Printf("服务器启动在 http://localhost%s", addr)
	log.Printf("管理界面: http://localhost%s/", addr)
	log.Printf("API 端点: http://localhost%s/api", addr)
	log.Printf("页面访问: http://localhost%s/page/:slug", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
