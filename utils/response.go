package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var EmailAlreadyExistsError = errors.New("email already exists")

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Status: status, Message: message})
}
