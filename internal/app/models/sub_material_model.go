package models

import (
	"time"

	"gorm.io/gorm"
)

type SubMaterial struct {
	gorm.Model
	ModuleID    uint      `gorm:"not null;index" json:"module_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Order       int       `gorm:"not null;default:0" json:"order"` // Order dalam module
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships - Setiap SubMaterial punya 3 komponen utama
	VideoMaterial  *VideoMaterial `gorm:"foreignKey:SubMaterialID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"video_material,omitempty"`
	ARExperimentID *uint         `gorm:"index" json:"ar_experiment_id,omitempty"`
	ARExperiment   *ARExperiment `gorm:"foreignKey:ARExperimentID" json:"ar_experiment,omitempty"`
	Prequizzes     []Prequiz     `gorm:"foreignKey:SubMaterialID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"prequizzes,omitempty"`
}
