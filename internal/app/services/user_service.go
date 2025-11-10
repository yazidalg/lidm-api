package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type UserServiceInterface interface {
	GetUserById(id int) (models.User, error)
	GetUserByIDUint(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateAccount(userID uint, name, email string) (models.User, error)
	UpdateUserRole(userID uint, roleName string) (models.User, error)
	DeleteUser(userID uint) error
	DecrementLife(userID uint) error
	ResetLivesIfNeeded(userID uint, maxLives int) error
	AddXP(userID uint, xp int32) error
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

func (s *userService) GetUserByIDUint(id uint) (*models.User, error) {
	return s.userRepository.GetUserByIDUint(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *userService) UpdateAccount(userID uint, name, email string) (models.User, error) {
	return s.userRepository.UpdateAccount(userID, name, email)
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

func (s *userService) DecrementLife(userID uint) error {
	return s.userRepository.DecrementLife(userID)
}

func (s *userService) ResetLivesIfNeeded(userID uint, maxLives int) error {
	return s.userRepository.ResetLivesIfNeeded(userID, maxLives)
}

func (s *userService) AddXP(userID uint, xp int32) error {
	return s.userRepository.AddXP(userID, xp)
}
