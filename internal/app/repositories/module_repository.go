package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ModuleRepositoryInterface interface {
	CreateModule(module *models.Module) (*models.Module, error)
	GetModuleByID(id uint32) (*models.Module, error)
	GetAllModules() ([]models.Module, error)
	UpdateModule(id uint32, module *models.Module) (*models.Module, error)
	DeleteModule(id uint32) error
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *moduleRepository {
	return &moduleRepository{db}
}

func (r *moduleRepository) CreateModule(module *models.Module) (*models.Module, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Create(module).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return module, nil
}

func (r *moduleRepository) GetModuleByID(id uint32) (*models.Module, error) {
	var module models.Module

	// Preload SubMaterials dengan komponen utama: Video, AR, dan Prequizzes (tanpa Module relationship)
	err := r.db.
		Preload("Lessons").                                                      // Legacy support
		Preload("SubMaterials", func(db *gorm.DB) *gorm.DB {                     // SubMaterials ordered by order field
			return db.Order("sub_materials.order ASC")
		}).
		Preload("SubMaterials.VideoMaterial").                                   // Video content
		Preload("SubMaterials.VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("SubMaterials.ARExperiment").                                    // AR experiments
		// Preload("SubMaterials.Prequizzes").                                      // Prequizzes - DISABLED temporarily
		First(&module, id).Error

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) GetAllModules() ([]models.Module, error) {
	var modules []models.Module

	// Preload SubMaterials dengan komponen utama: Video, AR, dan Prequizzes (tanpa Module relationship)
	err := r.db.
		Preload("Lessons").                                                      // Legacy support
		Preload("SubMaterials", func(db *gorm.DB) *gorm.DB {                     // SubMaterials ordered by order field
			return db.Order("sub_materials.order ASC")
		}).
		Preload("SubMaterials.VideoMaterial").                                   // Video content
		Preload("SubMaterials.VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("SubMaterials.ARExperiment").                                    // AR experiments
		// Preload("SubMaterials.Prequizzes").                                      // Prequizzes - Disabled again to debug
		Order("created_at ASC").                                                 // Order modules by creation
		Find(&modules).Error

	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *moduleRepository) UpdateModule(id uint32, module *models.Module) (*models.Module, error) {
	existingModule, err := r.GetModuleByID(id)

	if err != nil {
		return nil, err
	}

	// Update the module fields
	existingModule.Title = module.Title
	existingModule.Description = module.Description

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Save(&existingModule).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return existingModule, nil
}

func (r *moduleRepository) DeleteModule(id uint32) error {
	var module models.Module

	// Check if the module exists
	if err := r.db.First(&module, id).Error; err != nil {
		return err
	}

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	if err := tx.Delete(&module).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
