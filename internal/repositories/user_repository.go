package repositories

import (
	"user-service/internal/models"
	"user-service/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) Create(user *models.User) error {
	return database.GetDB().Create(user).Error
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetByEmailAndPassword(email, password string) (*models.User, error) {
	var user models.User
	result := database.GetDB().Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
