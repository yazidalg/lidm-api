package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	UserID    int
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	GameID    int
	Quiz      Quiz `gorm:"foreignKey:GameID"`
	Score     int
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
