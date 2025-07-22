package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuizRepositoryInterface interface {
	CreateQuiz(quiz *models.Quiz) (*models.Quiz, error)
	GetQuizByID(id uint) (*models.Quiz, error)
	GetQuizByInviteCode(inviteCode string) (*models.Quiz, error)
	GetAllQuizzes() ([]models.Quiz, error)
	UpdateQuiz(id uint, quiz *models.Quiz) (*models.Quiz, error)
	DeleteQuiz(id uint) error
	GetRandomQuestionsByModule(moduleID uint, count int) ([]models.Question, error)
}

type quizRepository struct {
	db *gorm.DB
}

// generateInviteCode generates a unique 8-character invite code
func (r *quizRepository) generateInviteCode() (string, error) {
	for attempts := 0; attempts < 10; attempts++ {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		code := hex.EncodeToString(bytes)[:8]

		// Check if code already exists
		var count int64
		r.db.Model(&models.Quiz{}).Where("invite_code = ?", code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique invite code")
}

func NewQuizRepository(db *gorm.DB) *quizRepository {
	return &quizRepository{db}
}

func (r *quizRepository) CreateQuiz(quiz *models.Quiz) (*models.Quiz, error) {
	// Generate invite code for multiplayer mode
	if quiz.Mode == "multiplayer" {
		inviteCode, err := r.generateInviteCode()
		if err != nil {
			return nil, err
		}
		quiz.InviteCode = inviteCode
	}

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
	if err := r.db.Preload("Questions").
		Preload("Winner").
		Preload("Host").
		Preload("Module").
		First(&createdQuiz, quiz.ID).Error; err != nil {
		return nil, err
	}

	return &createdQuiz, nil
}

func (r *quizRepository) GetQuizByID(id uint) (*models.Quiz, error) {
	var quiz models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants.User").
		Preload("Questions").
		Preload("Winner").
		Preload("Host").
		Preload("Module").
		First(&quiz, id).Error

	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (r *quizRepository) GetQuizByInviteCode(inviteCode string) (*models.Quiz, error) {
	var quiz models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants.User").
		Preload("Questions").
		Preload("Winner").
		Preload("Host").
		Preload("Module").
		Where("invite_code = ?", inviteCode).
		First(&quiz).Error

	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (r *quizRepository) GetAllQuizzes() ([]models.Quiz, error) {
	var quizzes []models.Quiz

	// Preload related entities
	err := r.db.Preload("Participants.User").
		Preload("Questions").
		Preload("Winner").
		Preload("Host").
		Preload("Module").
		Find(&quizzes).Error

	return quizzes, err
}

func (r *quizRepository) GetRandomQuestionsByModule(moduleID uint, count int) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Where("module_id = ?", moduleID).
		Order("RAND()").
		Limit(count).
		Find(&questions).Error
	return questions, err
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

	if quiz.WinnerID != nil {
		existingQuiz.WinnerID = quiz.WinnerID
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

	if err := tx.Unscoped().Delete(&quiz).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
