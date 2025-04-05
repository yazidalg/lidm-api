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
	VerificationToken string
	Point             int32
	Streak            int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"-:all"`
}
