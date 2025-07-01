package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type PrequizServiceInterface interface {
	CreatePrequiz(prequiz request.PrequizRequest) (*models.Prequiz, error)
	GetPrequizByID(id uint) (*models.Prequiz, error)
	GetAllPrequizzes() ([]models.Prequiz, error)
	UpdatePrequiz(id uint, prequiz request.PrequizRequest) (*models.Prequiz, error)
}

type prequizService struct {
	prequizRepo repositories.PrequizRepositoryInterface
	lessonRepo  repositories.LessonRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
}

func NewPrequizService(
	prequizRepo repositories.PrequizRepositoryInterface,
	lessonRepo repositories.LessonRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *prequizService {
	return &prequizService{
		prequizRepo: prequizRepo,
		lessonRepo:  lessonRepo,
		userRepo:    userRepo,
	}
}

func (s *prequizService) CreatePrequiz(prequiz request.PrequizRequest) (*models.Prequiz, error) {

	prequizModel := models.Prequiz{
		LessonID:      prequiz.LessonID,
		UserID:        prequiz.UserID,
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
		LessonID:      prequiz.LessonID,
		UserID:        prequiz.UserID,
		Question:      prequiz.Question,
		CorrectAnswer: prequiz.CorrectAnswer,
		Explanation:   prequiz.Explanation,
		Options:       models.Options(prequiz.Options),
	}

	return s.prequizRepo.UpdatePrquiz(id, &prequizModel)
}
