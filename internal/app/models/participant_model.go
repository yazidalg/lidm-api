package models

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	QuizID             uint            `gorm:"not null;index"`
	UserID             uint            `gorm:"not null;index"`
	TotalScore         int             `gorm:"default:0"`
	CorrectAnswers     int             `gorm:"default:0"`     // Jumlah jawaban benar
	WrongAnswers       int             `gorm:"default:0"`     // Jumlah jawaban salah
	ConsecutiveCorrect int             `gorm:"default:0"`     // Jawaban benar berturut-turut maksimal
	CurrentStreak      int             `gorm:"default:0"`     // Streak saat ini (untuk perhitungan)
	IsFinished         bool            `gorm:"default:false"` // Apakah sudah selesai mengerjakan
	FinishedAt         *gorm.DeletedAt // Waktu selesai mengerjakan

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Quiz Quiz `gorm:"foreignKey:QuizID"`
}
