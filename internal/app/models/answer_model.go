package models

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	QuestionID     uint
	Question       Question `gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID         uint
	User           User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OptionSelected string
	IsCorrect      bool
	Score          int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
