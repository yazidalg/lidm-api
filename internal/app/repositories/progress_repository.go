package repositories

import "gorm.io/gorm"

type ProgressRepositoryInterface interface {
	UpdateProgress(userID uint32, moduleID uint32, lessonID uint32) error
}

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *progressRepository {
	return &progressRepository{db: db}
}
