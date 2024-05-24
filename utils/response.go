package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrorResponse struct
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// TokenResponse struct
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SendErrorResponse sends an error response
func SendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Status: status, Message: message})
}

// SendTokenResponse sends a token response
func SendTokenResponse(c *gin.Context, accessToken string, refreshToken string) {
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
