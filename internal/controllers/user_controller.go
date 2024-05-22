package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/utils"
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

	err := uc.UserService.Register(&user)
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

	user, err := uc.UserService.Login(input.Email, input.Password)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	c.JSON(http.StatusOK, user)
}

//// GoogleAuth godoc
//// @Summary Google Auth
//// @Description Redirect to Google OAuth consent screen
//// @Tags user
//// @Produce  json
//// @Success 302
//// @Router /auth/google [get]
//func (uc *UserController) GoogleAuth(c *gin.Context) {
//	url := config.GoogleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
//	c.Redirect(http.StatusFound, url)
//}
//
//// GoogleAuthCallback godoc
//// @Summary Google Auth Callback
//// @Description Callback for Google OAuth consent screen
//// @Tags user
//// @Produce  json
//// @Success 302
//// @Router /auth/google/callback [get]
//func (uc *UserController) GoogleAuthCallback(c *gin.Context) {
//	code := c.Query("code")
//	token, err := uc.UserService.GetGoogleToken(code)
//	if err != nil {
//		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not get Google token")
//		return
//	}
//
//	user, err := uc.UserService.GetGoogleUserInfo(token)
//	if err != nil {
//		utils.SendErrorResponse(c, http.StatusInternalServerError, "Could not get Google user info")
//		return
//	}
//
//	c.JSON(http.StatusOK, user)
//}
