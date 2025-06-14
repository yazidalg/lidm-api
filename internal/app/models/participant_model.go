package models

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	QuizID     uint `gorm:"not null;index"`
	UserID     uint `gorm:"not null;index"`
	TotalScore int  `gorm:"default:0"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Quiz Quiz `gorm:"foreignKey:QuizID"`
}
