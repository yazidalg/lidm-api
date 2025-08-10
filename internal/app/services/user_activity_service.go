package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type UserActivityServiceInterface interface {
	LogActivity(userID uint, activityType, description string, metadata map[string]interface{}, ipAddress, userAgent string) error
	GetUserActivities(userID uint, limit int) ([]models.UserActivity, error)
	GetLastActivity(userID uint) (*models.UserActivity, error)
	GetRecentActivities(limit int) ([]models.UserActivity, error)
	GetMostActiveUsers(limit int) ([]repositories.UserActivitySummary, error)
	GetActivityByType(userID uint, activityType string, limit int) ([]models.UserActivity, error)
	GetTimeSinceLastActivity(userID uint) (string, error)
	UpdateUserStreak(userID uint) error
	GetUserStreak(userID uint) (int, int, error) // current streak, max streak, error
}

type userActivityService struct {
	activityRepo repositories.UserActivityRepositoryInterface
	userRepo     repositories.UserRepositoryInterface
}

func NewUserActivityService(activityRepo repositories.UserActivityRepositoryInterface, userRepo repositories.UserRepositoryInterface) UserActivityServiceInterface {
	return &userActivityService{
		activityRepo: activityRepo,
		userRepo:     userRepo,
	}
}

func (s *userActivityService) LogActivity(userID uint, activityType, description string, metadata map[string]interface{}, ipAddress, userAgent string) error {
	metadataJSON := ""
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err == nil {
			metadataJSON = string(metadataBytes)
		}
	}

	activity := &models.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		MetaData:     metadataJSON,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}

	return s.activityRepo.CreateActivity(activity)
}

func (s *userActivityService) GetUserActivities(userID uint, limit int) ([]models.UserActivity, error) {
	return s.activityRepo.GetUserActivities(userID, limit)
}

func (s *userActivityService) GetLastActivity(userID uint) (*models.UserActivity, error) {
	return s.activityRepo.GetLastActivity(userID)
}

func (s *userActivityService) GetRecentActivities(limit int) ([]models.UserActivity, error) {
	return s.activityRepo.GetRecentActivities(limit)
}

func (s *userActivityService) GetMostActiveUsers(limit int) ([]repositories.UserActivitySummary, error) {
	return s.activityRepo.GetMostActiveUsers(limit)
}

func (s *userActivityService) GetActivityByType(userID uint, activityType string, limit int) ([]models.UserActivity, error) {
	return s.activityRepo.GetActivityByType(userID, activityType, limit)
}

func (s *userActivityService) GetTimeSinceLastActivity(userID uint) (string, error) {
	lastActivity, err := s.activityRepo.GetLastActivity(userID)
	if err != nil {
		return "Tidak pernah", err
	}

	duration := time.Since(lastActivity.CreatedAt)

	if duration < time.Minute {
		return "Baru saja", nil
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 menit yang lalu", nil
		}
		return fmt.Sprintf("%d menit yang lalu", minutes), nil
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 jam yang lalu", nil
		}
		return fmt.Sprintf("%d jam yang lalu", hours), nil
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 hari yang lalu", nil
		}
		return fmt.Sprintf("%d hari yang lalu", days), nil
	} else if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 minggu yang lalu", nil
		}
		return fmt.Sprintf("%d minggu yang lalu", weeks), nil
	} else {
		return lastActivity.CreatedAt.Format("2 Jan 2006"), nil
	}
}

// UpdateUserStreak - Update user streak based on activity
func (s *userActivityService) UpdateUserStreak(userID uint) error {
	// Get user data
	user, err := s.userRepo.GetUserByIDUint(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Jika belum pernah aktif, set streak = 1
	if user.LastActiveDate == nil {
		user.CurrentStreak = 1
		user.MaxStreak = 1
		user.LastActiveDate = &today
		return s.userRepo.UpdateUser(user)
	}

	lastActiveDay := time.Date(user.LastActiveDate.Year(), user.LastActiveDate.Month(), user.LastActiveDate.Day(), 0, 0, 0, 0, user.LastActiveDate.Location())
	daysDiff := int(today.Sub(lastActiveDay).Hours() / 24)

	switch {
	case daysDiff == 0:
		// Sudah aktif hari ini, tidak perlu update streak
		return nil
	case daysDiff == 1:
		// Aktif berturut-turut, tambah streak
		user.CurrentStreak++
		if user.CurrentStreak > user.MaxStreak {
			user.MaxStreak = user.CurrentStreak
		}
	case daysDiff > 1:
		// Putus streak, mulai dari 1 lagi
		user.CurrentStreak = 1
	}

	user.LastActiveDate = &today
	return s.userRepo.UpdateUser(user)
}

// GetUserStreak - Get current and max streak for user
func (s *userActivityService) GetUserStreak(userID uint) (int, int, error) {
	user, err := s.userRepo.GetUserByIDUint(userID)
	if err != nil {
		return 0, 0, err
	}

	// Check if streak should be reset (if user hasn't been active today or yesterday)
	if user.LastActiveDate != nil {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		lastActiveDay := time.Date(user.LastActiveDate.Year(), user.LastActiveDate.Month(), user.LastActiveDate.Day(), 0, 0, 0, 0, user.LastActiveDate.Location())
		daysDiff := int(today.Sub(lastActiveDay).Hours() / 24)

		// Jika lebih dari 1 hari tidak aktif, reset streak
		if daysDiff > 1 {
			user.CurrentStreak = 0
			s.userRepo.UpdateUser(user)
			return 0, user.MaxStreak, nil
		}
	}

	return user.CurrentStreak, user.MaxStreak, nil
}
