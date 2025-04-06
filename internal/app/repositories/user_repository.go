package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserById(id int) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Find(&user, id).Error
	return user, err
}
