package models

import (
	"time"

	"gorm.io/gorm"
)

type QuizQuestion struct {
	gorm.Model
	SubMaterialID uint      `gorm:"not null;index" json:"sub_material_id"`
	Question      string    `gorm:"type:text;not null" json:"question"`
	OptionA       string    `gorm:"not null" json:"option_a"`
	OptionB       string    `gorm:"not null" json:"option_b"`
	OptionC       string    `gorm:"not null" json:"option_c"`
	OptionD       string    `gorm:"not null" json:"option_d"`
	CorrectAnswer string    `gorm:"size:1;not null" json:"correct_answer"` // A, B, C, atau D
	Description   string    `gorm:"type:text" json:"description"`
	Order         int       `gorm:"not null;default:0" json:"order"` // Order dalam sub material (1, 2, 3)
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Note: Removed SubMaterial relationship to avoid circular import
	// Use SubMaterialID to reference the parent
}
