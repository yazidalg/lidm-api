package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	UserID    uint
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuizID    uint
	Quiz      Quiz `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Score     int
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
