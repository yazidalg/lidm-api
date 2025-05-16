package services

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type QuizServiceInterface interface {
	GetQuizByID(id int) (*models.Quiz, error)
	GetAllQuizzes() (*[]models.Quiz, error)
	CreateQuiz(quiz request.QuizRequest) (*models.Quiz, error)
	UpdateQuiz(id int32, quiz models.Quiz) (*models.Quiz, error)
	DeleteQuiz(id int32, quiz models.Quiz) (*models.Quiz, error)
}

type quizService struct {
	quizRepository repositories.QuizRepositoryInterface
}

func NewQuizService(quizRepository repositories.QuizRepositoryInterface) *quizService {
	return &quizService{quizRepository}
}

func (s *quizService) GetQuizByID(id int) (*models.Quiz, error) {
	return s.quizRepository.GetQuizByID(id)
}

func (s *quizService) GetAllQuizzes() (*[]models.Quiz, error) {
	return s.quizRepository.GetAllQuizzes()
}

func (s *quizService) CreateQuiz(quiz request.QuizRequest) (*models.Quiz, error) {

	quizData := models.Quiz{
		Question:      quiz.Question,
		AnswerTime:    quiz.AnswerTime,
		ReadTime:      quiz.ReadTime,
		Options:       models.Options(quiz.Options),
		CorrectAnswer: quiz.CorrectAnswer,
		Explanation:   quiz.Explanation,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.quizRepository.CreateQuiz(quizData)
}

func (s *quizService) UpdateQuiz(id int32, quiz models.Quiz) (*models.Quiz, error) {

	quizModel := models.Quiz{
		Question:      quiz.Question,
		AnswerTime:    quiz.AnswerTime,
		ReadTime:      quiz.ReadTime,
		Options:       models.Options(quiz.Options),
		CorrectAnswer: quiz.CorrectAnswer,
		Explanation:   quiz.Explanation,
		UpdatedAt:     time.Now(),
	}

	return s.quizRepository.UpdateQuiz(id, quizModel)
}

func (s *quizService) DeleteQuiz(id int32, quiz models.Quiz) (*models.Quiz, error) {
	return s.quizRepository.DeleteQuiz(id, quiz)
}
