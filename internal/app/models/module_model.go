package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string
	Description string
	Thumbnail   string
	Icon        string
	OffsetX     uint16    `gorm:"default:0"` // Offset for the lesson in the module
	OffsetY     uint16    `gorm:"default:0"` // Offset for the lesson in the module
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Quizzes        []Quiz           `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"quizzes,omitempty"`         // Quiz yang terkait dengan module ini
	VideoMaterial  *VideoMaterial   `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"video_material,omitempty"`  // Video material yang terkait dengan module ini (single object)
	ARExperiments  []ARExperiment   `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"ar_experiments,omitempty"`  // AR experiments yang terkait dengan module ini (array)
	Prequizzes     []Prequiz        `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"prequizzes,omitempty"`      // Prequizzes terkait dengan module ini
	Flashcards     []Flashcard      `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"flashcards,omitempty"`      // Flashcards terkait dengan module ini
	ModuleProgress []ModuleProgress `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE" json:"module_progress,omitempty"` // Progress tracking untuk module ini
}
