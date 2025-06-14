package models

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	QuestionID     uint `gorm:"not null;index"`
	ParticipantID  uint `gorm:"not null;index"`
	OptionSelected string
	IsCorrect      bool
	Score          int
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relationships
	Question    Question    `gorm:"foreignKey:QuestionID"`
	Participant Participant `gorm:"foreignKey:ParticipantID"`
}
