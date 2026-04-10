package controllers

import (
	"net/http"
	"os/exec"
	"strconv"

	"admin-backend/models"
	"admin-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ImageURL    string `json:"image_url"`
	Order       int    `json:"order"`
	Status      int    `json:"status"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ImageURL    string `json:"image_url"`
	Order       int    `json:"order"`
	Status      int    `json:"status"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ImageURL    string `json:"image_url"`
	Order       int    `json:"order"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func restartService() error {
	cmd := exec.Command("sudo", "systemctl", "restart", "go-luosi")
	return cmd.Start()
}

func GetCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var categories []models.Category
	var total int64

	query := models.DB.Model(&models.Category{}).Where("deleted_at IS NULL")
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("`order` ASC, id DESC").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = CategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description,
			Icon:        utils.GetFullURL(cat.Icon),
			ImageURL:    utils.GetFullURL(cat.ImageURL),
			Order:       cat.Order,
			Status:      cat.Status,
			CreatedAt:   cat.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   cat.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": total,
		"page":  page,
	})
}

func GetAllCategories(c *gin.Context) {
	var categories []models.Category
	if err := models.DB.Where("deleted_at IS NULL AND status = ?", 1).Order("`order` ASC, id DESC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = CategoryResponse{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description,
			Icon:        utils.GetFullURL(cat.Icon),
			ImageURL:    utils.GetFullURL(cat.ImageURL),
			Order:       cat.Order,
			Status:      cat.Status,
			CreatedAt:   cat.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   cat.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		}
		return
	}

	response := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Icon:        utils.GetFullURL(category.Icon),
		ImageURL:    utils.GetFullURL(category.ImageURL),
		Order:       category.Order,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCategory models.Category
	if err := models.DB.Where("slug = ? AND deleted_at IS NULL", req.Slug).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类标识已存在"})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		ImageURL:    req.ImageURL,
		Order:       req.Order,
		Status:      req.Status,
	}

	if category.Status == 0 {
		category.Status = 1
	}

	if err := models.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	restartService()

	response := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Icon:        utils.GetFullURL(category.Icon),
		ImageURL:    utils.GetFullURL(category.ImageURL),
		Order:       category.Order,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    response,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		}
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Slug != "" && req.Slug != category.Slug {
		var existingCategory models.Category
		if err := models.DB.Where("slug = ? AND deleted_at IS NULL AND id != ?", req.Slug, category.ID).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "分类标识已存在"})
			return
		}
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	category.Description = req.Description
	category.Icon = req.Icon
	category.ImageURL = req.ImageURL
	category.Order = req.Order
	category.Status = req.Status

	if err := models.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	restartService()

	response := CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Icon:        utils.GetFullURL(category.Icon),
		ImageURL:    utils.GetFullURL(category.ImageURL),
		Order:       category.Order,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    response,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := models.DB.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类失败"})
		}
		return
	}

	if err := models.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	restartService()

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetCategoryCount(c *gin.Context) {
	var total int64
	models.DB.Model(&models.Category{}).Where("deleted_at IS NULL").Count(&total)
	c.JSON(http.StatusOK, gin.H{"count": total})
}
