package models

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	UserID     uint
	User       User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QuizID     uint
	Quiz       Quiz `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalScore int
	TotalXP    int
	Result     string
}
