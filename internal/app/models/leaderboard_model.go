package models

import (
	"time"

	"gorm.io/gorm"
)

type Leaderboard struct {
	gorm.Model
	UserID        uint  `gorm:"not null;index"`
	TotalScore    int64 `gorm:"default:0"`
	MatchesPlayed uint  `gorm:"default:0"`
	Wins          uint  `gorm:"default:0"`
	Losses        uint  `gorm:"default:0"`
	UpdatedAt     time.Time

	// Relationships
	User *User `gorm:"foreignKey:UserID"`
}
