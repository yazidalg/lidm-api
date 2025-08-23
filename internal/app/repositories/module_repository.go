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

	// Create the module first
	if err := tx.Create(module).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle VideoMaterial creation if provided
	if len(module.VideoMaterial) > 0 {
		for i := range module.VideoMaterial {
			// Set the ModuleID
			module.VideoMaterial[i].ModuleID = module.ID

			// Create the video material
			if err := tx.Create(&module.VideoMaterial[i]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return module, nil
}

func (r *moduleRepository) GetModuleByID(id uint32) (*models.Module, error) {
	var module models.Module

	// Preload komponen utama: Video, AR, dan Prequizzes
	err := r.db.
		Preload("VideoMaterial").                                          // Video content
		Preload("VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("ARExperiment"). // AR experiments
		// Preload("Prequizzes").                                      // Prequizzes - DISABLED temporarily
		First(&module, id).Error

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) GetAllModules() ([]models.Module, error) {
	var modules []models.Module

	// Preload komponen utama: Video, AR, dan Prequizzes
	err := r.db.
		Preload("VideoMaterial").                                          // Video content
		Preload("VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("ARExperiment"). // AR experiments
		// Preload("Prequizzes").                                      // Prequizzes - Disabled again to debug
		Order("created_at ASC"). // Order modules by creation
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
	existingModule.OffsetX = module.OffsetX
	existingModule.OffsetY = module.OffsetY
	existingModule.Icon = module.Icon
	existingModule.Thumbnail = module.Thumbnail

	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	// Update the module
	if err := tx.Save(&existingModule).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle VideoMaterial creation if provided
	if len(module.VideoMaterial) > 0 {
		for _, videoMat := range module.VideoMaterial {
			// Set the ModuleID
			videoMat.ModuleID = uint(id)

			// Create the video material
			if err := tx.Create(&videoMat).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return the updated module with preloaded relationships
	return r.GetModuleByID(id)
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
