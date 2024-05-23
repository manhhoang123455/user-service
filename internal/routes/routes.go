package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/internal/controllers"
	middleware "user-service/internal/middlewares"
	"user-service/internal/repositories"
	"user-service/internal/services"
	"user-service/pkg/database"
)

func RegisterRoutes(r *gin.Engine) {
	userRepo := repositories.NewUserRepository(database.GetDB())
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.GET("/google-login", userController.GoogleLogin)
	r.GET("/google-callback", userController.GoogleCallback)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("/init-superuser", userController.CreateSuperUser)
	}
}
