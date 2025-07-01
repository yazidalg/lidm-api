package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type ProgressServiceInterface interface {
	UpdateProgress(progress request.ProgressRequest) (*models.Progress, error)
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

func (s *progressService) UpdateProgress(progress request.ProgressRequest) (*models.Progress, error) {
	user, err := s.userRepo.GetUserById(int(progress.UserID))
	if err != nil || user.ID == 0 {
		return nil, err
	}

	lesson, err := s.lessonRepo.GetLessonByID(progress.LessonID)
	if err != nil || lesson == nil {
		return nil, err
	}

	progressData := &models.Progress{
		UserID:    uint(user.ID),
		LessonID:  uint(progress.LessonID),
		Completed: true, // Assuming the progress is marked as completed
	}

	return s.progressRepo.UpdateProgress(progressData)
}

func (s *progressService) GetAllProgress() ([]models.Progress, error) {
	return s.progressRepo.GetAllProgress()
}
