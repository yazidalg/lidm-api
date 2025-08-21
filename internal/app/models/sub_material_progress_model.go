package models

import "gorm.io/gorm"

// SubMaterialProgress tracks completion status for video quizzes and prequizzes within a sub_material
type SubMaterialProgress struct {
	gorm.Model
	UserID        uint `gorm:"not null;index" json:"user_id"`
	SubMaterialID uint `gorm:"not null;index" json:"sub_material_id"`
	
	// Video completion status
	VideoCompleted bool `gorm:"default:false" json:"video_completed"`
	
	// Prequiz completion status
	PrequizzesCompleted bool `gorm:"default:false" json:"prequizzes_completed"`
	
	// Overall sub_material completion (both video and prequizzes done)
	Completed bool `gorm:"default:false" json:"completed"`

	// Relationships
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	SubMaterial *SubMaterial `gorm:"foreignKey:SubMaterialID" json:"sub_material,omitempty"`
}

// Unique constraint to ensure one progress record per user per sub_material
func (SubMaterialProgress) TableName() string {
	return "sub_material_progresses"
}
