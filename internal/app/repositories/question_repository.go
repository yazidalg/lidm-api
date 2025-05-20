package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuestionRepositoryInterface interface {
	CreateQuestion(question *models.Question) (*models.Question, error)
	GetQuestionByID(id uint) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	UpdateQuestion(id int, question *models.Question) (*models.Question, error)
	DeleteQuestion(id uint) error
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *questionRepository {
	return &questionRepository{db}
}

func (r *questionRepository) CreateQuestion(question *models.Question) (*models.Question, error) {
	err := r.db.Create(question).Error
	return question, err
}

func (r *questionRepository) GetQuestionByID(id uint) (*models.Question, error) {
	var question models.Question
	err := r.db.First(&question, id).Error
	return &question, err
}

func (r *questionRepository) GetAllQuestions() ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Find(&questions).Error
	return questions, err
}

func (r *questionRepository) UpdateQuestion(id int, question *models.Question) (*models.Question, error) {
	err := r.db.Model(&models.Question{}).Where("id = ?", id).Updates(question).Error
	return question, err
}

func (r *questionRepository) DeleteQuestion(id uint) error {
	return r.db.Delete(&models.Question{}, id).Error
}
