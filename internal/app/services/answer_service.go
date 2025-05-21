package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"gorm.io/gorm"
)

type AnswerServiceInterface interface {
	CreateAnswer(request request.CreateAnswerRequest) (*models.Answer, error)
	GetAnswerByID(id int32) (*models.Answer, error)
	GetAllAnswers() ([]models.Answer, error)
	UpdateAnswer(id int32, answer *models.Answer) (*models.Answer, error)
	DeleteAnswer(id int32) error
}

type answerService struct {
	answerRepository repositories.AnswerRepositoryInterface
}

func NewAnswerService(answerRepository repositories.AnswerRepositoryInterface) *answerService {
	return &answerService{answerRepository}
}

func (s *answerService) CreateAnswer(request request.CreateAnswerRequest) (*models.Answer, error) {
	// Convert request to model
	answerData := models.Answer{
		QuestionID:     request.QuestionID,
		UserID:         request.UserID,
		OptionSelected: request.OptionSelected,
		IsCorrect:      request.IsCorrect,
		Score:          request.Score,
	}

	// Add quiz associations
	for _, quizID := range request.QuizID {
		answerData.QuizID = append(answerData.QuizID, models.Quiz{Model: gorm.Model{ID: quizID}})
	}

	// Delegate to repository
	return s.answerRepository.CreateAnswer(&answerData)
}

func (s *answerService) GetAnswerByID(id int32) (*models.Answer, error) {
	return s.answerRepository.GetAnswerByID(id)
}

func (s *answerService) GetAllAnswers() ([]models.Answer, error) {
	return s.answerRepository.GetAllAnswers()
}

func (s *answerService) UpdateAnswer(id int32, answer *models.Answer) (*models.Answer, error) {
	return s.answerRepository.UpdateAnswer(id, answer)
}

func (s *answerService) DeleteAnswer(id int32) error {
	return s.answerRepository.DeleteAnswer(id)
}
