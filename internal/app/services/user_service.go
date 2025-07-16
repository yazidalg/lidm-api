package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type UserServiceInterface interface {
	GetUserById(id int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUserRole(userID uint, roleName string) (models.User, error)
	DeleteUser(userID uint) error
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
	roleService    RoleServiceInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface, roleService RoleServiceInterface) *userService {
	return &userService{userRepository, roleService}
}

func (s *userService) GetUserById(id int) (models.User, error) {
	return s.userRepository.GetUserById(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *userService) UpdateUserRole(userID uint, roleName string) (models.User, error) {
	// Get role by name
	role, err := s.roleService.GetRoleByName(roleName)
	if err != nil {
		return models.User{}, err
	}

	return s.userRepository.UpdateUserRole(userID, role.ID)
}

func (s *userService) DeleteUser(userID uint) error {
	return s.userRepository.DeleteUser(userID)
}
