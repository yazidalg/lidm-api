package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	Status        string `gorm:"type:enum('pending','in_progress','completed', 'cancelled');default:'pending'"`
	Mode          string `gorm:"type:enum('single_player','multiplayer');default:'multiplayer'"` // Mode permainan
	ModuleID      *uint  `gorm:"index"`                                                          // Materi yang dipilih untuk quiz
	InviteCode    string `gorm:"unique;size:8"`                                                  // Kode unik untuk mengundang teman
	WinnerID      *uint  `gorm:"index"`
	HostUserID    uint   `gorm:"index"`      // User yang membuat quiz (untuk mode multiplayer)
	QuestionCount int    `gorm:"default:10"` // Jumlah pertanyaan (default 10)
	CreatedAt     *time.Time
	UpdatedAt     *time.Time

	// Relationships
	Participants []Participant `gorm:"foreignKey:QuizID"`
	Questions    []Question    `gorm:"many2many:quiz_questions;"`
	Winner       *User         `gorm:"foreignKey:WinnerID"`
	Host         User          `gorm:"foreignKey:HostUserID"`
	Module       *Module       `gorm:"foreignKey:ModuleID"`
}
