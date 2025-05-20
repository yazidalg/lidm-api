package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	ParticipantID uint          `gorm:"not null"`
	Participants  []Participant `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Answers       []Answer      `gorm:"many2many:quiz_answers;"`
	Questions     []*Question   `gorm:"many2many:quiz_questions;"`
	Status        string        `gorm:"type:enum('draft','active','finished');default:'draft'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
