package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuizRepositoryInterface interface {
	CreateQuiz(quiz *models.Quiz) (*models.Quiz, error)
	GetQuizByID(id uint) (*models.Quiz, error)
	GetAllQuizzes() ([]models.Quiz, error)
	UpdateQuiz(id uint, quiz *models.Quiz) (*models.Quiz, error)
	DeleteQuiz(id uint) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *quizRepository {
	return &quizRepository{db}
}

func (r *quizRepository) CreateQuiz(quiz *models.Quiz) (*models.Quiz, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Create(quiz).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Fetch the created quiz with all relationships loaded
	var createdQuiz models.Quiz
	if err := r.db.Preload("Questions").Preload("Winner").First(&createdQuiz, quiz.ID).Error; err != nil {
		return nil, err
	}

	return &createdQuiz, nil
}

func (r *quizRepository) GetQuizByID(id uint) (*models.Quiz, error) {
	var quiz models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants").
		Preload("Questions").
		Preload("Winner").
		First(&quiz, id).Error

	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (r *quizRepository) GetAllQuizzes() ([]models.Quiz, error) {
	var quizzes []models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants").
		Preload("Questions").
		Preload("Winner").
		Find(&quizzes).Error

	return quizzes, err
}

func (r *quizRepository) UpdateQuiz(id uint, quiz *models.Quiz) (*models.Quiz, error) {
	// Use a transaction
	tx := r.db.Begin()

	// First get the existing quiz
	var existingQuiz models.Quiz
	if err := tx.First(&existingQuiz, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if quiz.Status != "" {
		existingQuiz.Status = quiz.Status
	}

	// Save the updated quiz
	if err := tx.Save(&existingQuiz).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle questions association if needed
	if len(quiz.Questions) > 0 {
		// Replace questions association
		if err := tx.Model(&existingQuiz).Association("Questions").Replace(quiz.Questions); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return the updated quiz with preloaded associations
	return r.GetQuizByID(id)
}

func (r *quizRepository) DeleteQuiz(id uint) error {
	// Use a transaction
	tx := r.db.Begin()

	var quiz models.Quiz
	if err := tx.First(&quiz, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the quiz (GORM will handle the cascade delete for associated records)
	if err := tx.Delete(&quiz).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
