package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type AnswerRepositoryInterface interface {
	CreateAnswer(answer *models.Answer) (*models.Answer, error)
	GetAnswerByID(id int32) (*models.Answer, error)
	GetAllAnswers() ([]models.Answer, error)
	UpdateAnswer(id int32, answer *models.Answer) (*models.Answer, error)
	DeleteAnswer(id int32) error
}

type answerRepository struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) *answerRepository {
	return &answerRepository{db}
}

func (r *answerRepository) CreateAnswer(answer *models.Answer) (*models.Answer, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Create(answer).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return answer, nil
}

func (r *answerRepository) GetAnswerByID(id int32) (*models.Answer, error) {
	var answer models.Answer

	// Preload related entities
	err := r.db.Preload("Question").
		Preload("User").
		Preload("QuizID").
		First(&answer, id).Error

	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func (r *answerRepository) GetAllAnswers() ([]models.Answer, error) {
	var answers []models.Answer

	// Preload related entities
	err := r.db.Preload("Question").
		Preload("User").
		Preload("QuizID").
		Find(&answers).Error

	return answers, err
}

func (r *answerRepository) UpdateAnswer(id int32, answer *models.Answer) (*models.Answer, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	// Get existing record
	var existingAnswer models.Answer
	if err := tx.First(&existingAnswer, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Remove existing associations with quizzes
	if err := tx.Model(&existingAnswer).Association("QuizID").Clear(); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update the answer fields
	if err := tx.Model(&existingAnswer).Updates(map[string]interface{}{
		"question_id":     answer.QuestionID,
		"user_id":         answer.UserID,
		"option_selected": answer.OptionSelected,
		"is_correct":      answer.IsCorrect,
		"score":           answer.Score,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add quiz associations
	if err := tx.Model(&existingAnswer).Association("QuizID").Replace(answer.QuizID); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload the answer with associations
	var updatedAnswer models.Answer
	if err := r.db.Preload("Question").
		Preload("User").
		Preload("QuizID").
		First(&updatedAnswer, id).Error; err != nil {
		return nil, err
	}

	return &updatedAnswer, nil
}

func (r *answerRepository) DeleteAnswer(id int32) error {
	// Use a transaction
	tx := r.db.Begin()

	var answer models.Answer
	if err := tx.First(&answer, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Clear associations before deleting
	if err := tx.Model(&answer).Association("QuizID").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Delete the answer
	if err := tx.Delete(&answer).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
