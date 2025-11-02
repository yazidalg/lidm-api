package services

import (
	"errors"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type ForgotPasswordServiceInterface interface {
	RequestPasswordReset(req request.ForgotPasswordRequest) error
	VerifyOTP(req request.VerifyOTPRequest) (bool, error)
	ResetPassword(req request.ResetPasswordRequest) error
}

type forgotPasswordService struct {
	forgotPasswordRepo repositories.ForgotPasswordRepositoryInterface
	authRepo           repositories.AuthRepositoryInterface
}

func NewForgotPasswordService(
	forgotPasswordRepo repositories.ForgotPasswordRepositoryInterface,
	authRepo repositories.AuthRepositoryInterface,
) *forgotPasswordService {
	return &forgotPasswordService{
		forgotPasswordRepo: forgotPasswordRepo,
		authRepo:           authRepo,
	}
}

func (s *forgotPasswordService) RequestPasswordReset(req request.ForgotPasswordRequest) error {
	// Find user by email
	user, err := s.authRepo.GetVerifiedUser(req.Email)
	if err != nil {
		return errors.New("user not found or not verified")
	}

	// Generate OTP
	otp := utils.GenerateOTP(6)
	expiryTime := utils.GetExpiryTime()

	// Save OTP to database
	otpRecord := &models.ForgotPassword{
		UserID:    user.ID,
		Email:     user.Email,
		OTP:       otp,
		ExpiresAt: expiryTime,
		Used:      false,
	}

	_, err = s.forgotPasswordRepo.CreateOTP(otpRecord)
	if err != nil {
		return errors.New("failed to create OTP")
	}

	// Send email with OTP
	return utils.SendPasswordResetEmail(user.Email, otp)
}

func (s *forgotPasswordService) VerifyOTP(req request.VerifyOTPRequest) (bool, error) {
	// Check if OTP is valid
	otpRecord, err := s.forgotPasswordRepo.GetValidOTP(req.Email, req.OTP)
	if err != nil {
		return false, err
	}

	if otpRecord == nil {
		return false, errors.New("invalid or expired OTP")
	}

	return true, nil
}

func (s *forgotPasswordService) ResetPassword(req request.ResetPasswordRequest) error {
	// Get user by email
	user, err := s.authRepo.GetByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update user password
	user.Password = string(hashedPassword)
	_, err = s.authRepo.UpdatePassword(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
