package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ProgressRepositoryInterface interface {
	UpdateProgress(progress *models.Progress) (*models.Progress, error)
	GetAllProgress() ([]models.Progress, error)
}

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *progressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) UpdateProgress(progress *models.Progress) (*models.Progress, error) {
	tx := r.db.Begin()

	if err := tx.Save(progress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *progressRepository) GetAllProgress() ([]models.Progress, error) {
	var progress []models.Progress

	err := r.db.Preload("Users").Preload("Lessons").Find(&progress).Error
	if err != nil {
		return nil, err
	}

	return progress, nil
}
