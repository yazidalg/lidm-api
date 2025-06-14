package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type QuestionServiceInterface interface {
	CreateQuestion(request request.CreateQuestionRequest) (*models.Question, error)
	GetQuestionByID(id int32) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	UpdateQuestion(id int32, request request.UpdateQuestionRequest) (*models.Question, error)
	DeleteQuestion(id int32) error
	GetRandomQuestion(count int) (*[]models.Question, error)
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

func (s *questionService) GetQuestionByID(id int32) (*models.Question, error) {
	return s.questionRepository.GetQuestionByID(id)
}

func (s *questionService) GetAllQuestions() ([]models.Question, error) {
	return s.questionRepository.GetAllQuestions()
}

func (s *questionService) UpdateQuestion(id int32, request request.UpdateQuestionRequest) (*models.Question, error) {
	questionData := models.Question{
		Question:      request.Question,
		AnswerTime:    request.AnswerTime,
		ReadTime:      request.ReadTime,
		CorrectAnswer: request.CorrectAnswer,
		Explanation:   request.Explanation,
		Options:       models.Options(request.Options),
	}

	return s.questionRepository.UpdateQuestion(id, &questionData)
}

func (s *questionService) DeleteQuestion(id int32) error {
	return s.questionRepository.DeleteQuestion(id)
}

func (s *questionService) GetRandomQuestion(count int) (*[]models.Question, error) {
	return s.questionRepository.GetRandomQuestion(count)
}
