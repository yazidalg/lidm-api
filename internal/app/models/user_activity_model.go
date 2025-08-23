package models

import (
	"time"

	"gorm.io/gorm"
)

type UserActivity struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	ActivityType string `gorm:"not null;index"` // login, quiz_join, quiz_complete, module_view, module_complete
	Description  string // Human readable description
	MetaData     string `gorm:"type:json"` // JSON data for additional context
	IPAddress    string
	UserAgent    string
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}

// Activity type constants
const (
	ActivityTypeLogin        = "masuk"
	ActivityTypeLogout       = "keluar"
	ActivityTypeQuizJoin     = "gabung_kuis"
	ActivityTypeQuizComplete = "selesai_kuis"
	ActivityTypeQuizAnswer   = "jawab_kuis"

	ActivityTypeModuleView     = "lihat_modul"
	ActivityTypeModuleComplete = "selesai_modul"
	ActivityTypeProfileUpdate  = "update_profil"
)

// Helper methods for activity descriptions
func (ua *UserActivity) GetFormattedDescription() string {
	switch ua.ActivityType {
	case ActivityTypeLogin:
		return "Masuk ke sistem"
	case ActivityTypeLogout:
		return "Keluar dari sistem"
	case ActivityTypeQuizJoin:
		return "Bergabung dengan sesi kuis"
	case ActivityTypeQuizComplete:
		return "Menyelesaikan kuis"
	case ActivityTypeQuizAnswer:
		return "Menjawab pertanyaan kuis"

	case ActivityTypeModuleView:
		return "Melihat modul"
	case ActivityTypeModuleComplete:
		return "Menyelesaikan modul"
	case ActivityTypeProfileUpdate:
		return "Memperbarui informasi profil"
	default:
		return ua.Description
	}
}
