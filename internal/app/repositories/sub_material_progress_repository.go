package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type SubMaterialProgressRepositoryInterface interface {
	CreateOrUpdateProgress(progress *models.SubMaterialProgress) (*models.SubMaterialProgress, error)
	GetByUserAndSubMaterial(userID uint, subMaterialID uint) (*models.SubMaterialProgress, error)
	GetByUserAndModule(userID uint, moduleID uint) ([]models.SubMaterialProgress, error)
	UpdateVideoCompleted(userID uint, subMaterialID uint, completed bool) error
	UpdatePrequizzesCompleted(userID uint, subMaterialID uint, completed bool) error
	UpdateOverallCompleted(userID uint, subMaterialID uint) error
}

type subMaterialProgressRepository struct {
	db *gorm.DB
}

func NewSubMaterialProgressRepository(db *gorm.DB) SubMaterialProgressRepositoryInterface {
	return &subMaterialProgressRepository{db: db}
}

func (r *subMaterialProgressRepository) CreateOrUpdateProgress(progress *models.SubMaterialProgress) (*models.SubMaterialProgress, error) {
	// Try to find existing progress
	existing, err := r.GetByUserAndSubMaterial(progress.UserID, progress.SubMaterialID)
	if err != nil {
		// If not found, create new
		if err := r.db.Create(progress).Error; err != nil {
			return nil, err
		}
		return progress, nil
	}

	// Update existing
	existing.VideoCompleted = progress.VideoCompleted
	existing.PrequizzesCompleted = progress.PrequizzesCompleted
	existing.Completed = progress.Completed

	if err := r.db.Save(existing).Error; err != nil {
		return nil, err
	}
	return existing, nil
}

func (r *subMaterialProgressRepository) GetByUserAndSubMaterial(userID uint, subMaterialID uint) (*models.SubMaterialProgress, error) {
	var progress models.SubMaterialProgress
	err := r.db.Where("user_id = ? AND sub_material_id = ?", userID, subMaterialID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

func (r *subMaterialProgressRepository) GetByUserAndModule(userID uint, moduleID uint) ([]models.SubMaterialProgress, error) {
	var progresses []models.SubMaterialProgress
	err := r.db.Joins("JOIN sub_materials ON sub_material_progresses.sub_material_id = sub_materials.id").
		Where("sub_material_progresses.user_id = ? AND sub_materials.module_id = ?", userID, moduleID).
		Find(&progresses).Error
	return progresses, err
}

func (r *subMaterialProgressRepository) UpdateVideoCompleted(userID uint, subMaterialID uint, completed bool) error {
	// Get or create progress record
	progress, err := r.GetByUserAndSubMaterial(userID, subMaterialID)
	if err != nil {
		// Create new if not exists
		progress = &models.SubMaterialProgress{
			UserID:        userID,
			SubMaterialID: subMaterialID,
			VideoCompleted: completed,
		}
		if err := r.db.Create(progress).Error; err != nil {
			return err
		}
	} else {
		// Update existing
		progress.VideoCompleted = completed
		if err := r.db.Save(progress).Error; err != nil {
			return err
		}
	}

	// Update overall completion status
	return r.UpdateOverallCompleted(userID, subMaterialID)
}

func (r *subMaterialProgressRepository) UpdatePrequizzesCompleted(userID uint, subMaterialID uint, completed bool) error {
	// Get or create progress record
	progress, err := r.GetByUserAndSubMaterial(userID, subMaterialID)
	if err != nil {
		// Create new if not exists
		progress = &models.SubMaterialProgress{
			UserID:             userID,
			SubMaterialID:      subMaterialID,
			PrequizzesCompleted: completed,
		}
		if err := r.db.Create(progress).Error; err != nil {
			return err
		}
	} else {
		// Update existing
		progress.PrequizzesCompleted = completed
		if err := r.db.Save(progress).Error; err != nil {
			return err
		}
	}

	// Update overall completion status
	return r.UpdateOverallCompleted(userID, subMaterialID)
}

func (r *subMaterialProgressRepository) UpdateOverallCompleted(userID uint, subMaterialID uint) error {
	progress, err := r.GetByUserAndSubMaterial(userID, subMaterialID)
	if err != nil {
		return err
	}

	// Mark as completed if both video and prequizzes are done
	progress.Completed = progress.VideoCompleted && progress.PrequizzesCompleted
	return r.db.Save(progress).Error
}
