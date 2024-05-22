package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"user-service/config"
	_ "user-service/docs" // Import để load các docs đã tạo bởi swag
	"user-service/internal/models"
	"user-service/internal/routes"
	"user-service/pkg/database"
)

// @title User Service API
// @version 1.0
// @description API documentation for User Service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Load cấu hình
	config.LoadConfig()

	// Khởi tạo kết nối cơ sở dữ liệu với connection pool
	database.InitDB()

	// Tự động migrate các bảng
	db := database.GetDB()
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Tạo router mới
	r := gin.Default()

	// Đăng ký routes
	routes.RegisterRoutes(r)

	// Thêm route cho Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Chạy server trên cổng được chỉ định
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
