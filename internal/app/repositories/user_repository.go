package repositories

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserById(id int) (models.User, error)
	GetUserByIDUint(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	UpdateUserRole(userID uint, roleID uint) (models.User, error)
	DeleteUser(userID uint) error
	DecrementLife(userID uint) error
	ResetLivesIfNeeded(userID uint, maxLives int) error
	AddXP(userID uint, xp int32) error
}

type userRepository struct {
	db *gorm.DB
}

// GetAllUsers implements UserRepositoryInterface.
func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Role").Find(&users).Error
	return users, err
}

// UpdateUserRole implements UserRepositoryInterface.
func (r *userRepository) UpdateUserRole(userID uint, roleID uint) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).Where("id = ?", userID).Update("role_id", roleID).Error
	if err != nil {
		return user, err
	}
	// Fetch updated user with role
	err = r.db.Preload("Role").First(&user, userID).Error
	return user, err
}

func (r *userRepository) DeleteUser(userID uint) error {
	return r.db.Delete(&models.User{}, userID).Error
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Find(&user, id).Error
	return user, err
}

// GetUserByIDUint - Get user by uint ID
func (r *userRepository) GetUserByIDUint(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser - Update user data
func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// DecrementLife - kurangi nyawa user sebanyak 1 jika masih > 0
func (r *userRepository) DecrementLife(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ? AND lives > 0", userID).UpdateColumn("lives", gorm.Expr("lives - 1")).Error
}

// ResetLivesIfNeeded - reset nyawa user ke maxLives jika hari sudah berganti
func (r *userRepository) ResetLivesIfNeeded(userID uint, maxLives int) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	now := time.Now()
	if user.LifeResetAt == nil || (user.LifeResetAt.Year() != now.Year() || user.LifeResetAt.YearDay() != now.YearDay()) {
		return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"lives":         maxLives,
			"life_reset_at": now,
		}).Error
	}
	return nil
}

// AddXP - tambah XP user
func (r *userRepository) AddXP(userID uint, xp int32) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("total_xp", gorm.Expr("total_xp + ?", xp)).Error
}
