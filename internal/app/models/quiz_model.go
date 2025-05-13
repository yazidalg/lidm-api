package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	ID            int `json:"id" gorm:"primaryKey"`
	Question      string
	AnswerTime    time.Time
	ReadTime      time.Time
	Options       Options `gorm:"embedded;"`
	CorrectAnswer string
	Explanation   string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Options struct {
	OptionA string
	OptionB string
	OptionC string
	OptionD string
}
