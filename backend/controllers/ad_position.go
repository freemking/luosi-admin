package controllers

import (
	"net/http"
	"strconv"

	"admin-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdPositionRequest 广告位请求结构
type AdPositionRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Status      int    `json:"status"`
}

// UpdateAdPositionRequest 更新广告位请求结构
type UpdateAdPositionRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Status      int    `json:"status"`
}

// AdPositionResponse 广告位响应结构
type AdPositionResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GetAdPositions 获取广告位列表
func GetAdPositions(c *gin.Context) {
	var positions []models.AdPosition
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	err := models.DB.Offset(offset).Limit(pageSize).Find(&positions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位列表失败"})
		return
	}

	var total int64
	models.DB.Model(&models.AdPosition{}).Where("deleted_at IS NULL").Count(&total)

	var response []AdPositionResponse
	for _, pos := range positions {
		response = append(response, convertAdPositionToResponse(pos))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetAdPositionCount 获取广告位总数
func GetAdPositionCount(c *gin.Context) {
	var total int64
	models.DB.Model(&models.AdPosition{}).Where("deleted_at IS NULL").Count(&total)
	c.JSON(http.StatusOK, gin.H{"count": total})
}

// GetAdPosition 获取单个广告位
func GetAdPosition(c *gin.Context) {
	id := c.Param("id")
	var position models.AdPosition
	if err := models.DB.First(&position, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告位不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": convertAdPositionToResponse(position)})
}

// CreateAdPosition 创建广告位
func CreateAdPosition(c *gin.Context) {
	var req AdPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.AdPosition
	if err := models.DB.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "广告位编码已存在"})
		return
	}

	position := models.AdPosition{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Width:       req.Width,
		Height:      req.Height,
		Status:      req.Status,
	}

	if position.Status == 0 && position.Status == 1 {
		position.Status = 1
	}

	if err := models.DB.Create(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建广告位失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    convertAdPositionToResponse(position),
	})
}

// UpdateAdPosition 更新广告位
func UpdateAdPosition(c *gin.Context) {
	id := c.Param("id")
	var position models.AdPosition
	if err := models.DB.First(&position, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告位不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位失败"})
		}
		return
	}

	var req UpdateAdPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Code != "" && req.Code != position.Code {
		var existing models.AdPosition
		if err := models.DB.Where("code = ? AND id != ?", req.Code, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "广告位编码已存在"})
			return
		}
		position.Code = req.Code
	}

	if req.Name != "" {
		position.Name = req.Name
	}
	if req.Description != "" {
		position.Description = req.Description
	}
	if req.Width != 0 {
		position.Width = req.Width
	}
	if req.Height != 0 {
		position.Height = req.Height
	}
	if req.Status != 0 || req.Status == 0 {
		position.Status = req.Status
	}

	if err := models.DB.Save(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新广告位失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    convertAdPositionToResponse(position),
	})
}

// DeleteAdPosition 删除广告位
func DeleteAdPosition(c *gin.Context) {
	id := c.Param("id")
	var position models.AdPosition
	if err := models.DB.First(&position, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告位不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位失败"})
		}
		return
	}

	if err := models.DB.Delete(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除广告位失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAllAdPositions 获取所有启用的广告位（用于下拉选择）
func GetAllAdPositions(c *gin.Context) {
	var positions []models.AdPosition
	err := models.DB.Where("status = 1").Find(&positions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位列表失败"})
		return
	}

	var response []AdPositionResponse
	for _, pos := range positions {
		response = append(response, convertAdPositionToResponse(pos))
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// convertAdPositionToResponse converts an ad position model to response
func convertAdPositionToResponse(position models.AdPosition) AdPositionResponse {
	return AdPositionResponse{
		ID:          position.ID,
		Code:        position.Code,
		Name:        position.Name,
		Description: position.Description,
		Width:       position.Width,
		Height:      position.Height,
		Status:      position.Status,
		CreatedAt:   position.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   position.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}