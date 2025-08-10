package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	OffsetX     float64   `json:"offset_x"`
	OffsetY     float64   `json:"offset_y"`
	Icon        string    `json:"icon"`
	Thumbnail   string    `json:"thumbnail"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Lessons      []Lesson      `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"lessons,omitempty"` // Legacy relationship
	SubMaterials []SubMaterial `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"sub_materials,omitempty"`
}
