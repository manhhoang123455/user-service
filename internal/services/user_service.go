package services

import (
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

func (us *UserService) RegisterUser(user *models.User) error {
	err := us.UserRepository.CheckEmailExist(user)
	if err != nil {
		return err
	}
	return us.UserRepository.Create(user)
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	return us.UserRepository.GetByEmail(email)
}

func (us *UserService) CreateUser(user *models.User) error {
	err := us.UserRepository.CheckEmailExist(user)
	if err != nil {
		return err
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return us.UserRepository.Create(user)
}

func (us *UserService) CreateSuperUser(email, password, name string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	superUser := models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Role:     "admin",
	}
	return us.UserRepository.Create(&superUser)
}

func (us *UserService) CreateAuthProvider(authProvider *models.AuthProvider) error {
	return us.UserRepository.CreateAuthProvider(authProvider)
}

func (us *UserService) GetUserByProviderID(providerID string) (*models.User, error) {
	return us.UserRepository.GetUserByProviderID(providerID)
}

func (us *UserService) AuthenticateUser(email string, password string) (*models.User, error) {
	user, err := us.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, utils.ErrInvalidCredentials
	}
	return user, nil
}

func (us *UserService) HandleGoogleCallback(code string) (*models.User, error) {
	googleUser, err := utils.GetGoogleUserInfo(code)
	if err != nil {
		return nil, err
	}

	user, err := us.GetUserByEmail(googleUser.Email)
	if err != nil {
		// If user does not exist, create a new user
		newUser := models.User{
			Email:    googleUser.Email,
			Name:     googleUser.Name,
			Role:     "user",
			Password: "",
		}
		if err := us.UserRepository.Create(&newUser); err != nil {
			return nil, err
		}
		return &newUser, nil
	}

	return user, nil
}

func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	return us.UserRepository.GetByID(id)
}

func (us *UserService) UpdateUser(user *models.User) error {
	err := us.UserRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}
