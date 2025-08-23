package models

import "gorm.io/gorm"

type Prequiz struct {
	gorm.Model
	ModuleID      uint    `gorm:"not null;index"` // Foreign key to Module
	Question      string  `gorm:"not null"`       // Quiz question
	Options       Options `gorm:"embedded;"`
	CorrectAnswer string  `gorm:"not null"`  // Correct answer for the quiz
	Explanation   string  `gorm:"type:text"` // Explanation for the correct answer

	// Relationships
	UserAnswers []PrequizUserAnswer `gorm:"foreignKey:PrequizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_answers,omitempty"`
}

// Model untuk menyimpan jawaban user pada prequiz
type PrequizUserAnswer struct {
	gorm.Model
	PrequizID  uint   `gorm:"not null;index" json:"prequiz_id"` // Foreign key ke Prequiz
	UserID     uint   `gorm:"not null;index" json:"user_id"`    // Foreign key ke User
	Answer     string `gorm:"not null" json:"answer"`           // Jawaban yang dipilih user (A/B/C/D)
	IsCorrect  bool   `gorm:"not null" json:"is_correct"`       // Apakah jawaban benar
	AnsweredAt int64  `gorm:"not null" json:"answered_at"`      // Timestamp kapan dijawab (Unix timestamp)

	// Relationships
	Prequiz *Prequiz `gorm:"foreignKey:PrequizID" json:"prequiz,omitempty"`
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
