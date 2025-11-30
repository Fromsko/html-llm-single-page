package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	GinMode        string
	SiteName       string
	AuthorName     string
	SupportedSites []string
	DBPath         string
}

var AppConfig *Config

func LoadConfig() {
	// 尝试加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用环境变量或默认值")
	}

	AppConfig = &Config{
		Port:       getEnv("PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		SiteName:   getEnv("SITE_NAME", "HTML Page Manager"),
		AuthorName: getEnv("AUTHOR_NAME", "Your Name"),
		DBPath:     getEnv("DB_PATH", "./data/pages.db"),
	}

	sites := getEnv("SUPPORTED_SITES", "github.com,gitea.com,gitlab.com")
	AppConfig.SupportedSites = strings.Split(sites, ",")
	for i, site := range AppConfig.SupportedSites {
		AppConfig.SupportedSites[i] = strings.TrimSpace(site)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
