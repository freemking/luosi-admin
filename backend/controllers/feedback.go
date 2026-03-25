package controllers

import (
	"net/http"
	"strconv"

	"admin-backend/models"

	"github.com/gin-gonic/gin"
)

// GetFeedbacks 获取反馈列表
func GetFeedbacks(c *gin.Context) {
	var feedbacks []models.Feedback

	result := models.DB.Order("created_at DESC").Find(&feedbacks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feedbacks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// GetFeedback 获取单个反馈
func GetFeedback(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback ID"})
		return
	}

	var feedback models.Feedback
	result := models.DB.First(&feedback, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Feedback not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedback": feedback})
}

// GetFeedbackCount 获取反馈总数
func GetFeedbackCount(c *gin.Context) {
	var count int64

	result := models.DB.Model(&models.Feedback{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get feedback count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
