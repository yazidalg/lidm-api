package models

import (
	"time"

	"gorm.io/gorm"
)

type Flashcard struct {
	gorm.Model
	ModuleID  uint      `gorm:"not null;index" json:"module_id"`
	FrontText string    `gorm:"type:text;not null" json:"front_text"`
	BackText  string    `gorm:"type:text;not null" json:"back_text"`
	Order     int       `gorm:"not null;default:0" json:"order"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Direct relationship to Module
	// Use ModuleID to reference the parent
	UserFlashcardProgresses []UserFlashcardProgress `gorm:"foreignKey:FlashcardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_progresses,omitempty"`
}
