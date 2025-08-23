package services

import (
	"errors"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type PrequizServiceInterface interface {
	CreatePrequiz(prequiz request.PrequizRequest) (*models.Prequiz, error)
	GetPrequizByID(id uint) (*models.Prequiz, error)
	GetAllPrequizzes() ([]models.Prequiz, error)
	UpdatePrequiz(id uint, prequiz request.PrequizRequest) (*models.Prequiz, error)
	GetUserPrequizAnswers(userID uint) ([]models.PrequizUserAnswer, error)
	SubmitPrequizAnswer(userID uint, answer request.PrequizAnswerRequest) (*models.PrequizUserAnswer, error)
	GetPrequizzesByModule(moduleID uint) ([]models.Prequiz, error)
}

type prequizService struct {
	prequizRepo repositories.PrequizRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
}

func NewPrequizService(
	prequizRepo repositories.PrequizRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *prequizService {
	return &prequizService{
		prequizRepo: prequizRepo,
		userRepo:    userRepo,
	}
}

func (s *prequizService) CreatePrequiz(prequiz request.PrequizRequest) (*models.Prequiz, error) {

	prequizModel := models.Prequiz{
		ModuleID:      prequiz.ModuleID,
		Question:      prequiz.Question,
		CorrectAnswer: prequiz.CorrectAnswer,
		Explanation:   prequiz.Explanation,
		Options:       models.Options(prequiz.Options),
	}

	return s.prequizRepo.CreatePrequiz(&prequizModel)
}

func (s *prequizService) GetPrequizByID(id uint) (*models.Prequiz, error) {
	return s.prequizRepo.GetPrequizByID(id)
}

func (s *prequizService) GetAllPrequizzes() ([]models.Prequiz, error) {
	return s.prequizRepo.GetAllPrequizzes()
}

func (s *prequizService) UpdatePrequiz(id uint, prequiz request.PrequizRequest) (*models.Prequiz, error) {

	prequizModel := models.Prequiz{
		ModuleID:      prequiz.ModuleID,
		Question:      prequiz.Question,
		CorrectAnswer: prequiz.CorrectAnswer,
		Explanation:   prequiz.Explanation,
		Options:       models.Options(prequiz.Options),
	}

	return s.prequizRepo.UpdatePrquiz(id, &prequizModel)
}

func (s *prequizService) GetUserPrequizAnswers(userID uint) ([]models.PrequizUserAnswer, error) {
	return s.prequizRepo.GetUserPrequizAnswers(userID)
}

func (s *prequizService) SubmitPrequizAnswer(userID uint, answer request.PrequizAnswerRequest) (*models.PrequizUserAnswer, error) {
	// Check if user already answered this prequiz
	hasAnswered, err := s.prequizRepo.HasUserAnsweredPrequiz(userID, answer.PrequizID)
	if err != nil {
		return nil, err
	}

	if hasAnswered {
		return nil, errors.New("user has already answered this prequiz")
	}

	// Get the prequiz to check correct answer
	prequiz, err := s.prequizRepo.GetPrequizByID(answer.PrequizID)
	if err != nil {
		return nil, err
	}

	// Check if answer is correct
	isCorrect := prequiz.CorrectAnswer == answer.SelectedAnswer

	userAnswer := models.PrequizUserAnswer{
		PrequizID:  answer.PrequizID,
		UserID:     userID,
		Answer:     answer.SelectedAnswer,
		IsCorrect:  isCorrect,
		AnsweredAt: time.Now().Unix(),
	}

	return s.prequizRepo.SubmitPrequizAnswer(&userAnswer)
}

func (s *prequizService) GetPrequizzesByModule(moduleID uint) ([]models.Prequiz, error) {
	return s.prequizRepo.GetPrequizzesByModule(moduleID, 0) // 0 means no limit
}
