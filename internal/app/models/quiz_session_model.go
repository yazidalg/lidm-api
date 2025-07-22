package models

import (
	"time"

	"gorm.io/gorm"
)

// QuizSession menyimpan detail setiap sesi pertanyaan dalam quiz
type QuizSession struct {
	gorm.Model
	QuizID        uint      `gorm:"not null;index"`
	ParticipantID uint      `gorm:"not null;index"`
	QuestionID    uint      `gorm:"not null;index"`
	UserAnswer    string    // Jawaban yang dipilih user (A, B, C, atau D)
	IsCorrect     bool      `gorm:"default:false"`
	ResponseTime  int32     // Waktu response dalam milidetik
	PointsEarned  int       `gorm:"default:0"` // Point yang didapat dari pertanyaan ini
	AnsweredAt    time.Time `gorm:"autoCreateTime"`

	// Relationships
	Quiz        Quiz        `gorm:"foreignKey:QuizID"`
	Participant Participant `gorm:"foreignKey:ParticipantID"`
	Question    Question    `gorm:"foreignKey:QuestionID"`
}
