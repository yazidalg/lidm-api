package models

import (
	"gorm.io/gorm"
)

type ForgetPasswordModel struct {
	gorm.Model
	UserId    uint
	User      User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Email     string `gorm:"unique;not null"`
	OTP       string `gorm:"not null"`
	ExpiresAt int64
}
