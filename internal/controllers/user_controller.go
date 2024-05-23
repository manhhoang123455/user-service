package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
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
// @Success 200 {object} models.User
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
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, token)
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

	token, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not fetch user data")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	data, _ := io.ReadAll(response.Body)

	var userInfo map[string]interface{}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not fetch user data")
		return
	}

	email, ok := userInfo["email"].(string)
	if !ok {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not fetch user data")
		return
	}
	providerID, ok := userInfo["id"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Provider ID not found in user info"})
		return
	}
	user, err := uc.UserService.GetUserByProviderID(providerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from database"})
		return
	}
	if user == nil {
		// Create new user
		user = &models.User{
			Email: email,
			Name:  userInfo["given_name"].(string) + " " + userInfo["family_name"].(string),
		}
		if err := uc.UserService.CreateUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		authProvider := &models.AuthProvider{
			UserID:     user.ID,
			Provider:   "google",
			ProviderID: providerID,
		}
		if err := uc.UserService.CreateAuthProvider(authProvider); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth provider"})
			return
		}
	}
	c.JSON(http.StatusOK, user)
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
