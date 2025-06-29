package models

import "gorm.io/gorm"

type Prequiz struct {
	gorm.Model
	CourseID      uint    `gorm:"not null;index"` // Foreign key to Course
	ModuleID      uint    `gorm:"not null;index"` // Foreign key to Module
	Options       Options `gorm:"type:text"`      // JSON string containing quiz options
	CorrectAnswer string  `gorm:"not null"`       // Correct answer for the quiz
	Explanation   string  `gorm:"type:text"`      // Explanation for the correct answer

	Module *Module `gorm:"foreignKey:ModuleID"` // Relationship to Module
}
