package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	LastLogin time.Time `json:"last_login" gorm:"default:null"`
	Role      string    `gorm:"not null"`
}

type AuthProvider struct {
	gorm.Model
	UserID     uint   `gorm:"not null"`
	Provider   string `gorm:"not null"` // e.g., 'google', 'facebook'
	ProviderID string `gorm:"not null"` // e.g., Google ID, Facebook ID
}

type RegisterInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
