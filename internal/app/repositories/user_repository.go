package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserById(id int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUserRole(userID uint, role string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetAllUsers implements UserRepositoryInterface.
func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// UpdateUserRole implements UserRepositoryInterface.
func (r *userRepository) UpdateUserRole(userID uint, role string) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).Where("id = ?", userID).Update("role", role).Error
	return user, err
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Find(&user, id).Error
	return user, err
}
