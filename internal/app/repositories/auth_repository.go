package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	RegisterUser(user models.User) (models.User, error)
	LoginUser(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByVerificationToken(token string) (models.User, error)
	VerifyUser(user models.User) (models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

func (r *authRepository) RegisterUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *authRepository) LoginUser(id int) (models.User, error) {
	var user models.User

	err := r.db.Find(&user, id).Error

	return user, err
}

func (r *authRepository) GetByEmail(email string) (models.User, error) {
	var user models.User

	err := r.db.First(&user, "email = ?", email).Error

	return user, err
}

func (r *authRepository) GetByVerificationToken(token string) (models.User, error) {
	var user models.User

	err := r.db.Where("verification_token = ?", token).First(&user).Error

	return user, err
}

func (r *authRepository) VerifyUser(user models.User) (models.User, error) {
	user.IsVerified = true
	user.VerificationToken = ""
	err := r.db.Save(&user).Error
	return user, err
}
