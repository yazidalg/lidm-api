package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type ModuleRepositoryInterface interface {
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

func (r *moduleRepository) GetModuleByID(id uint32) (*models.Module, error) {
	var module models.Module

	// Preload related entities if necessary
	err := r.db.First(&module, id).Error

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) GetAllModules() ([]models.Module, error) {
	var modules []models.Module

	// Preload related entities if necessary
	err := r.db.Find(&modules).Error

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
	existingModule.SortOrder = module.SortOrder

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
