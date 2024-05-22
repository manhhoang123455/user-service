package services

import (
	"errors"
	"user-service/internal/models"
	"user-service/internal/repositories"
	"user-service/utils"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{UserRepository: ur}
}

func (us *UserService) Register(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return us.UserRepository.CreateUser(user)
}

func (us *UserService) Login(email, password string) (*models.User, error) {
	user, err := us.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("incorrect email or password")
	}
	return user, nil
}
