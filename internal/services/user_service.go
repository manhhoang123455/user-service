package services

import (
	"user-service/internal/models"
	"user-service/internal/repositories"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func NewUserService(ur *repositories.UserRepository) *UserService {
	return &UserService{UserRepository: ur}
}

func (us *UserService) RegisterUser(user *models.User) error {
	return us.UserRepository.Create(user)
}

func (us *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	return us.UserRepository.GetByEmailAndPassword(email, password)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return us.UserRepository.GetByEmail(email)
}

func (us *UserService) CreateUser(user *models.User) error {
	return us.UserRepository.Create(user)
}
