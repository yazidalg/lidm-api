package repositories

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type UserActivityRepositoryInterface interface {
	CreateActivity(activity *models.UserActivity) error
	GetUserActivities(userID uint, limit int) ([]models.UserActivity, error)
	GetLastActivity(userID uint) (*models.UserActivity, error)
	GetRecentActivities(limit int) ([]models.UserActivity, error)
	GetMostActiveUsers(limit int) ([]UserActivitySummary, error)
	GetActivityByType(userID uint, activityType string, limit int) ([]models.UserActivity, error)
}

type userActivityRepository struct {
	db *gorm.DB
}

type UserActivitySummary struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	TotalCount   int64     `json:"total_count"`
	LastActivity time.Time `json:"last_activity"`
}

func NewUserActivityRepository(db *gorm.DB) UserActivityRepositoryInterface {
	return &userActivityRepository{db: db}
}

func (r *userActivityRepository) CreateActivity(activity *models.UserActivity) error {
	return r.db.Create(activity).Error
}

func (r *userActivityRepository) GetUserActivities(userID uint, limit int) ([]models.UserActivity, error) {
	var activities []models.UserActivity

	query := r.db.Where("user_id = ?", userID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&activities).Error
	return activities, err
}

func (r *userActivityRepository) GetLastActivity(userID uint) (*models.UserActivity, error) {
	var activity models.UserActivity

	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&activity).Error

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (r *userActivityRepository) GetRecentActivities(limit int) ([]models.UserActivity, error) {
	var activities []models.UserActivity

	query := r.db.Preload("User").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&activities).Error
	return activities, err
}

func (r *userActivityRepository) GetMostActiveUsers(limit int) ([]UserActivitySummary, error) {
	var results []UserActivitySummary

	query := r.db.Table("user_activities").
		Select("user_id, users.name as username, COUNT(*) as total_count, MAX(user_activities.created_at) as last_activity").
		Joins("JOIN users ON users.id = user_activities.user_id").
		Group("user_id, users.name").
		Order("total_count DESC, last_activity DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Scan(&results).Error
	return results, err
}

func (r *userActivityRepository) GetActivityByType(userID uint, activityType string, limit int) ([]models.UserActivity, error) {
	var activities []models.UserActivity

	query := r.db.Where("user_id = ? AND activity_type = ?", userID, activityType).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&activities).Error
	return activities, err
}
