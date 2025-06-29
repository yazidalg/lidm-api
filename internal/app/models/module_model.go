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
	Thumbnail   string
	CourseID    uint      `gorm:"not null;index"` // Foreign key to Course
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relationship to Lesson
}
