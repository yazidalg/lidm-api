package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title       string
	Description string
	Thumbnail   string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Modules []Module `gorm:"foreignKey:CourseID"` // Relationship to Module
}
