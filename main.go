package main

import (
	"html-manager/config"
	"html-manager/handlers"
	"html-manager/models"
	"html-manager/templates"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// 版本信息（构建时注入）
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
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

	// 加载嵌入的 HTML 模板
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		log.Fatal("模板加载失败:", err)
	}
	router.SetHTMLTemplate(tmpl)

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
