package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuizRepositoryInterface interface {
	GetQuizByID(id int) (*models.Quiz, error)
	GetAllQuizzes() (*[]models.Quiz, error)
	CreateQuiz(quiz models.Quiz) (*models.Quiz, error)
	UpdateQuiz(quiz models.Quiz) (*models.Quiz, error)
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *quizRepository {
	return &quizRepository{db}
}

func (r *quizRepository) GetQuizByID(id int) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.Find(&quiz, id).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) GetAllQuizzes() (*[]models.Quiz, error) {
	var quizzes []models.Quiz
	err := r.db.Find(&quizzes).Error
	if err != nil {
		return nil, err
	}
	return &quizzes, nil
}

func (r *quizRepository) CreateQuiz(quiz models.Quiz) (*models.Quiz, error) {
	err := r.db.Create(&quiz).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) UpdateQuiz(quiz models.Quiz) (*models.Quiz, error) {
	err := r.db.Save(&quiz).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}
