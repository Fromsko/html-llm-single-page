package handlers

import (
	"database/sql"
	"fmt"
	"html-manager/config"
	"html-manager/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	db *sql.DB
}

func NewAPIHandler(db *sql.DB) *APIHandler {
	return &APIHandler{db: db}
}

func SetupRouter(db *sql.DB) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.GinMode)

	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiHandler := NewAPIHandler(db)

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/pages", apiHandler.GetPages)
		api.GET("/page/:slug", apiHandler.GetPage)    // 通过 slug 获取页面
		api.GET("/pages/:id", apiHandler.GetPageByID) // 通过 ID 获取页面，用于编辑
		api.POST("/pages", apiHandler.CreatePage)
		api.PUT("/pages/:id", apiHandler.UpdatePage)
		api.DELETE("/pages/:id", apiHandler.DeletePage)
		api.POST("/upload", apiHandler.UploadFile)
		api.GET("/download/:id", apiHandler.DownloadFile)
	}

	// 页面路由 - 展示 HTML 页面
	r.GET("/page/:slug", apiHandler.ServePage)

	// 主页 - 管理界面
	r.GET("/", apiHandler.ServeAdmin)

	// 静态文件
	r.Static("/static", "./static")

	return r
}

func (h *APIHandler) GetPages(c *gin.Context) {
	pages, err := models.GetAllPages(h.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pages})
}

func (h *APIHandler) GetPage(c *gin.Context) {
	slug := c.Param("slug")
	page, err := models.GetPage(h.db, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *APIHandler) GetPageByID(c *gin.Context) {
	id := c.Param("id")
	pageID := 0
	fmt.Sscanf(id, "%d", &pageID)

	// 查询数据库获取页面信息
	query := `SELECT id, title, slug, content, description, author, created_at, updated_at FROM pages WHERE id = ?`
	var page models.Page
	err := h.db.QueryRow(query, pageID).Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.Description, &page.Author, &page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *APIHandler) CreatePage(c *gin.Context) {
	var page models.Page
	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置作者
	if page.Author == "" {
		page.Author = config.AppConfig.AuthorName
	}

	// 生成 slug
	if page.Slug == "" {
		page.Slug = generateSlug(page.Title)
	}

	if err := page.Create(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": page})
}

func (h *APIHandler) UpdatePage(c *gin.Context) {
	id := c.Param("id")

	var updateData struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取现有页面
	query := `SELECT id, title, slug, content, description, author, created_at, updated_at FROM pages WHERE id = ?`
	var page models.Page
	err := h.db.QueryRow(query, id).Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.Description, &page.Author, &page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// 更新数据
	page.Title = updateData.Title
	page.Content = updateData.Content
	page.Description = updateData.Description

	if err := page.Update(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

func (h *APIHandler) DeletePage(c *gin.Context) {
	id := c.Param("id")
	pageID := 0
	fmt.Sscanf(id, "%d", &pageID)

	if err := models.DeletePage(h.db, pageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

func (h *APIHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件类型
	if !strings.HasSuffix(file.Filename, ".html") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only HTML files are allowed"})
		return
	}

	// 读取文件内容
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileContent.Close()

	content := make([]byte, file.Size)
	_, err = fileContent.Read(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建页面
	title := strings.TrimSuffix(file.Filename, ".html")
	slug := generateSlug(title)

	page := models.Page{
		Title:   title,
		Slug:    slug,
		Content: string(content),
		Author:  config.AppConfig.AuthorName,
	}

	if err := page.Create(h.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": page})
}

func (h *APIHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	pageID := 0
	fmt.Sscanf(id, "%d", &pageID)

	query := `SELECT title, slug, content FROM pages WHERE id = ?`
	var title, slug, content string
	err := h.db.QueryRow(query, pageID).Scan(&title, &slug, &content)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	filename := slug + ".html"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, content)
}

func (h *APIHandler) ServePage(c *gin.Context) {
	slug := c.Param("slug")
	page, err := models.GetPage(h.db, slug)
	if err != nil {
		c.String(http.StatusNotFound, "Page not found")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, page.Content)
}

func (h *APIHandler) ServeAdmin(c *gin.Context) {
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"SiteName": config.AppConfig.SiteName,
		"Author":   config.AppConfig.AuthorName,
	})
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// 保留字母、数字和连字符
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slugResult := result.String()

	// 如果 slug 为空，使用时间戳
	if slugResult == "" {
		slugResult = fmt.Sprintf("page-%d", time.Now().Unix())
	}

	return slugResult
}
