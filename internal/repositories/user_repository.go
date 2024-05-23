package repositories

import (
	"errors"
	"gorm.io/gorm"
	"user-service/internal/models"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) Create(user *models.User) error {
	return ur.DB.Create(user).Error
}

func (ur *UserRepository) CreateAuthProvider(authProvider *models.AuthProvider) error {
	return ur.DB.Create(authProvider).Error
}

func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := ur.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByProviderID(providerID string) (*models.User, error) {
	var user models.User
	if err := ur.DB.Where("provider_id = ?", providerID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) CheckEmailExist(user *models.User) error {
	var count int64
	ur.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return errors.New("email already exists")
	}
	return nil
}
