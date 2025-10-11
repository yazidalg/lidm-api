package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type FlashcardRepositoryInterface interface {
	GetAllFlashcards() ([]models.Flashcard, error)
	GetFlashcardByID(id uint) (*models.Flashcard, error)
	GetFlashcardsByModule(moduleID uint) ([]models.Flashcard, error)
	CreateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
	UpdateFlashcard(id uint, flashcard *models.Flashcard) (*models.Flashcard, error)
	DeleteFlashcard(id uint) error
}

type flashcardRepository struct {
	db *gorm.DB
}

func NewFlashcardRepository(db *gorm.DB) FlashcardRepositoryInterface {
	return &flashcardRepository{db: db}
}

func (r *flashcardRepository) GetAllFlashcards() ([]models.Flashcard, error) {
	var flashcards []models.Flashcard
	err := r.db.Order("module_id ASC, `order` ASC").Find(&flashcards).Error
	return flashcards, err
}

func (r *flashcardRepository) GetFlashcardByID(id uint) (*models.Flashcard, error) {
	var flashcard models.Flashcard
	err := r.db.First(&flashcard, id).Error
	if err != nil {
		return nil, err
	}
	return &flashcard, nil
}

func (r *flashcardRepository) GetFlashcardsByModule(moduleID uint) ([]models.Flashcard, error) {
	var flashcards []models.Flashcard
	err := r.db.Where("module_id = ?", moduleID).Order("`order` ASC").Find(&flashcards).Error
	return flashcards, err
}

func (r *flashcardRepository) CreateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	if err := r.db.Create(flashcard).Error; err != nil {
		return nil, err
	}
	return flashcard, nil
}

func (r *flashcardRepository) UpdateFlashcard(id uint, flashcard *models.Flashcard) (*models.Flashcard, error) {
	if err := r.db.Model(&models.Flashcard{}).Where("id = ?", id).Updates(flashcard).Error; err != nil {
		return nil, err
	}
	return flashcard, nil
}

func (r *flashcardRepository) DeleteFlashcard(id uint) error {
	return r.db.Delete(&models.Flashcard{}, id).Error
}
