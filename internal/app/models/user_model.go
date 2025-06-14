package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email             string `gorm:"unique"`
	Name              string
	Password          string
	Class             string
	IsVerified        bool
	VerificationToken string
	Point             int32
	TotalXP           int32
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Leaderboard  Leaderboard   `gorm:"foreignKey:UserID"`
	Participants []Participant `gorm:"foreignKey:UserID"`
}
