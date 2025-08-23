package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string
	Description string
	Thumbnail   string
	Icon        string
	OffsetX     uint16    `gorm:"default:0"` // Offset for the lesson in the module
	OffsetY     uint16    `gorm:"default:0"` // Offset for the lesson in the module
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Quizzes       []Quiz          `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"quizzes,omitempty"`        // Quiz yang terkait dengan module ini
	VideoMaterial []VideoMaterial `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"video_material,omitempty"` // Video material yang terkait dengan module ini
	ARExperiment  *ARExperiment   `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"ar_experiment,omitempty"`  // AR experiment yang terkait dengan module ini
}
