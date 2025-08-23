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
	GetUserPrequizAnswers(userID uint) ([]models.PrequizUserAnswer, error)
	SubmitPrequizAnswer(answer *models.PrequizUserAnswer) (*models.PrequizUserAnswer, error)
	HasUserAnsweredPrequiz(userID uint, prequizID uint) (bool, error)
	GetPrequizzesByModule(moduleID uint, limit int) ([]models.Prequiz, error)
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

	err := r.db.First(&prequiz, id).Error
	if err != nil {
		return nil, err // Other error
	}

	return &prequiz, nil
}

func (r *prequizRepository) GetAllPrequizzes() ([]models.Prequiz, error) {
	var prequizzes []models.Prequiz

	err := r.db.Find(&prequizzes).Error
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

func (r *prequizRepository) GetUserPrequizAnswers(userID uint) ([]models.PrequizUserAnswer, error) {
	var answers []models.PrequizUserAnswer

	err := r.db.Preload("User").Preload("Prequiz").Where("user_id = ?", userID).Find(&answers).Error
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (r *prequizRepository) SubmitPrequizAnswer(answer *models.PrequizUserAnswer) (*models.PrequizUserAnswer, error) {
	tx := r.db.Begin()

	if err := tx.Create(answer).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return answer, nil
}

func (r *prequizRepository) HasUserAnsweredPrequiz(userID uint, prequizID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.PrequizUserAnswer{}).Where("user_id = ? AND prequiz_id = ?", userID, prequizID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *prequizRepository) GetPrequizzesByModule(moduleID uint, limit int) ([]models.Prequiz, error) {
	var prequizzes []models.Prequiz

	query := r.db.Where("module_id = ?", moduleID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&prequizzes).Error
	if err != nil {
		return nil, err
	}

	return prequizzes, nil
}
