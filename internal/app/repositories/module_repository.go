package repositories

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ModuleRepositoryInterface interface {
	CreateModule(module *models.Module) (*models.Module, error)
	GetModuleByID(id uint32) (*models.Module, error)
	GetAllModules() ([]models.Module, error)
	UpdateModule(id uint32, module *models.Module) (*models.Module, error)
	UpdateModuleWithVideo(id uint32, module *models.Module) (*models.Module, error)
	DeleteModule(id uint32) error
	CreateARExperiment(arExperiment *models.ARExperiment) (*models.ARExperiment, error)
	GetFlashcardsByModule(moduleID uint) ([]models.Flashcard, error)
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
	if module.VideoMaterial != nil {
		// Set the ModuleID
		module.VideoMaterial.ModuleID = module.ID

		// Ensure VideoMaterial doesn't have a pre-set ID to avoid conflicts
		module.VideoMaterial.ID = 0

		// Reset timestamps to let GORM handle them
		module.VideoMaterial.CreatedAt = time.Time{}
		module.VideoMaterial.UpdatedAt = time.Time{}
		module.VideoMaterial.DeletedAt = gorm.DeletedAt{}

		// Handle VideoQuizzes if provided
		if len(module.VideoMaterial.VideoQuizzes) > 0 {
			for i := range module.VideoMaterial.VideoQuizzes {
				// Reset ID and timestamps for each quiz
				module.VideoMaterial.VideoQuizzes[i].ID = 0
				module.VideoMaterial.VideoQuizzes[i].CreatedAt = time.Time{}
				module.VideoMaterial.VideoQuizzes[i].UpdatedAt = time.Time{}
				module.VideoMaterial.VideoQuizzes[i].DeletedAt = gorm.DeletedAt{}
			}
		}

		// Create the video material
		if err := tx.Create(module.VideoMaterial).Error; err != nil {
			tx.Rollback()
			return nil, err
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

	// Preload semua komponen yang diperlukan untuk satu module: Video, AR, Prequizzes, dan Flashcards
	err := r.db.
		Preload("VideoMaterial").                                          // Video content
		Preload("VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("ARExperiment").                           // AR experiments
		Preload("Prequizzes", func(db *gorm.DB) *gorm.DB { // Prequizzes ordered by created_at
			return db.Order("prequizzes.created_at ASC")
		}).
		Preload("Flashcards", func(db *gorm.DB) *gorm.DB { // Flashcards ordered by order
			return db.Order("flashcards.`order` ASC")
		}).
		First(&module, id).Error

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) GetAllModules() ([]models.Module, error) {
	var modules []models.Module

	// Preload komponen utama: Video, AR, Prequizzes, dan Flashcards
	err := r.db.
		Preload("VideoMaterial").                                          // Video content
		Preload("VideoMaterial.VideoQuizzes", func(db *gorm.DB) *gorm.DB { // Video quizzes ordered by timestamp
			return db.Order("video_quizzes.timestamp_start ASC")
		}).
		Preload("ARExperiment").                           // AR experiments
		Preload("Prequizzes", func(db *gorm.DB) *gorm.DB { // Load all prequizzes ordered by created_at
			return db.Order("prequizzes.created_at ASC")
		}).
		Preload("Flashcards", func(db *gorm.DB) *gorm.DB { // Load flashcards ordered by order
			return db.Order("flashcards.`order` ASC")
		}).
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
	if module.VideoMaterial != nil {
		// Set the ModuleID
		module.VideoMaterial.ModuleID = uint(id)

		// Create the video material
		if err := tx.Create(module.VideoMaterial).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Return the updated module with preloaded relationships
	return r.GetModuleByID(id)
}

func (r *moduleRepository) UpdateModuleWithVideo(id uint32, module *models.Module) (*models.Module, error) {
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

	// Update the module using Model and Updates to avoid timestamp issues
	if err := tx.Model(&existingModule).Updates(map[string]interface{}{
		"title":       existingModule.Title,
		"description": existingModule.Description,
		"offset_x":    existingModule.OffsetX,
		"offset_y":    existingModule.OffsetY,
		"icon":        existingModule.Icon,
		"thumbnail":   existingModule.Thumbnail,
		"updated_at":  time.Now(), // Set current timestamp
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle VideoMaterial creation if provided
	if module.VideoMaterial != nil {
		// Reset ID and timestamps to let GORM auto-generate
		module.VideoMaterial.ID = 0
		module.VideoMaterial.CreatedAt = time.Time{}
		module.VideoMaterial.UpdatedAt = time.Time{}
		module.VideoMaterial.DeletedAt = gorm.DeletedAt{}

		// Set the ModuleID
		module.VideoMaterial.ModuleID = uint(id)

		// Create the video material
		if err := tx.Create(module.VideoMaterial).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Handle ARExperiment creation if provided
	if module.ARExperiment != nil {
		// Reset ID and timestamps to let GORM auto-generate
		module.ARExperiment.ID = 0
		module.ARExperiment.CreatedAt = time.Time{}
		module.ARExperiment.UpdatedAt = time.Time{}
		module.ARExperiment.DeletedAt = gorm.DeletedAt{}

		// Set the ModuleID
		module.ARExperiment.ModuleID = uint(id)

		// Create the AR experiment
		if err := tx.Create(module.ARExperiment).Error; err != nil {
			tx.Rollback()
			return nil, err
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

func (r *moduleRepository) CreateARExperiment(arExperiment *models.ARExperiment) (*models.ARExperiment, error) {
	// Use a transaction to ensure data integrity
	tx := r.db.Begin()

	// Create the ARExperiment
	if err := tx.Create(arExperiment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return arExperiment, nil
}

func (r *moduleRepository) GetFlashcardsByModule(moduleID uint) ([]models.Flashcard, error) {
	var flashcards []models.Flashcard
	err := r.db.Where("module_id = ?", moduleID).Order("`order` ASC").Find(&flashcards).Error
	return flashcards, err
}
