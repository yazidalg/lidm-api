package models

import (
	"time"

	"gorm.io/gorm"
)

type ARExperiment struct {
	gorm.Model
	Title     string    `gorm:"not null" json:"title"`
	LinkAR    string    `gorm:"not null" json:"link_ar"`
	LinkEmbed string    `gorm:"not null" json:"link_embed"`
	OffsetX   float64   `gorm:"default:0" json:"offset_x"`
	OffsetY   float64   `gorm:"default:0" json:"offset_y"`
	ModuleID  uint      `gorm:"not null;index" json:"module_id"` // Foreign key to Module
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Removed SubMaterials relationship to avoid circular import
	// SubMaterials can reference this via ARExperimentID
}
