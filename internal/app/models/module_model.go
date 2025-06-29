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
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE"` // Relationship to Lesson
}
