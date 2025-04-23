package repositories

import (
	"errors"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ForgotPasswordRepositoryInterface interface {
	CreateOTP(otp *models.ForgotPassword) (*models.ForgotPassword, error)
	GetValidOTP(email, otp string) (*models.ForgotPassword, error)
	CheckUsedOTP(otp *models.ForgotPassword) (*models.ForgotPassword, error)
}

type forgotPasswordRepository struct {
	db *gorm.DB
}

func NewForgotPasswordRepository(db *gorm.DB) *forgotPasswordRepository {
	return &forgotPasswordRepository{db}
}

func (r *forgotPasswordRepository) CreateOTP(otp *models.ForgotPassword) (*models.ForgotPassword, error) {
	err := r.db.Create(&otp).Error

	return otp, err
}

func (r *forgotPasswordRepository) GetValidOTP(email, otp string) (*models.ForgotPassword, error) {
	var record models.ForgotPassword
	err := r.db.Where("email = ? AND otp = ? AND used = ? AND expires_at > ?", email, otp, false, time.Now().Unix()).
		First(&record).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &record, err
}

func (r *forgotPasswordRepository) CheckUsedOTP(otp *models.ForgotPassword) (*models.ForgotPassword, error) {
	otp.Used = true
	err := r.db.Save(otp).Error
	return otp, err
}
