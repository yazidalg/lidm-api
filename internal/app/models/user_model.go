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
	Role              string    `gorm:"default:'user'"` // Default role is 'user', can be 'admin' or 'superadmin'
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Leaderboard  Leaderboard   `gorm:"foreignKey:UserID"`
	Participants []Participant `gorm:"foreignKey:UserID"`
	Progress     Progress      `gorm:"foreignKey:UserID"`
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsUser() bool {
	return u.Role == RoleUser
}
