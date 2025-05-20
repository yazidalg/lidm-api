package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	Participants []Participant `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Questions    []Question    `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Answers      []Answer      `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status       string        `gorm:"type:enum('draft','active','finished');default:'draft'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
