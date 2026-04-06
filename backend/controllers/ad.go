package controllers

import (
	"net/http"
	"strconv"
	"time"

	"admin-backend/models"
	"admin-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdRequest 广告请求结构
type AdRequest struct {
	PositionID uint   `json:"position_id" binding:"required"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	ImageURL   string `json:"image_url" binding:"required"`
	LinkURL    string `json:"link_url"`
	Order      int    `json:"order"`
	Status     int    `json:"status"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// UpdateAdRequest 更新广告请求结构
type UpdateAdRequest struct {
	PositionID uint   `json:"position_id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	ImageURL   string `json:"image_url"`
	LinkURL    string `json:"link_url"`
	Order      int    `json:"order"`
	Status     int    `json:"status"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}

// AdResponse 广告响应结构
type AdResponse struct {
	ID           uint               `json:"id"`
	PositionID   uint               `json:"position_id"`
	PositionName string             `json:"position_name"`
	Title        string             `json:"title"`
	Subtitle     string             `json:"subtitle"`
	ImageURL     string             `json:"image_url"`
	LinkURL      string             `json:"link_url"`
	Order        int                `json:"order"`
	Status       int                `json:"status"`
	StartTime    *string            `json:"start_time"`
	EndTime      *string            `json:"end_time"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

// GetAds 获取广告列表
func GetAds(c *gin.Context) {
	var ads []models.Ad
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	positionID := c.DefaultQuery("position_id", "")

	offset := (page - 1) * pageSize

	query := models.DB.Model(&models.Ad{}).Where("deleted_at IS NULL")
	if positionID != "" {
		query = query.Where("position_id = ?", positionID)
	}

	err := query.Preload("AdPosition").Offset(offset).Limit(pageSize).Find(&ads).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告列表失败"})
		return
	}

	var total int64
	countQuery := models.DB.Model(&models.Ad{}).Where("deleted_at IS NULL")
	if positionID != "" {
		countQuery = countQuery.Where("position_id = ?", positionID)
	}
	countQuery.Count(&total)

	var response []AdResponse
	for _, ad := range ads {
		response = append(response, convertAdToResponse(ad))
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetAdCount 获取广告总数
func GetAdCount(c *gin.Context) {
	var total int64
	positionID := c.DefaultQuery("position_id", "")
	query := models.DB.Model(&models.Ad{}).Where("deleted_at IS NULL")
	if positionID != "" {
		query = query.Where("position_id = ?", positionID)
	}
	query.Count(&total)
	c.JSON(http.StatusOK, gin.H{"count": total})
}

// GetAd 获取单个广告
func GetAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Ad
	if err := models.DB.Preload("AdPosition").First(&ad, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": convertAdToResponse(ad)})
}

// CreateAd 创建广告
func CreateAd(c *gin.Context) {
	var req AdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var position models.AdPosition
	if err := models.DB.First(&position, req.PositionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "广告位不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位失败"})
		}
		return
	}

	ad := models.Ad{
		PositionID: req.PositionID,
		Title:      req.Title,
		Subtitle:   req.Subtitle,
		ImageURL:   req.ImageURL,
		LinkURL:    req.LinkURL,
		Order:      req.Order,
		Status:     req.Status,
	}

	if ad.Status == 0 && ad.Status == 1 {
		ad.Status = 1
	}

	if req.StartTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err == nil {
			ad.StartTime = &t
		}
	}
	if req.EndTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err == nil {
			ad.EndTime = &t
		}
	}

	if err := models.DB.Create(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建广告失败"})
		return
	}

	models.DB.Preload("AdPosition").First(&ad, ad.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"data":    convertAdToResponse(ad),
	})
}

// UpdateAd 更新广告
func UpdateAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Ad
	if err := models.DB.First(&ad, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		}
		return
	}

	var req UpdateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PositionID != 0 && req.PositionID != ad.PositionID {
		var position models.AdPosition
		if err := models.DB.First(&position, req.PositionID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "广告位不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告位失败"})
			}
			return
		}
		ad.PositionID = req.PositionID
	}

	if req.Title != "" {
		ad.Title = req.Title
	}
	ad.Subtitle = req.Subtitle
	if req.ImageURL != "" {
		ad.ImageURL = req.ImageURL
	}
	if req.LinkURL != "" {
		ad.LinkURL = req.LinkURL
	}
	if req.Order != 0 || req.Order == 0 {
		ad.Order = req.Order
	}
	if req.Status != 0 || req.Status == 0 {
		ad.Status = req.Status
	}

	if req.StartTime != "" {
		if req.StartTime == "null" || req.StartTime == "" {
			ad.StartTime = nil
		} else {
			t, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
			if err == nil {
				ad.StartTime = &t
			}
		}
	}
	if req.EndTime != "" {
		if req.EndTime == "null" || req.EndTime == "" {
			ad.EndTime = nil
		} else {
			t, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
			if err == nil {
				ad.EndTime = &t
			}
		}
	}

	if err := models.DB.Save(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新广告失败"})
		return
	}

	models.DB.Preload("AdPosition").First(&ad, ad.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"data":    convertAdToResponse(ad),
	})
}

// DeleteAd 删除广告
func DeleteAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Ad
	if err := models.DB.First(&ad, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		}
		return
	}

	if err := models.DB.Delete(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除广告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAdsByPosition 获取指定广告位的广告
func GetAdsByPosition(c *gin.Context) {
	code := c.Param("code")
	var position models.AdPosition
	if err := models.DB.Where("code = ? AND status = 1", code).First(&position).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "广告位不存在或已禁用"})
		return
	}

	var ads []models.Ad
	models.DB.Where("position_id = ? AND status = 1 AND deleted_at IS NULL", position.ID).Order("`order` ASC").Find(&ads)

	var response []AdResponse
	for _, ad := range ads {
		response = append(response, convertAdToResponse(ad))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// convertAdToResponse converts an ad model to response
func convertAdToResponse(ad models.Ad) AdResponse {
	response := AdResponse{
		ID:         ad.ID,
		PositionID: ad.PositionID,
		Title:      ad.Title,
		Subtitle:   ad.Subtitle,
		ImageURL:   utils.GetFullURL(ad.ImageURL),
		LinkURL:    ad.LinkURL,
		Order:      ad.Order,
		Status:     ad.Status,
		CreatedAt:  ad.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  ad.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if ad.StartTime != nil {
		formatted := ad.StartTime.Format("2006-01-02 15:04:05")
		response.StartTime = &formatted
	}
	if ad.EndTime != nil {
		formatted := ad.EndTime.Format("2006-01-02 15:04:05")
		response.EndTime = &formatted
	}

	if ad.AdPosition.ID != 0 {
		response.PositionName = ad.AdPosition.Name
	}

	return response
}