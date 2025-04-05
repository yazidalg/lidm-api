package services

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type AuthServiceInterface interface {
	RegisterUser(user request.UserRegisterRequest) (models.User, error)
	LoginUser(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
}

type authService struct {
	authRepository repositories.AuthRepositoryInterface
}

func NewAuthService(authRepository repositories.AuthRepositoryInterface) *authService {
	return &authService{authRepository}
}

func (s *authService) RegisterUser(user request.UserRegisterRequest) (models.User, error) {
	userData := models.User{
		Name:       user.Name,
		Email:      user.Email,
		Class:      user.Class,
		Password:   user.Password,
		IsVerified: false,
		Point:      0,
		Streak:     0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.authRepository.RegisterUser(userData)
}

func (s *authService) LoginUser(id int) (models.User, error) {
	return s.authRepository.LoginUser(id)
}

func (s *authService) GetByEmail(email string) (models.User, error) {
	return s.authRepository.GetByEmail(email)
}
