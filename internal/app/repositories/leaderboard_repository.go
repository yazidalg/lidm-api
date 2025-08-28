package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type LeaderboardRepositoryInterface interface {
	GetLeaderboardData(moduleID *uint, quizType string) ([]models.Participant, error)
	GetUserScoreByModule(userID uint, moduleID *uint, quizType string) (int64, error)
}

type LeaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db}
}

func (r *LeaderboardRepository) GetLeaderboardData(moduleID *uint, quizType string) ([]models.Participant, error) {
	var participants []models.Participant

	query := r.db.Preload("User").Preload("Quiz")

	// Join with quiz table to filter by mode/type
	if quizType != "" {
		query = query.Joins("JOIN quizzes ON participants.quiz_id = quizzes.id").
			Where("quizzes.mode = ?", quizType)
	}

	// Filter by module if provided
	if moduleID != nil {
		if quizType != "" {
			query = query.Where("quizzes.module_id = ?", *moduleID)
		} else {
			query = query.Joins("JOIN quizzes ON participants.quiz_id = quizzes.id").
				Where("quizzes.module_id = ?", *moduleID)
		}
	}

	// Only get finished participants
	query = query.Where("participants.is_finished = ?", true)

	// Order by total score descending
	err := query.Order("participants.total_score DESC").Find(&participants).Error

	return participants, err
}

func (r *LeaderboardRepository) GetUserScoreByModule(userID uint, moduleID *uint, quizType string) (int64, error) {
	var totalScore int64

	query := r.db.Model(&models.Participant{}).
		Select("COALESCE(SUM(total_score), 0)").
		Where("user_id = ? AND is_finished = ?", userID, true)

	// Join with quiz table for filtering
	if quizType != "" || moduleID != nil {
		query = query.Joins("JOIN quizzes ON participants.quiz_id = quizzes.id")

		if quizType != "" {
			query = query.Where("quizzes.mode = ?", quizType)
		}

		if moduleID != nil {
			query = query.Where("quizzes.module_id = ?", *moduleID)
		}
	}

	err := query.Scan(&totalScore).Error
	return totalScore, err
}
