package models

import (
	"time"

	"gorm.io/gorm"
)

type ModuleProgress struct {
	gorm.Model
	UserID     uint      `gorm:"not null;index;uniqueIndex:idx_user_module" json:"user_id"`        // Foreign key to User
	ModuleID   uint      `gorm:"not null;index;uniqueIndex:idx_user_module" json:"module_id"`      // Foreign key to Module
	IsUnlocked bool      `gorm:"default:false" json:"is_unlocked"`                                 // Whether the module is unlocked for this user
	IsComplete bool      `gorm:"default:false" json:"is_complete"`                                 // Whether the module is completed
	StartedAt  *time.Time `json:"started_at"`                                                       // When user first accessed this module
	CompletedAt *time.Time `json:"completed_at"`                                                    // When user completed this module
	Progress   float32   `gorm:"default:0" json:"progress"`                                        // Progress percentage (0-100)

	// Relationships
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Module *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}

// Helper methods
func (mp *ModuleProgress) CalculateProgress(db *gorm.DB) float32 {
	if mp.ModuleID == 0 || mp.UserID == 0 {
		return 0
	}

	var totalItems float32 = 0
	var completedItems float32 = 0

	// Count total prequizzes in module
	var totalPrequizzes int64
	db.Model(&Prequiz{}).Where("module_id = ?", mp.ModuleID).Count(&totalPrequizzes)
	totalItems += float32(totalPrequizzes)

	// Count completed prequizzes by user
	var completedPrequizzes int64
	db.Table("prequiz_user_answers").
		Joins("JOIN prequizzes ON prequiz_user_answers.prequiz_id = prequizzes.id").
		Where("prequizzes.module_id = ? AND prequiz_user_answers.user_id = ?", mp.ModuleID, mp.UserID).
		Count(&completedPrequizzes)
	completedItems += float32(completedPrequizzes)

	// Count total video quizzes in module
	var totalVideoQuizzes int64
	db.Table("video_quizzes").
		Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Where("video_materials.module_id = ?", mp.ModuleID).
		Count(&totalVideoQuizzes)
	totalItems += float32(totalVideoQuizzes)

	// Count completed video quizzes by user
	var completedVideoQuizzes int64
	db.Table("video_quiz_user_answers").
		Joins("JOIN video_quizzes ON video_quiz_user_answers.video_quiz_id = video_quizzes.id").
		Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Where("video_materials.module_id = ? AND video_quiz_user_answers.user_id = ?", mp.ModuleID, mp.UserID).
		Count(&completedVideoQuizzes)
	completedItems += float32(completedVideoQuizzes)

	if totalItems == 0 {
		return 0
	}

	progress := (completedItems / totalItems) * 100
	if progress > 100 {
		progress = 100
	}

	return progress
}

// MarkAsStarted sets the started timestamp if not already set
func (mp *ModuleProgress) MarkAsStarted() {
	if mp.StartedAt == nil {
		now := time.Now()
		mp.StartedAt = &now
	}
}

// MarkAsCompleted sets the module as complete with timestamp
func (mp *ModuleProgress) MarkAsCompleted() {
	mp.IsComplete = true
	mp.Progress = 100
	if mp.CompletedAt == nil {
		now := time.Now()
		mp.CompletedAt = &now
	}
}

// CheckAndUnlockNextModule unlocks the next module if this one is completed
func (mp *ModuleProgress) CheckAndUnlockNextModule(db *gorm.DB) error {
	if !mp.IsComplete {
		return nil
	}

	// Find the next module (by ID order)
	var nextModule Module
	err := db.Where("id > ?", mp.ModuleID).Order("id ASC").First(&nextModule).Error
	if err != nil {
		// No next module found, that's okay
		return nil
	}

	// Check if next module is already unlocked
	var nextProgress ModuleProgress
	err = db.Where("user_id = ? AND module_id = ?", mp.UserID, nextModule.ID).First(&nextProgress).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new progress entry for next module
		nextProgress = ModuleProgress{
			UserID:     mp.UserID,
			ModuleID:   nextModule.ID,
			IsUnlocked: true,
			IsComplete: false,
			Progress:   0,
		}
		return db.Create(&nextProgress).Error
	} else if err != nil {
		return err
	} else if !nextProgress.IsUnlocked {
		// Update existing entry to unlock
		nextProgress.IsUnlocked = true
		return db.Save(&nextProgress).Error
	}

	return nil
}
