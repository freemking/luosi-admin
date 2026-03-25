package main

import (
	"log"

	"admin-backend/models"
	"admin-backend/routes"
	"admin-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	err := models.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化默认用户
	initDefaultUser()

	// 设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	log.Println("Admin backend server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDefaultUser 初始化默认用户
func initDefaultUser() {
	// 检查是否已有用户
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	// 创建默认超级管理员用户
	hashedPassword, err := utils.HashPassword("1qaz@WSX")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	user := models.User{
		Username: "admin",
		Password: hashedPassword,
		Role:     "super",
	}

	result := models.DB.Create(&user)
	if result.Error != nil {
		log.Printf("Failed to create default user: %v", result.Error)
	} else {
		log.Println("Default admin user created successfully")
	}
}
