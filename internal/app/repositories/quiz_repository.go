package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuizRepositoryInterface interface {
	GetQuizByID(id int) (*models.Quiz, error)
	GetAllQuizzes() (*[]models.Quiz, error)
	CreateQuiz(quiz models.Quiz) (*models.Quiz, error)
	UpdateQuiz(id int32, quiz models.Quiz) (*models.Quiz, error)
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *quizRepository {
	return &quizRepository{db}
}

func (r *quizRepository) GetQuizByID(id int) (*models.Quiz, error) {
	var quiz models.Quiz
	err := r.db.First(&quiz, "id = ?", id).Error
	return &quiz, err
}

func (r *quizRepository) GetAllQuizzes() (*[]models.Quiz, error) {
	var quizzes []models.Quiz
	err := r.db.Find(&quizzes).Error
	return &quizzes, err
}

func (r *quizRepository) CreateQuiz(quiz models.Quiz) (*models.Quiz, error) {
	err := r.db.Create(&quiz).Error
	return &quiz, err
}

func (r *quizRepository) UpdateQuiz(id int32, quiz models.Quiz) (*models.Quiz, error) {
	err := r.db.Model(&quiz).Where("id = ?", id).Updates(quiz).Error
	return &quiz, err
}
