package services

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

type AuthServiceInterface interface {
	RegisterUser(user request.UserRegisterRequest) (models.User, error)
	LoginUser(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByVerificationToken(token string) (models.User, error)
	VerifyUser(user models.User) (models.User, error)
	GetVerifiedUser(email string) (models.User, error)
}

type authService struct {
	authRepository repositories.AuthRepositoryInterface
	roleService    RoleServiceInterface
}

func NewAuthService(authRepository repositories.AuthRepositoryInterface, roleService RoleServiceInterface) *authService {
	return &authService{authRepository, roleService}
}

func (s *authService) RegisterUser(user request.UserRegisterRequest) (models.User, error) {
	token := utils.GenerateToken()

	// Get role by name, default to 'user' if not provided
	roleName := user.RoleName
	if roleName == "" {
		roleName = models.RoleUserName
	}

	role, err := s.roleService.GetRoleByName(roleName)
	if err != nil {
		// If role not found, use default user role
		role, err = s.roleService.GetRoleByName(models.RoleUserName)
		if err != nil {
			return models.User{}, err
		}
	}

	userData := models.User{
		Name:              user.Name,
		Email:             user.Email,
		Class:             user.Class,
		Password:          user.Password,
		RoleID:            role.ID, // Use role ID instead of role string
		IsVerified:        false,
		VerificationToken: token,
		Point:             0,
		TotalXP:           0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	createdUser, err := s.authRepository.RegisterUser(userData)

	if err != nil {
		return userData, err
	}

	return createdUser, nil
}

func (s *authService) LoginUser(id int) (models.User, error) {
	return s.authRepository.LoginUser(id)
}

func (s *authService) GetByEmail(email string) (models.User, error) {
	return s.authRepository.GetByEmail(email)
}

func (s *authService) GetByVerificationToken(token string) (models.User, error) {
	return s.authRepository.GetByVerificationToken(token)
}

func (s *authService) VerifyUser(user models.User) (models.User, error) {
	return s.authRepository.VerifyUser(user)
}

func (s *authService) GetVerifiedUser(email string) (models.User, error) {
	return s.authRepository.GetVerifiedUser(email)
}
