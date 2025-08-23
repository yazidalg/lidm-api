package models

import "gorm.io/gorm"

type Resource struct {
	gorm.Model
	LessonID  uint   `gorm:"not null;index"` // Foreign key to Lesson
	Type      string `gorm:"not null"`       // Type of resource (e.g., "video", "document", "link")
	URL       string `gorm:"not null"`       // URL or path to the resource
	CreatedAt string `gorm:"autoCreateTime"` // Timestamp when the resource was created
	UpdatedAt string `gorm:"autoUpdateTime"` // Timestamp when the resource was last updated
}
