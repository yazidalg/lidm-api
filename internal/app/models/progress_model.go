package models

import (
	"time"

	"gorm.io/gorm"
)

type Progress struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index"` // Foreign key to User
	Completed   bool       `gorm:"default:false"`  // Indicates if the progress is completed
	CompletedAt *time.Time // Timestamp when the progress was completed

	// Relationships
	User *User `gorm:"foreignKey:UserID"` // Relationship to User
}
