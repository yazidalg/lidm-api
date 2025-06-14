package repositories

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ParticipantRepositoryInterface interface {
	CreateParticipant(participant *models.Participant) (*models.Participant, error)
	GetParticipantByID(id int32) (*models.Participant, error)
	GetAllParticipants() ([]models.Participant, error)
	GetParticipantsByQuizID(quizID uint) ([]models.Participant, error)
	GetParticipantsByUserID(userID uint) ([]models.Participant, error)
	UpdateParticipant(id int32, participant *models.Participant) (*models.Participant, error)
	DeleteParticipant(id int32) error
}

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) *participantRepository {
	return &participantRepository{db}
}

func (r *participantRepository) CreateParticipant(participant *models.Participant) (*models.Participant, error) {
	// Begin transaction for database integrity
	tx := r.db.Begin()

	// First check if the quiz exists
	var quiz models.Quiz
	if err := tx.Where("id = ?", participant.QuizID).First(&quiz).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("quiz with ID %d not found: %w", participant.QuizID, err)
	}

	// Also check if the user exists
	var user models.User
	if err := tx.Where("id = ?", participant.UserID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("user with ID %d not found: %w", participant.UserID, err)
	}

	// Create the participant
	if err := tx.Create(participant).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Now load the complete participant with all relations
	var createdParticipant models.Participant
	if err := r.db.
		Preload("User").
		First(&createdParticipant, participant.ID).Error; err != nil {
		return nil, err
	}

	// Manually set the quiz since preloading didn't work correctly
	createdParticipant.Quiz = quiz

	return &createdParticipant, nil
}

func (r *participantRepository) GetParticipantByID(id int32) (*models.Participant, error) {
	var participant models.Participant

	// Preload related entities
	err := r.db.Preload("User.Leaderboard").
		Preload("Quiz.Questions").
		Preload("Quiz.Participants").
		Preload("Quiz.Winner").
		First(&participant, id).Error

	if err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *participantRepository) GetAllParticipants() ([]models.Participant, error) {
	var participants []models.Participant

	// Preload related entities
	err := r.db.Preload("User.Leaderboard").
		Preload("Quiz.Questions").
		Preload("Quiz.Participants").
		Preload("Quiz.Winner").
		Find(&participants).Error

	return participants, err
}

func (r *participantRepository) GetParticipantsByQuizID(quizID uint) ([]models.Participant, error) {
	var participants []models.Participant

	// Preload related entities
	err := r.db.Preload("User.Leaderboard").
		Preload("Quiz.Questions").
		Preload("Quiz.Participants").
		Preload("Quiz.Winner").
		Where("quiz_id = ?", quizID).
		Find(&participants).Error

	return participants, err
}

func (r *participantRepository) GetParticipantsByUserID(userID uint) ([]models.Participant, error) {
	var participants []models.Participant

	// Preload related entities
	err := r.db.Preload("User.Leaderboard").
		Preload("Quiz.Questions").
		Preload("Quiz.Participants").
		Preload("Quiz.Winner").
		Where("user_id = ?", userID).
		Find(&participants).Error

	return participants, err
}

func (r *participantRepository) UpdateParticipant(id int32, participant *models.Participant) (*models.Participant, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	// Get existing record
	var existingParticipant models.Participant
	if err := tx.First(&existingParticipant, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update the participant fields
	if err := tx.Model(&existingParticipant).Updates(map[string]interface{}{
		"user_id":     participant.UserID,
		"quiz_id":     participant.QuizID,
		"total_score": participant.TotalScore,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload the participant with associations
	var updatedParticipant models.Participant
	if err := r.db.Preload("User.Leaderboard").
		Preload("Quiz.Questions").
		Preload("Quiz.Participants").
		Preload("Quiz.Winner").
		First(&updatedParticipant, id).Error; err != nil {
		return nil, err
	}

	return &updatedParticipant, nil
}

func (r *participantRepository) DeleteParticipant(id int32) error {
	// Use a transaction
	tx := r.db.Begin()

	var participant models.Participant
	if err := tx.First(&participant, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Hard delete the participant (physically remove from database)
	if err := tx.Unscoped().Delete(&participant).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
