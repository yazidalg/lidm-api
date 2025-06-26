package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type CourseServiceInterface interface {
	CreateCourse(request request.CreateCourseRequest) (*models.Course, error)
	GetCourseByID(id uint32) (*models.Course, error)
	GetAllCourses() ([]models.Course, error)
	UpdateCourse(id uint32, course request.UpdateCourseRequest) (*models.Course, error)
	DeleteCourse(id uint32) error
}

type courseService struct {
	courseRepository repositories.CourseRepositoryInterface
}

func NewCourseService(courseRepository repositories.CourseRepositoryInterface) *courseService {
	return &courseService{courseRepository}
}

func (s *courseService) CreateCourse(request request.CreateCourseRequest) (*models.Course, error) {
	// Convert request to model
	courseData := models.Course{
		Title:       request.Title,
		Description: request.Description,
		Thumbnail:   request.Thumbnail,
	}

	// Delegate to repository
	return s.courseRepository.CreateCourse(&courseData)
}

func (s *courseService) GetCourseByID(id uint32) (*models.Course, error) {
	return s.courseRepository.GetCourseByID(id)
}

func (s *courseService) GetAllCourses() ([]models.Course, error) {
	return s.courseRepository.GetAllCourses()
}

func (s *courseService) UpdateCourse(id uint32, course request.UpdateCourseRequest) (*models.Course, error) {
	courseData := models.Course{
		Title:       course.Title,
		Description: course.Description,
		Thumbnail:   course.Thumbnail,
	}

	return s.courseRepository.UpdateCourse(id, &courseData)
}

func (s *courseService) DeleteCourse(id uint32) error {
	return s.courseRepository.DeleteCourse(id)
}
