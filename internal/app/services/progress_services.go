package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type ProgressServiceInterface interface {
	CreateProgress(progress models.Progress) (*models.Progress, error)
	UpdateProgress(id uint, progress models.Progress) (*models.Progress, error)
	GetByUserAndLesson(userID uint, lessonID uint) (*models.Progress, error)
	GetProgressByID(id uint) (*models.Progress, error)
	GetAllProgress() ([]models.Progress, error)
}

type progressService struct {
	progressRepo repositories.ProgressRepositoryInterface
	userRepo     repositories.UserRepositoryInterface
	lessonRepo   repositories.LessonRepositoryInterface
}

func NewProgressService(
	progressRepo repositories.ProgressRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	lessonRepo repositories.LessonRepositoryInterface,
) *progressService {
	return &progressService{
		progressRepo: progressRepo,
		userRepo:     userRepo,
		lessonRepo:   lessonRepo,
	}
}

func (s *progressService) CreateProgress(progress models.Progress) (*models.Progress, error) {
	return s.progressRepo.CreateProgress(&progress)
}

func (s *progressService) UpdateProgress(id uint, progress models.Progress) (*models.Progress, error) {
	return s.progressRepo.UpdateProgress(id, &progress)
}

func (s *progressService) GetProgressByID(id uint) (*models.Progress, error) {
	return s.progressRepo.GetProgressById(id)
}

func (s *progressService) GetAllProgress() ([]models.Progress, error) {
	return s.progressRepo.GetAllProgress()
}

func (s *progressService) GetByUserAndLesson(userID uint, lessonID uint) (*models.Progress, error) {
	return s.progressRepo.GetByUserAndLesson(userID, lessonID)
}
