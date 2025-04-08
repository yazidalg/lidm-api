package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type ForgotPasswordServiceInterface interface {
	GenerateAndSendOTP(email string) (models.ForgetPasswordModel, error)
	ResetPasswordWithOTP(email, otp, newPassword string) (models.ForgetPasswordModel, error)
}

type forgotPasswordService struct {
	forgotPasswordRepository repositories.ForgotPasswordRepositoryInterface
}

func NewForgotPasswordService(forgotPasswordRepository repositories.ForgotPasswordRepositoryInterface) *forgotPasswordService {
	return &forgotPasswordService{forgotPasswordRepository}
}
