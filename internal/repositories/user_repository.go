package repositories

import (
	"user-service/internal/models"
	"user-service/pkg/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	db := database.GetDB()
	return db.Create(user).Error
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
