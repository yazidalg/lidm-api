package repositories

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ProgressRepositoryInterface interface {
	UpdateProgress(id uint, progress *models.Progress) (*models.Progress, error)
	CreateProgress(progress *models.Progress) (*models.Progress, error)
	GetByUserAndLesson(userID uint, lessonID uint) (*models.Progress, error)
	GetProgressById(id uint) (*models.Progress, error)
	GetAllProgress() ([]models.Progress, error)
}

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *progressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) CreateProgress(progress *models.Progress) (*models.Progress, error) {
	tx := r.db.Begin()

	if err := tx.Create(progress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *progressRepository) GetByUserAndLesson(userID uint, lessonID uint) (*models.Progress, error) {
	var progress models.Progress

	err := r.db.Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(&progress).Error
	if err != nil {
		return nil, err // Other error
	}

	return &progress, nil
}

func (r *progressRepository) GetProgressById(id uint) (*models.Progress, error) {
	var progress models.Progress

	err := r.db.Preload("User").Preload("Lesson").First(&progress, id).Error

	if err != nil {
		return nil, err // Other error
	}

	return &progress, nil
}

func (r *progressRepository) UpdateProgress(id uint, progress *models.Progress) (*models.Progress, error) {
	existingProgress, err := r.GetProgressById(id)
	now := time.Now()

	if err != nil {
		return nil, err // Other error
	}

	existingProgress.UserID = progress.UserID
	existingProgress.LessonID = progress.LessonID
	existingProgress.Completed = true
	existingProgress.CompletedAt = &now

	tx := r.db.Begin()

	if err := tx.Save(&existingProgress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return existingProgress, nil
}

func (r *progressRepository) GetAllProgress() ([]models.Progress, error) {
	var progress []models.Progress

	err := r.db.Preload("User").Preload("Lesson").Find(&progress).Error
	if err != nil {
		return nil, err
	}

	return progress, nil
}
