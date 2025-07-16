package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type UserServiceInterface interface {
	GetUserById(id int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUserRole(userID uint, role string) (models.User, error)
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) *userService {
	return &userService{userRepository}
}

func (s *userService) GetUserById(id int) (models.User, error) {
	return s.userRepository.GetUserById(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *userService) UpdateUserRole(userID uint, role string) (models.User, error) {
	return s.userRepository.UpdateUserRole(userID, role)
}
