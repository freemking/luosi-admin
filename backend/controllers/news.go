package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"admin-backend/models"
	"admin-backend/utils"

	"github.com/gin-gonic/gin"
)

// NewsRequest 新闻请求结构
type NewsRequest struct {
	Title       string `json:"title"`
	CoverImage  string `json:"cover_image"`
	PublishDate string `json:"publish_date"` // YYYY-MM-DD format
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Status      int    `json:"status"` // 1: 已发布, 0: 草稿
}

// NewsResponse 新闻响应结构
type NewsResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	CoverImage  string `json:"cover_image"`
	PublishDate string `json:"publish_date"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// convertNewsToResponse converts a news model to response with full URL
func convertNewsToResponse(news models.News) NewsResponse {
	publishDate := ""
	if news.PublishDate.Valid {
		publishDate = news.PublishDate.Time.Format("2006-01-02")
	}
	
	return NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		CoverImage:  utils.GetFullURL(news.CoverImage),
		PublishDate: publishDate,
		Summary:     news.Summary,
		Content:     news.Content,
		Status:      news.Status,
		CreatedAt:   news.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   news.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// GetNewsList 获取新闻列表
func GetNewsList(c *gin.Context) {
	var newsList []models.News

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	var total int64
	models.DB.Model(&models.News{}).Count(&total)

	// 分页查询
	result := models.DB.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&newsList)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news list"})
		return
	}

	// Convert to response with full URLs
	var newsResponses []NewsResponse
	for _, news := range newsList {
		newsResponses = append(newsResponses, convertNewsToResponse(news))
	}

	c.JSON(http.StatusOK, gin.H{
		"news":  newsResponses,
		"total": total,
	})
}

// GetNews 获取单个新闻
func GetNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	var news models.News
	result := models.DB.First(&news, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"news": convertNewsToResponse(news)})
}

// CreateNews 创建新闻
func CreateNews(c *gin.Context) {
	var req NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析发布日期
	var publishDate models.NullTime
	if req.PublishDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.PublishDate)
		if err == nil {
			publishDate = models.NullTime{NullTime: sql.NullTime{Time: parsedDate, Valid: true}}
		}
	}

	// 创建新闻
	news := models.News{
		Title:       req.Title,
		CoverImage:  req.CoverImage,
		PublishDate: publishDate,
		Summary:     req.Summary,
		Content:     req.Content,
		Status:      req.Status,
	}

	result := models.DB.Create(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"news": convertNewsToResponse(news)})
}

// UpdateNews 更新新闻
func UpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	// 检查新闻是否存在
	var news models.News
	result := models.DB.First(&news, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	var req NewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新新闻信息
	if req.Title != "" {
		news.Title = req.Title
	}
	news.CoverImage = req.CoverImage
	if req.PublishDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.PublishDate)
		if err == nil {
			news.PublishDate = models.NullTime{NullTime: sql.NullTime{Time: parsedDate, Valid: true}}
		}
	} else {
		// Clear the publish date if empty (set to NULL)
		news.PublishDate = models.NullTime{NullTime: sql.NullTime{Valid: false}}
	}
	news.Summary = req.Summary
	news.Content = req.Content
	news.Status = req.Status

	result = models.DB.Save(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"news": convertNewsToResponse(news)})
}

// DeleteNews 删除新闻
func DeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid news ID"})
		return
	}

	// 检查新闻是否存在
	var news models.News
	result := models.DB.First(&news, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	// 删除新闻
	result = models.DB.Delete(&news)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

// GetNewsCount 获取新闻总数
func GetNewsCount(c *gin.Context) {
	var count int64

	result := models.DB.Model(&models.News{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
