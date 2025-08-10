package models

import (
	"time"

	"gorm.io/gorm"
)

// UserFlashcardProgress implements FSRS (Free Spaced Repetition Scheduler) algorithm
type UserFlashcardProgress struct {
	gorm.Model
	UserID      uint `gorm:"not null;index" json:"user_id"`
	FlashcardID uint `gorm:"not null;index" json:"flashcard_id"`

	// FSRS Algorithm Fields
	Stability  float64    `gorm:"default:0" json:"stability"`  // Memory stability
	Difficulty float64    `gorm:"default:0" json:"difficulty"` // Item difficulty
	Elapsed    int        `gorm:"default:0" json:"elapsed"`    // Days since last review
	Scheduled  int        `gorm:"default:0" json:"scheduled"`  // Scheduled days for next review
	Reps       int        `gorm:"default:0" json:"reps"`       // Number of repetitions
	Lapses     int        `gorm:"default:0" json:"lapses"`     // Number of lapses
	State      int        `gorm:"default:0" json:"state"`      // Card state (New=0, Learning=1, Review=2, Relearning=3)
	LastReview *time.Time `json:"last_review"`                 // Last review date
	Due        time.Time  `json:"due"`                         // Next due date

	// Review History
	ReviewCount  int     `gorm:"default:0" json:"review_count"`
	AverageGrade float64 `gorm:"default:0" json:"average_grade"` // Average rating given by user

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Flashcard Flashcard `gorm:"foreignKey:FlashcardID" json:"flashcard,omitempty"`
}

// FSRS States
const (
	StateNew        = 0 // New card
	StateLearning   = 1 // Learning phase
	StateReview     = 2 // Review phase
	StateRelearning = 3 // Relearning after lapse
)

// FSRS Grades/Ratings
const (
	GradeAgain = 1 // Failed recall
	GradeHard  = 2 // Difficult recall
	GradeGood  = 3 // Normal recall
	GradeEasy  = 4 // Easy recall
)

// Helper methods
func (ufp *UserFlashcardProgress) IsNew() bool {
	return ufp.State == StateNew
}

func (ufp *UserFlashcardProgress) IsLearning() bool {
	return ufp.State == StateLearning
}

func (ufp *UserFlashcardProgress) IsReview() bool {
	return ufp.State == StateReview
}

func (ufp *UserFlashcardProgress) IsRelearning() bool {
	return ufp.State == StateRelearning
}

func (ufp *UserFlashcardProgress) IsDue() bool {
	return time.Now().After(ufp.Due) || time.Now().Equal(ufp.Due)
}
