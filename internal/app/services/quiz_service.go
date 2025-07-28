package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"gorm.io/gorm"
)

type QuizServiceInterface interface {
	CreateQuiz(request request.CreateQuizRequest) (*models.Quiz, error)
	GetQuizByID(id uint) (*models.Quiz, error)
	GetAllQuizzes() ([]models.Quiz, error)
	UpdateQuiz(id uint, request request.UpdateQuizRequest) (*models.Quiz, error)
	GetQuizByInviteCode(inviteCode string) (*models.Quiz, error)
	DeleteQuiz(id uint) error
}

type quizService struct {
	quizRepository repositories.QuizRepositoryInterface
}

func NewQuizService(quizRepository repositories.QuizRepositoryInterface) *quizService {
	return &quizService{quizRepository}
}

// GetQuizByInviteCode implements QuizServiceInterface.
func (s *quizService) GetQuizByInviteCode(code string) (*models.Quiz, error) {
	return s.quizRepository.GetQuizByInviteCode(code)
}

func (s *quizService) CreateQuiz(request request.CreateQuizRequest) (*models.Quiz, error) {
	// Convert request to model
	quizData := models.Quiz{
		Status:        request.Status,
		Mode:          request.Mode,
		ModuleID:      &request.ModuleID,
		HostUserID:    request.HostUserID,
		InviteCode:    request.InviteCode,
		QuestionCount: request.QuestionCount,
	}

	// Set default question count if not provided
	if quizData.QuestionCount == 0 {
		quizData.QuestionCount = 5
	}

	// Get random questions from the module if not provided
	if len(request.QuestionsIDs) == 0 {
		// Let the repository handle getting random questions
		quizData.Questions = []models.Question{} // Will be populated by repository
	} else {
		// Convert question IDs to Question pointers
		questions := make([]models.Question, len(request.QuestionsIDs))
		for i, id := range request.QuestionsIDs {
			questions[i] = models.Question{
				Model: gorm.Model{
					ID: id,
				},
			}
		}
		quizData.Questions = questions
	}

	// Set default status if not provided
	if quizData.Status == "" {
		quizData.Status = "pending"
	}

	// Delegate to repository
	return s.quizRepository.CreateQuiz(&quizData)
}

func (s *quizService) GetQuizByID(id uint) (*models.Quiz, error) {
	return s.quizRepository.GetQuizByID(id)
}

func (s *quizService) GetAllQuizzes() ([]models.Quiz, error) {
	return s.quizRepository.GetAllQuizzes()
}

func (s *quizService) UpdateQuiz(id uint, request request.UpdateQuizRequest) (*models.Quiz, error) {
	// Initialize quiz model for update
	quizData := models.Quiz{}

	if request.Status != "" {
		quizData.Status = request.Status
	}

	if request.WinnerID != nil {
		quizData.WinnerID = request.WinnerID
	}

	// Convert question IDs to Question pointers if provided
	if len(request.QuestionsIDs) > 0 {
		questions := make([]models.Question, len(request.QuestionsIDs))
		for i, id := range request.QuestionsIDs {
			questions[i] = models.Question{
				Model: gorm.Model{
					ID: id,
				},
			}
		}
		quizData.Questions = questions
	}

	// Delegate to repository
	return s.quizRepository.UpdateQuiz(id, &quizData)
}

func (s *quizService) DeleteQuiz(id uint) error {
	return s.quizRepository.DeleteQuiz(id)
}
