package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type LessonRepositoryInterface interface {
	CreateLesson(lesson *models.Lesson) (*models.Lesson, error)
	GetLessonByID(id uint32) (*models.Lesson, error)
	GetAllLessons() ([]models.Lesson, error)
	UpdateLesson(id uint32, lesson *models.Lesson) (*models.Lesson, error)
	DeleteLesson(id uint32) error
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *lessonRepository {
	return &lessonRepository{db}
}

func (r *lessonRepository) CreateLesson(lesson *models.Lesson) (*models.Lesson, error) {
	tx := r.db.Begin()

	if err := tx.Create(lesson).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *lessonRepository) GetLessonByID(id uint32) (*models.Lesson, error) {
	var lesson models.Lesson

	// Preload related entities if necessary
	err := r.db.Preload("Module").First(&lesson, id).Error

	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

func (r *lessonRepository) GetAllLessons() ([]models.Lesson, error) {
	var lessons []models.Lesson

	// Preload related entities if necessary
	err := r.db.Preload("Module").Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *lessonRepository) UpdateLesson(id uint32, lesson *models.Lesson) (*models.Lesson, error) {
	existingLesson, err := r.GetLessonByID(id)

	if err != nil {
		return nil, err
	}

	// Update the lesson fields
	existingLesson.Title = lesson.Title
	existingLesson.Content = lesson.Content
	existingLesson.SortOrder = lesson.SortOrder

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Save(&existingLesson).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return existingLesson, nil
}

func (r *lessonRepository) DeleteLesson(id uint32) error {
	var lesson models.Lesson

	// Find the lesson by ID
	if err := r.db.First(&lesson, id).Error; err != nil {
		return err
	}

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Delete(&lesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
