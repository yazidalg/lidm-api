package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type UserServiceInterface interface {
	GetUserById(id int) (models.User, error)
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
