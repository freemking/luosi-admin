package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	// 解析多部分表单
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// 获取所有上传的文件
	files := form.File["image"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	// 处理每个文件
	var urls []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif, webp are allowed"})
			return
		}

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String()[:8], ext)
		filepath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		imageURL := fmt.Sprintf("/uploads/%s", filename)
		urls = append(urls, imageURL)
	}

	// 如果只上传了一个文件，保持原来的返回格式
	if len(urls) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"url":  urls[0],
			"name": files[0].Filename,
		})
		return
	}

	// 多个文件返回URL数组
	c.JSON(http.StatusOK, gin.H{
		"urls": urls,
	})
}
