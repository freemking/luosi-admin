package controllers

import (
	"net/http"
	"strconv"

	"admin-backend/models"
	"admin-backend/utils"

	"github.com/gin-gonic/gin"
)

// ProductRequest 产品请求结构
type ProductRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Category    string            `json:"category" binding:"required"`
	Standard    string            `json:"standard"`
	Material    string            `json:"material"`
	Images      []ProductImageRequest `json:"images"`
}

// ProductImageRequest 产品图片请求结构
type ProductImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
	Order    int    `json:"order"`
}

// UpdateProductRequest 更新产品请求结构
type UpdateProductRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Standard    string            `json:"standard"`
	Material    string            `json:"material"`
	Images      []ProductImageRequest `json:"images"`
}

// ProductResponse 产品响应结构
type ProductResponse struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Standard    string                 `json:"standard"`
	Material    string                 `json:"material"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Images      []ProductImageResponse `json:"images"`
}

// ProductImageResponse 产品图片响应结构
type ProductImageResponse struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	ImageURL  string `json:"image_url"`
	Order     int    `json:"order"`
}

// convertProductToResponse converts a product model to response with full URLs
func convertProductToResponse(product models.Product) ProductResponse {
	response := ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Standard:    product.Standard,
		Material:    product.Material,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	for _, img := range product.Images {
		response.Images = append(response.Images, ProductImageResponse{
			ID:        img.ID,
			ProductID: img.ProductID,
			ImageURL:  utils.GetFullURL(img.ImageURL),
			Order:     img.Order,
		})
	}
	
	return response
}

// GetProducts 获取产品列表
func GetProducts(c *gin.Context) {
	var products []models.Product

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
	models.DB.Model(&models.Product{}).Count(&total)

	// 分页查询
	result := models.DB.Preload("Images").Limit(pageSize).Offset(offset).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	// Convert to response with full URLs
	var productResponses []ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, convertProductToResponse(product))
	}

	c.JSON(http.StatusOK, gin.H{
		"products": productResponses,
		"total":    total,
	})
}

// GetProduct 获取单个产品
func GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	result := models.DB.Preload("Images").First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": convertProductToResponse(product)})
}

// CreateProduct 创建产品
func CreateProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建产品
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Standard:    req.Standard,
		Material:    req.Material,
	}

	result := models.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// 创建产品图片
	for _, img := range req.Images {
		productImage := models.ProductImage{
			ProductID: product.ID,
			ImageURL:  img.ImageURL,
			Order:     img.Order,
		}
		models.DB.Create(&productImage)
	}

	// 重新加载产品及其图片
	models.DB.Preload("Images").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{"product": convertProductToResponse(product)})
}

// UpdateProduct 更新产品
func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// 检查产品是否存在
	var product models.Product
	result := models.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新产品信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Standard != "" {
		product.Standard = req.Standard
	}
	if req.Material != "" {
		product.Material = req.Material
	}

	result = models.DB.Save(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 更新产品图片
	if len(req.Images) > 0 {
		// 删除旧图片
		models.DB.Where("product_id = ?", product.ID).Delete(&models.ProductImage{})

		// 添加新图片
		for _, img := range req.Images {
			productImage := models.ProductImage{
				ProductID: product.ID,
				ImageURL:  img.ImageURL,
				Order:     img.Order,
			}
			models.DB.Create(&productImage)
		}
	}

	// 重新加载产品及其图片
	models.DB.Preload("Images").First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{"product": convertProductToResponse(product)})
}

// DeleteProduct 删除产品
func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// 检查产品是否存在
	var product models.Product
	result := models.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 删除产品图片
	models.DB.Where("product_id = ?", product.ID).Delete(&models.ProductImage{})

	// 删除产品
	result = models.DB.Delete(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetProductCount 获取产品总数
func GetProductCount(c *gin.Context) {
	var count int64

	result := models.DB.Model(&models.Product{}).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
