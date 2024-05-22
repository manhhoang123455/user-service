package routes

import (
	"github.com/gin-gonic/gin"
	"user-service/internal/controllers"
	"user-service/internal/repositories"
	"user-service/internal/services"
)

func RegisterRoutes(r *gin.Engine) {
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	//r.GET("/auth/google", userController.GoogleAuth)
	//r.GET("/auth/google/callback", userController.GoogleAuthCallback)
}
