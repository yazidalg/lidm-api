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
	UpdateAccount(userID uint, name, email string) (models.User, error)
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

// UpdateAccount - Update user account (name and email only)
func (r *userRepository) UpdateAccount(userID uint, name, email string) (models.User, error) {
	var user models.User
	err := r.db.Model(&user).Where("id = ?", userID).Updates(map[string]interface{}{
		"name":  name,
		"email": email,
	}).Error
	if err != nil {
		return user, err
	}
	// Fetch updated user with role
	err = r.db.Preload("Role").First(&user, userID).Error
	return user, err
}

// DecrementLife - kurangi nyawa user sebanyak 1 jika masih > 0
func (r *userRepository) DecrementLife(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ? AND lives > 0", userID).UpdateColumn("lives", gorm.Expr("lives - 1")).Error
}

// ResetLivesIfNeeded - reset nyawa user ke maxLives jika sudah melewati waktu life_reset_at
func (r *userRepository) ResetLivesIfNeeded(userID uint, maxLives int) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	now := time.Now()
	wib, _ := time.LoadLocation("Asia/Jakarta")
	nowWIB := now.In(wib)

	// Reset jika:
	// 1. LifeResetAt belum pernah diset (null), ATAU
	// 2. Waktu sekarang sudah melewati life_reset_at (sudah waktunya reset)
	needsReset := user.LifeResetAt == nil || nowWIB.After(*user.LifeResetAt)

	if needsReset {
		// Set life_reset_at ke besok jam 00:00 WIB
		nextResetTime := time.Date(nowWIB.Year(), nowWIB.Month(), nowWIB.Day()+1, 0, 0, 0, 0, wib)

		return r.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"lives":         maxLives,
			"life_reset_at": nextResetTime,
		}).Error
	}
	return nil
}

// isSameDay checks if two times are on the same calendar day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// AddXP - tambah XP user
func (r *userRepository) AddXP(userID uint, xp int32) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).UpdateColumn("total_xp", gorm.Expr("total_xp + ?", xp)).Error
}
