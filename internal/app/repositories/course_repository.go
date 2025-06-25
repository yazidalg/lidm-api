package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type CourseRepositoryInterface interface {
	CreateCourse(course *models.Course) (*models.Course, error)
	GetCourseByID(id uint32) (*models.Course, error)
	GetAllCourses() ([]models.Course, error)
	UpdateCourse(id uint32, course *models.Course) (*models.Course, error)
	DeleteCourse(id uint32) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *courseRepository {
	return &courseRepository{db}
}

func (r *courseRepository) CreateCourse(course *models.Course) (*models.Course, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Create(course).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return course, nil
}

func (r *courseRepository) GetCourseByID(id uint32) (*models.Course, error) {
	var course models.Course

	// Preload related entities if necessary
	err := r.db.First(&course, id).Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) GetAllCourses() ([]models.Course, error) {
	var courses []models.Course

	// Preload related entities if necessary
	err := r.db.Find(&courses).Error

	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *courseRepository) UpdateCourse(id uint32, course *models.Course) (*models.Course, error) {
	var existingCourse models.Course

	// Find the existing course
	if err := r.db.First(&existingCourse, id).Error; err != nil {
		return nil, err
	}

	// Update the course fields
	existingCourse.Title = course.Title
	existingCourse.Description = course.Description
	existingCourse.Thumbnail = course.Thumbnail

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Save(&existingCourse).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &existingCourse, nil
}

func (r *courseRepository) DeleteCourse(id uint32) error {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	// Delete the course
	if err := tx.Delete(&models.Course{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
