package models

import "gorm.io/gorm"

type Prequiz struct {
	gorm.Model
	LessonID      uint    `gorm:"not null;index"` // Foreign key to Course
	UserID        uint    `gorm:"not null;index"` // Foreign key to User
	Question      string  `gorm:"not null"`       // Quiz question
	Options       Options `gorm:"embedded;"`
	CorrectAnswer string  `gorm:"not null"`  // Correct answer for the quiz
	Explanation   string  `gorm:"type:text"` // Explanation for the correct answer

	Lesson *Lesson `gorm:"foreignKey:LessonID"` // Relationship to Lesson
}
