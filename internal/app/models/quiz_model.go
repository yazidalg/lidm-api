package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	Status    string `gorm:"type:enum('pending','in_progress','completed', 'cancelled');default:'pending'"`
	WinnerID  *uint  `gorm:"index"`
	CreatedAt *time.Time
	UpdatedAt *time.Time

	// Relationships
	Participants []Participant `gorm:"foreignKey:QuizID"`
	Questions    []Question    `gorm:"many2many:quiz_questions;"`
	Winner       *User         `gorm:"foreignKey:WinnerID"`
}
