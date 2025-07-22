package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuizSessionRepositoryInterface interface {
	CreateQuizSession(session *models.QuizSession) (*models.QuizSession, error)
	GetQuizSessionsByQuizID(quizID uint) ([]models.QuizSession, error)
	GetQuizSessionsByParticipantID(participantID uint) ([]models.QuizSession, error)
	GetQuizSession(quizID, participantID, questionID uint) (*models.QuizSession, error)
	UpdateQuizSession(session *models.QuizSession) (*models.QuizSession, error)
	GetQuizSessionStats(quizID uint) (map[uint]models.Participant, error)
}

type quizSessionRepository struct {
	db *gorm.DB
}

func NewQuizSessionRepository(db *gorm.DB) *quizSessionRepository {
	return &quizSessionRepository{db}
}

func (r *quizSessionRepository) CreateQuizSession(session *models.QuizSession) (*models.QuizSession, error) {
	err := r.db.Create(session).Error
	if err != nil {
		return nil, err
	}

	// Load relationships
	err = r.db.Preload("Quiz").Preload("Participant.User").Preload("Question").First(session, session.ID).Error
	return session, err
}

func (r *quizSessionRepository) GetQuizSessionsByQuizID(quizID uint) ([]models.QuizSession, error) {
	var sessions []models.QuizSession
	err := r.db.Preload("Quiz").
		Preload("Participant.User").
		Preload("Question").
		Where("quiz_id = ?", quizID).
		Order("created_at ASC").
		Find(&sessions).Error
	return sessions, err
}

func (r *quizSessionRepository) GetQuizSessionsByParticipantID(participantID uint) ([]models.QuizSession, error) {
	var sessions []models.QuizSession
	err := r.db.Preload("Quiz").
		Preload("Participant.User").
		Preload("Question").
		Where("participant_id = ?", participantID).
		Order("created_at ASC").
		Find(&sessions).Error
	return sessions, err
}

func (r *quizSessionRepository) GetQuizSession(quizID, participantID, questionID uint) (*models.QuizSession, error) {
	var session models.QuizSession
	err := r.db.Preload("Quiz").
		Preload("Participant.User").
		Preload("Question").
		Where("quiz_id = ? AND participant_id = ? AND question_id = ?", quizID, participantID, questionID).
		First(&session).Error
	return &session, err
}

func (r *quizSessionRepository) UpdateQuizSession(session *models.QuizSession) (*models.QuizSession, error) {
	err := r.db.Save(session).Error
	if err != nil {
		return nil, err
	}

	// Load relationships
	err = r.db.Preload("Quiz").Preload("Participant.User").Preload("Question").First(session, session.ID).Error
	return session, err
}

func (r *quizSessionRepository) GetQuizSessionStats(quizID uint) (map[uint]models.Participant, error) {
	var participants []models.Participant
	err := r.db.Preload("User").Where("quiz_id = ?", quizID).Find(&participants).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uint]models.Participant)
	for _, participant := range participants {
		result[participant.UserID] = participant
	}

	return result, nil
}
