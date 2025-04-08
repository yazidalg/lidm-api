package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ForgotPasswordRepositoryInterface interface {
	CreateOTP(otp models.ForgetPasswordModel) (models.ForgetPasswordModel, error)
	GetLatestOTP(userId uint, otp string) (models.ForgetPasswordModel, error)
	DeleteOTP(id uint) (models.ForgetPasswordModel, error)
	UpdateUserPassword(userId uint, newPassword string) (models.User, error)
}

type forgotPasswordRepository struct {
	db gorm.DB
}

func NewForgotPasswordRepository(db gorm.DB) *forgotPasswordRepository {
	return &forgotPasswordRepository{db}
}

func (r *forgotPasswordRepository) CreateOTP(otp models.ForgetPasswordModel) (models.ForgetPasswordModel, error) {
	err := r.db.Create(otp).Error

	return otp, err
}

func (r *forgotPasswordRepository) GetLatestOTP(userId uint, otp string) (models.ForgetPasswordModel, error) {
	var forgotPasswordModel models.ForgetPasswordModel
	err := r.db.Where("user_id = ? AND otp = ?", userId, otp).Order("created_at desc").First(&forgotPasswordModel).Error

	return forgotPasswordModel, err
}

func (r *forgotPasswordRepository) DeleteOTP(id uint) (models.ForgetPasswordModel, error) {
	var forgotPasswordModel models.ForgetPasswordModel
	err := r.db.Delete(&forgotPasswordModel, id).Error

	return forgotPasswordModel, err
}

func (r *forgotPasswordRepository) UpdateUserPassword(userId uint, newPassword string) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).Where("id = ?", userId).Update("password", newPassword).Error

	return user, err
}
