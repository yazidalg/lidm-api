package models

import (
	"time"

	"gorm.io/gorm"
)

type Progress struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index"` // Foreign key to User
	LessonID    uint       `gorm:"not null;index"` // Foreign key to Lesson
	Completed   bool       `gorm:"default:false"`  // Indicates if the lesson is completed
	CompletedAt *time.Time // Timestamp when the lesson was completed

	// Relationships
	User   User   `gorm:"foreignKey:UserID"`   // Relationship to User
	Lesson Lesson `gorm:"foreignKey:LessonID"` // Relationship to Lesson
}
