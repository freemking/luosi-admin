package controllers

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"path/filepath"

	"admin-backend/utils"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"golang.org/x/image/webp"
)

const (
	maxDimension = 1920
	maxFileSize  = 2 * 1024 * 1024 // 2MB
)

func UploadImage(c *gin.Context) {
	// Get upload type (products, news)
	uploadType := c.DefaultQuery("type", "products")

	// Validate upload type
	validTypes := map[string]bool{"products": true, "news": true, "ads": true, "categories": true, "ckeditor": true}
	if !validTypes[uploadType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid upload type. Use 'products', 'news', 'ads', 'categories' or 'ckeditor'"})
		return
	}

	// Get the file from form
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 2MB limit"})
		return
	}

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif, webp are allowed"})
		return
	}

	// Read file content
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Decode and resize image
	resizedData, err := resizeImageIfNeeded(fileData, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process image: %v", err)})
		return
	}

	// Upload to Qiniu
	relativePath, err := utils.UploadToQiniu(resizedData, uploadType, ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload to CDN: %v", err)})
		return
	}

	// Return the relative path and full URL
	c.JSON(http.StatusOK, gin.H{
		"url":      relativePath,
		"full_url": utils.GetFullURL(relativePath),
		"name":     header.Filename,
	})
}

// resizeImageIfNeeded decodes the image and resizes it if dimensions exceed maxDimension
// Also applies compression to all images
func resizeImageIfNeeded(data []byte, ext string) ([]byte, error) {
	// Decode image based on format
	var img image.Image
	var err error

	reader := bytes.NewReader(data)

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(reader)
	case ".png":
		img, err = png.Decode(reader)
	case ".gif":
		img, err = gif.Decode(reader)
	case ".webp":
		img, err = webp.Decode(reader)
	default:
		return data, nil // Return original if format not supported
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if dimensions exceed maxDimension
	if width > maxDimension || height > maxDimension {
		img = imaging.Fit(img, maxDimension, maxDimension, imaging.Lanczos)
	}

	// Encode back with compression
	var buf bytes.Buffer
	switch ext {
	case ".jpg", ".jpeg":
		// JPEG: quality 85 for good compression with minimal quality loss
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	case ".png":
		// PNG: lossless compression (BestCompression level)
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&buf, img)
	case ".gif":
		// GIF: lossless
		err = gif.Encode(&buf, img, nil)
	case ".webp":
		// webp encoding is not directly supported, convert to jpeg
		ext = ".jpg"
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}
