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
	Lives             int        `gorm:"default:5" json:"lives"`                 // Nyawa (hanya berkurang di mode single player)
	LifeResetAt       *time.Time `json:"life_reset_at"`                           // Waktu terakhir reset nyawa harian
	ProfilePicture    string     `gorm:"type:text" json:"profile_picture"`
	RoleID            uint       `gorm:"not null"` // Remove default value
	CurrentStreak     int        `gorm:"default:0" json:"current_streak"`     // Streak hari berturut-turut
	MaxStreak         int        `gorm:"default:0" json:"max_streak"`         // Streak terpanjang
	LastActiveDate    *time.Time `json:"last_active_date"`                    // Tanggal terakhir aktif
	CreatedAt         time.Time  `gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime"`

	// Relationships
	Role         Role          `gorm:"foreignKey:RoleID" json:"role"`
	Leaderboard  Leaderboard   `gorm:"foreignKey:UserID"`
	Participants []Participant `gorm:"foreignKey:UserID"`
	Progress     Progress      `gorm:"foreignKey:UserID"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func (u *User) IsAdmin() bool {
	return u.Role.Name == RoleAdmin
}

func (u *User) IsUser() bool {
	return u.Role.Name == RoleUser
}
