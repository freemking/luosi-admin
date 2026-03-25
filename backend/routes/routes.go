package routes

import (
	"admin-backend/controllers"
	"admin-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务 - 上传的图片
	r.Static("/uploads", "./uploads")

	// 公开路由
	public := r.Group("/api")
	{
		// 验证码
		public.GET("/captcha", controllers.GetCaptcha)
		// 登录
		public.POST("/login", controllers.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 图片上传
		protected.POST("/upload", controllers.UploadImage)

		// 用户信息
		protected.GET("/user/info", controllers.GetUserInfo)

		// 产品管理
		protected.GET("/products", controllers.GetProducts)
		protected.GET("/products/count", controllers.GetProductCount)
		protected.GET("/products/:id", controllers.GetProduct)
		protected.POST("/products", controllers.CreateProduct)
		protected.PUT("/products/:id", controllers.UpdateProduct)
		protected.DELETE("/products/:id", controllers.DeleteProduct)

		// 反馈管理
		protected.GET("/feedbacks", controllers.GetFeedbacks)
		protected.GET("/feedbacks/count", controllers.GetFeedbackCount)
		protected.GET("/feedbacks/:id", controllers.GetFeedback)

		// 用户管理
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/users/:id", controllers.GetUser)
		protected.POST("/users", controllers.CreateUser)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)
	}
}
