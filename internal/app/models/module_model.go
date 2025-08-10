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
	Lessons      []Lesson      `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"lessons,omitempty"` // Legacy relationship
	SubMaterials []SubMaterial `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"sub_materials,omitempty"`
}
