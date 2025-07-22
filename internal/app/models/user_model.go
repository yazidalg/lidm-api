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
	RoleID            uint      `gorm:"not null"` // Remove default value
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

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
