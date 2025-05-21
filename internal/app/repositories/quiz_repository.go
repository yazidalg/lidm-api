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

	return quiz, nil
}

func (r *quizRepository) GetQuizByID(id uint) (*models.Quiz, error) {
	var quiz models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants").
		Preload("Answers").
		Preload("Questions").
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
		Preload("Answers").
		Preload("Questions").
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

	// Update quiz fields
	if quiz.ParticipantID != 0 {
		existingQuiz.ParticipantID = quiz.ParticipantID
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
	if quiz.Questions != nil && len(quiz.Questions) > 0 {
		// Replace questions association
		if err := tx.Model(&existingQuiz).Association("Questions").Replace(quiz.Questions); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Handle answers association if needed
	if quiz.Answers != nil && len(quiz.Answers) > 0 {
		// Replace answers association
		if err := tx.Model(&existingQuiz).Association("Answers").Replace(quiz.Answers); err != nil {
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
