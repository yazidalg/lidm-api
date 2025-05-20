package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type QuestionServiceInterface interface {
	CreateQuestion(request request.CreateQuestionRequest) (*models.Question, error)
	GetQuestionByID(id uint) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	UpdateQuestion(id int, question *models.Question) (*models.Question, error)
	DeleteQuestion(id uint) error
}

type questionService struct {
	questionRepository repositories.QuestionRepositoryInterface
}

func NewQuestionService(questionRepository repositories.QuestionRepositoryInterface) *questionService {
	return &questionService{questionRepository}
}

func (s *questionService) CreateQuestion(request request.CreateQuestionRequest) (*models.Question, error) {

	questionData := models.Question{
		Question:      request.Question,
		AnswerTime:    request.AnswerTime,
		ReadTime:      request.ReadTime,
		CorrectAnswer: request.CorrectAnswer,
		Explanation:   request.Explanation,
		Options:       models.Options(request.Options),
	}

	return s.questionRepository.CreateQuestion(&questionData)
}

func (s *questionService) GetQuestionByID(id uint) (*models.Question, error) {
	return s.questionRepository.GetQuestionByID(id)
}

func (s *questionService) GetAllQuestions() ([]models.Question, error) {
	return s.questionRepository.GetAllQuestions()
}

func (s *questionService) UpdateQuestion(id int, question *models.Question) (*models.Question, error) {
	return s.questionRepository.UpdateQuestion(id, question)
}

func (s *questionService) DeleteQuestion(id uint) error {
	return s.questionRepository.DeleteQuestion(id)
}
