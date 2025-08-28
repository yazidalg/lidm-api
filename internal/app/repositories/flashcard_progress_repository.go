package repositories

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type FlashcardProgressRepositoryInterface interface {
	Create(progress *models.UserFlashcardProgress) (*models.UserFlashcardProgress, error)
	Update(progress *models.UserFlashcardProgress) (*models.UserFlashcardProgress, error)
	GetByUserAndFlashcard(userID, flashcardID uint) (*models.UserFlashcardProgress, error)
	GetDueByUser(userID uint, dueTime time.Time) ([]models.UserFlashcardProgress, error)
	GetAllByUser(userID uint) ([]models.UserFlashcardProgress, error)
	Delete(id uint) error
}

type flashcardProgressRepository struct {
	db *gorm.DB
}

func NewFlashcardProgressRepository(db *gorm.DB) FlashcardProgressRepositoryInterface {
	return &flashcardProgressRepository{db: db}
}

func (r *flashcardProgressRepository) Create(progress *models.UserFlashcardProgress) (*models.UserFlashcardProgress, error) {
	if err := r.db.Create(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *flashcardProgressRepository) Update(progress *models.UserFlashcardProgress) (*models.UserFlashcardProgress, error) {
	if err := r.db.Save(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *flashcardProgressRepository) GetByUserAndFlashcard(userID, flashcardID uint) (*models.UserFlashcardProgress, error) {
	var progress models.UserFlashcardProgress
	err := r.db.Where("user_id = ? AND flashcard_id = ?", userID, flashcardID).
		Preload("User").
		Preload("Flashcard").
		First(&progress).Error

	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *flashcardProgressRepository) GetDueByUser(userID uint, dueTime time.Time) ([]models.UserFlashcardProgress, error) {
	var progresses []models.UserFlashcardProgress
	err := r.db.Where("user_id = ? AND due <= ?", userID, dueTime).
		Preload("Flashcard").
		Order("due ASC").
		Find(&progresses).Error

	return progresses, err
}

func (r *flashcardProgressRepository) GetAllByUser(userID uint) ([]models.UserFlashcardProgress, error) {
	var progresses []models.UserFlashcardProgress
	err := r.db.Where("user_id = ?", userID).
		Preload("Flashcard").
		Find(&progresses).Error

	return progresses, err
}

func (r *flashcardProgressRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserFlashcardProgress{}, id).Error
}
