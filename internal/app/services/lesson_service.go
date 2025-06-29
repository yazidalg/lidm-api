package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type LessonServiceInterface interface {
	CreateLesson(lesson request.LessonRequest) (*models.Lesson, error)
	GetLessonByID(id uint32) (*models.Lesson, error)
	GetAllLessons() ([]models.Lesson, error)
	UpdateLesson(id uint32, lesson request.LessonRequest) (*models.Lesson, error)
	DeleteLesson(id uint32) error
}

type lessonService struct {
	lessonRepo repositories.LessonRepositoryInterface
}

func NewLessonService(lessonRepo repositories.LessonRepositoryInterface) *lessonService {
	return &lessonService{lessonRepo: lessonRepo}
}

func (s *lessonService) CreateLesson(lesson request.LessonRequest) (*models.Lesson, error) {
	lessonModel := &models.Lesson{
		ModuleID:  lesson.ModuleID,
		Title:     lesson.Title,
		Content:   lesson.Content,
		SortOrder: lesson.SortOrder,
	}

	return s.lessonRepo.CreateLesson(lessonModel)
}

func (s *lessonService) GetLessonByID(id uint32) (*models.Lesson, error) {
	return s.lessonRepo.GetLessonByID(id)
}

func (s *lessonService) GetAllLessons() ([]models.Lesson, error) {
	return s.lessonRepo.GetAllLessons()
}

func (s *lessonService) UpdateLesson(id uint32, lesson request.LessonRequest) (*models.Lesson, error) {
	lessonModel := &models.Lesson{
		ModuleID:  lesson.ModuleID,
		Title:     lesson.Title,
		Content:   lesson.Content,
		SortOrder: lesson.SortOrder,
	}

	return s.lessonRepo.UpdateLesson(id, lessonModel)
}

func (s *lessonService) DeleteLesson(id uint32) error {
	return s.lessonRepo.DeleteLesson(id)
}
