package models

import (
	"time"

	"gorm.io/gorm"
)

type Lesson struct {
	gorm.Model
	ModuleID  uint `gorm:"not null;index"` // Foreign key to Module
	Title     string
	Content   string
	SortOrder uint16    `gorm:"default:0"` // Order of the lesson within the
	OffsetX   uint16    `gorm:"default:0"` // Offset for the lesson in the module
	OffsetY   uint16    `gorm:"default:0"` // Offset for the lesson in the module
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Module   *Module    `gorm:"foreignKey:ModuleID"` // Relationship to Module
	Progress []Progress `gorm:"foreignKey:LessonID"` // Relationship to Progress
}
