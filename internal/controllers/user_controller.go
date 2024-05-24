package controllers

import (
	"fmt"
	"net/http"
	"user-service/config"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{UserService: us}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body models.RegisterInput true "User Registration Data"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /register [post]
func (uc *UserController) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := uc.UserService.RegisterUser(&user)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Email already exists")
		} else {
			utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not create user")
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body models.LoginInput true "User Login Data"
// @Success 200 {object} utils.TokenResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /login [post]
func (uc *UserController) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	user, err := uc.UserService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	accessToken, err := utils.GenerateToken(user.ID, user.Role, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Role, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	utils.SendTokenResponse(c, accessToken, refreshToken)
}

// GoogleLogin godoc
// @Summary Login with Google
// @Description Login a user with Google OAuth2
// @Tags User
// @Success 302
// @Router /google-login [get]
func (uc *UserController) GoogleLogin(c *gin.Context) {
	authURL := config.GoogleOAuthConfig.AuthCodeURL("randomstate")
	fmt.Println(authURL)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback godoc
// @Summary Google OAuth2 callback
// @Description Callback for Google OAuth2 login
// @Tags User
// @Success 200 {string} string "data"
// @Failure 500 {object} gin.H
// @Router /google-callback [get]
func (uc *UserController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	user, err := uc.UserService.HandleGoogleCallback(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle Google callback"})
		return
	}

	accessToken, err := utils.GenerateToken(user.ID, user.Role, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Role, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (uc *UserController) CreateSuperUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserService.CreateSuperUser(input.Email, input.Password, input.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Super user created successfully"})
}
