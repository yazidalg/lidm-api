package response

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

// PrequizWithStatus represents a prequiz with additional status information
type PrequizWithStatus struct {
	ID                uint           `json:"ID"`
	CreatedAt         time.Time      `json:"CreatedAt"`
	UpdatedAt         time.Time      `json:"UpdatedAt"`
	DeletedAt         gorm.DeletedAt `json:"DeletedAt"`
	ModuleID          uint           `json:"ModuleID"`
	Question          string         `json:"Question"`
	Options           models.Options `json:"Options"`
	CorrectAnswer     string         `json:"CorrectAnswer"`
	Explanation       string         `json:"Explanation"`
	IsAlreadyAnswered bool           `json:"isAlreadyAnswered"`
}

// PrequizStatus represents the overall status of prequizzes
type PrequizStatus struct {
	Answered  int  `json:"answered"`
	Total     int  `json:"total"`
	Completed bool `json:"completed"`
}

// PrequizzesByModuleResponse represents the response for getting prequizzes by module
type PrequizzesByModuleResponse struct {
	Success       bool                `json:"success"`
	Message       string              `json:"success"`
	PrequizStatus PrequizStatus       `json:"prequiz_status"`
	Prequizzes    []PrequizWithStatus `json:"prequizzes"`
}
