package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string
	Description string
	SortOrder   uint16
	CourseID    uint      `gorm:"not null;index"` // Foreign key to Course
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Course *Course `gorm:"foreignKey:CourseID"` // Relationship to Course
}
