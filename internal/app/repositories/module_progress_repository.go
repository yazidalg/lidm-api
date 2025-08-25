package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ModuleProgressRepositoryInterface interface {
	CreateModuleProgress(progress *models.ModuleProgress) (*models.ModuleProgress, error)
	GetByUserAndModule(userID, moduleID uint) (*models.ModuleProgress, error)
	GetAllByUser(userID uint) ([]models.ModuleProgress, error)
	UpdateProgress(progress *models.ModuleProgress) (*models.ModuleProgress, error)
	UnlockModule(userID, moduleID uint) error
	InitializeFirstModuleForUser(userID uint) error
	GetUnlockedModulesForUser(userID uint) ([]models.ModuleProgress, error)
}

type moduleProgressRepository struct {
	db *gorm.DB
}

func NewModuleProgressRepository(db *gorm.DB) *moduleProgressRepository {
	return &moduleProgressRepository{db: db}
}

func (r *moduleProgressRepository) CreateModuleProgress(progress *models.ModuleProgress) (*models.ModuleProgress, error) {
	tx := r.db.Begin()

	if err := tx.Create(progress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *moduleProgressRepository) GetByUserAndModule(userID, moduleID uint) (*models.ModuleProgress, error) {
	var progress models.ModuleProgress
	err := r.db.Where("user_id = ? AND module_id = ?", userID, moduleID).
		Preload("User").
		Preload("Module").
		First(&progress).Error

	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *moduleProgressRepository) GetAllByUser(userID uint) ([]models.ModuleProgress, error) {
	var progresses []models.ModuleProgress
	err := r.db.Where("user_id = ?", userID).
		Preload("Module").
		Order("module_id ASC").
		Find(&progresses).Error

	return progresses, err
}

func (r *moduleProgressRepository) UpdateProgress(progress *models.ModuleProgress) (*models.ModuleProgress, error) {
	tx := r.db.Begin()

	// Calculate new progress
	progress.Progress = progress.CalculateProgress(r.db)

	// Check if module should be marked as complete
	if progress.Progress >= 100 && !progress.IsComplete {
		progress.MarkAsCompleted()
	}

	if err := tx.Save(progress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Unlock next module if this one is completed
	if progress.IsComplete {
		if err := progress.CheckAndUnlockNextModule(r.db); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *moduleProgressRepository) UnlockModule(userID, moduleID uint) error {
	progress, err := r.GetByUserAndModule(userID, moduleID)
	if err == gorm.ErrRecordNotFound {
		// Create new progress entry
		progress = &models.ModuleProgress{
			UserID:     userID,
			ModuleID:   moduleID,
			IsUnlocked: true,
			IsComplete: false,
			Progress:   0,
		}
		progress.MarkAsStarted()
		_, err = r.CreateModuleProgress(progress)
		return err
	} else if err != nil {
		return err
	}

	if !progress.IsUnlocked {
		progress.IsUnlocked = true
		progress.MarkAsStarted()
		_, err = r.UpdateProgress(progress)
	}

	return err
}

func (r *moduleProgressRepository) InitializeFirstModuleForUser(userID uint) error {
	// Check if user already has any module progress
	var count int64
	r.db.Model(&models.ModuleProgress{}).Where("user_id = ?", userID).Count(&count)
	
	if count > 0 {
		return nil // User already has progress, don't initialize
	}

	// Find the first module (lowest ID)
	var firstModule models.Module
	err := r.db.Order("id ASC").First(&firstModule).Error
	if err != nil {
		return err
	}

	// Unlock the first module for this user
	return r.UnlockModule(userID, firstModule.ID)
}

func (r *moduleProgressRepository) GetUnlockedModulesForUser(userID uint) ([]models.ModuleProgress, error) {
	var progresses []models.ModuleProgress
	err := r.db.Where("user_id = ? AND is_unlocked = ?", userID, true).
		Preload("Module").
		Order("module_id ASC").
		Find(&progresses).Error

	return progresses, err
}
