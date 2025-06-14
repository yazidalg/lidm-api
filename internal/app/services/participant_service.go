package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type ParticipantServiceInterface interface {
	CreateParticipant(request request.CreateParticipantRequest) (*models.Participant, error)
	GetParticipantByID(id int32) (*models.Participant, error)
	GetAllParticipants() ([]models.Participant, error)
	GetParticipantsByQuizID(quizID uint) ([]models.Participant, error)
	GetParticipantsByUserID(userID uint) ([]models.Participant, error)
	UpdateParticipant(id int32, request request.UpdateParticipantRequest) (*models.Participant, error)
	DeleteParticipant(id int32) error
}

type participantService struct {
	participantRepository repositories.ParticipantRepositoryInterface
}

func NewParticipantService(participantRepository repositories.ParticipantRepositoryInterface) *participantService {
	return &participantService{participantRepository}
}

func (s *participantService) CreateParticipant(request request.CreateParticipantRequest) (*models.Participant, error) {
	// Convert request to model
	participantData := models.Participant{
		UserID:     request.UserID,
		QuizID:     request.QuizID,
		TotalScore: request.TotalScore,
	}

	// Delegate to repository
	return s.participantRepository.CreateParticipant(&participantData)
}

func (s *participantService) GetParticipantByID(id int32) (*models.Participant, error) {
	return s.participantRepository.GetParticipantByID(id)
}

func (s *participantService) GetAllParticipants() ([]models.Participant, error) {
	return s.participantRepository.GetAllParticipants()
}

func (s *participantService) GetParticipantsByQuizID(quizID uint) ([]models.Participant, error) {
	return s.participantRepository.GetParticipantsByQuizID(quizID)
}

func (s *participantService) GetParticipantsByUserID(userID uint) ([]models.Participant, error) {
	return s.participantRepository.GetParticipantsByUserID(userID)
}

func (s *participantService) UpdateParticipant(id int32, request request.UpdateParticipantRequest) (*models.Participant, error) {
	// Get the existing participant
	existingParticipant, err := s.participantRepository.GetParticipantByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided in the request
	if request.UserID != 0 {
		existingParticipant.UserID = request.UserID
	}

	if request.QuizID != 0 {
		existingParticipant.QuizID = request.QuizID
	}

	if request.TotalScore != 0 {
		existingParticipant.TotalScore = request.TotalScore
	}

	// Delegate to repository
	return s.participantRepository.UpdateParticipant(id, existingParticipant)
}

func (s *participantService) DeleteParticipant(id int32) error {
	return s.participantRepository.DeleteParticipant(id)
}
