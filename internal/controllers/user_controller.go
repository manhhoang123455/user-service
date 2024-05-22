package controllers

import (
	"encoding/json"
	"fmt"
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

	err := uc.UserService.CreateUser(&user)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.SendErrorResponse(c, http.StatusConflict, "Email already exists")
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
	c.JSON(http.StatusOK, user)
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
	user, err := uc.UserService.GetUserByEmail(email)
	if err != nil && err.Error() != "record not found" {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if user == nil {
		user = &models.User{
			Email:    email,
			Name:     userInfo["given_name"].(string) + " " + userInfo["family_name"].(string),
			Password: "123456",
		}
		err := uc.UserService.CreateUser(user)
		if err != nil {
			return
		}
	}
	c.JSON(http.StatusOK, user)
}
