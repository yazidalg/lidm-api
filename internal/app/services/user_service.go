package services

import (
	"errors"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetUserById(id int) (models.User, error)
	GetUserByIDUint(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUserRole(userID uint, roleName string) (models.User, error)
	UpdateUserAccount(userID uint, name string, email string) (models.User, error)
	DeleteUser(userID uint) error
	DecrementLife(userID uint) error
	ResetLivesIfNeeded(userID uint, maxLives int) error
	AddXP(userID uint, xp int32) error
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
	roleService    RoleServiceInterface
	authRepo       repositories.AuthRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface, roleService RoleServiceInterface, authRepo repositories.AuthRepositoryInterface) *userService {
	return &userService{userRepository, roleService, authRepo}
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

func (s *userService) UpdateUserRole(userID uint, roleName string) (models.User, error) {
	// Get role by name
	role, err := s.roleService.GetRoleByName(roleName)
	if err != nil {
		return models.User{}, err
	}

	return s.userRepository.UpdateUserRole(userID, role.ID)
}

func (s *userService) UpdateUserAccount(userID uint, name string, email string) (models.User, error) {
	// Get existing user
	user, err := s.userRepository.GetUserByIDUint(userID)
	if err != nil {
		return models.User{}, err
	}

	// Check if email is being changed
	if user.Email != email {
		// Check if new email already exists for another user
		existingUser, err := s.authRepo.GetByEmail(email)
		if err == nil && existingUser.ID != userID {
			return models.User{}, errors.New("email already exists")
		}
		// If error is not "record not found", it's a real error
		if err != nil && err != gorm.ErrRecordNotFound {
			return models.User{}, err
		}
	}

	// Update name and email
	user.Name = name
	user.Email = email

	// Save updated user
	err = s.userRepository.UpdateUser(user)
	if err != nil {
		return models.User{}, err
	}

	// Fetch updated user with role
	updatedUser, err := s.userRepository.GetUserByIDUint(userID)
	if err != nil {
		return models.User{}, err
	}

	return *updatedUser, nil
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
