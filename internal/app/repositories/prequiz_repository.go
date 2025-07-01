package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type PrequizRepositoryInterface interface {
	CreatePrequiz(prequiz *models.Prequiz) (*models.Prequiz, error)
	GetPrequizByID(id uint) (*models.Prequiz, error)
	GetAllPrequizzes() ([]models.Prequiz, error)
	UpdatePrquiz(id uint, prequiz *models.Prequiz) (*models.Prequiz, error)
}

type prequizRepository struct {
	db *gorm.DB
}

func NewPrequizRepository(db *gorm.DB) *prequizRepository {
	return &prequizRepository{db: db}
}

func (r *prequizRepository) CreatePrequiz(prequiz *models.Prequiz) (*models.Prequiz, error) {
	tx := r.db.Begin()

	if err := tx.Create(prequiz).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return prequiz, nil
}

func (r *prequizRepository) GetPrequizByID(id uint) (*models.Prequiz, error) {
	var prequiz models.Prequiz

	err := r.db.Preload("User").Preload("Lesson").First(&prequiz, id).Error
	if err != nil {
		return nil, err // Other error
	}

	return &prequiz, nil
}

func (r *prequizRepository) GetAllPrequizzes() ([]models.Prequiz, error) {
	var prequizzes []models.Prequiz

	err := r.db.Preload("User").Preload("Lesson").Find(&prequizzes).Error
	if err != nil {
		return nil, err // Other error
	}

	return prequizzes, nil
}

func (r *prequizRepository) UpdatePrquiz(id uint, prequiz *models.Prequiz) (*models.Prequiz, error) {
	existingPrequiz, err := r.GetPrequizByID(id)

	if err != nil {
		return nil, err // Other error
	}

	tx := r.db.Begin()

	if err := tx.Model(&models.Prequiz{}).Where("id = ?", id).Updates(existingPrequiz).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return prequiz, nil
}
