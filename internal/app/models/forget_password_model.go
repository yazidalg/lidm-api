package models

import (
	"gorm.io/gorm"
)

type ForgotPassword struct {
	gorm.Model
	UserID    uint
	User      User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Email     string
	OTP       string
	ExpiresAt int64
	Used      bool `gorm:"default:false"`
}
