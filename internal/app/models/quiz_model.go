package models

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	Question      string
	AnswerTime    int32
	ReadTime      int32
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
