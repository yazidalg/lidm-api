package models

import (
	"time"

	"gorm.io/gorm"
)

type VideoMaterial struct {
	gorm.Model
	SubMaterialID uint      `gorm:"not null;uniqueIndex" json:"sub_material_id"` // One-to-one dengan SubMaterial
	Title         string    `gorm:"not null" json:"title"`
	YoutubeLink   string    `gorm:"not null" json:"youtube_link"`
	Duration      int       `json:"duration"` // Duration in seconds (optional)
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	VideoQuizzes []VideoQuiz `gorm:"foreignKey:VideoMaterialID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"video_quizzes,omitempty"`
	// Note: Removed SubMaterial relationship to avoid circular import
	// Use SubMaterialID to reference the parent
}
