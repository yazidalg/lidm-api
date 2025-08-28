package models

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	ModuleID      *uint `gorm:"index"` // Materi/topik pertanyaan
	Question      string
	AnswerTime    int32
	ReadTime      int32
	Options       Options `gorm:"embedded;"`
	CorrectAnswer string
	QuestionType  string // "hots" or "regular"
	Explanation   string
	Quiz          []Quiz    `gorm:"many2many:quiz_questions;"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Module *Module `gorm:"foreignKey:ModuleID"`
}

type Options struct {
	OptionA string
	OptionB string
	OptionC string
	OptionD string
}
