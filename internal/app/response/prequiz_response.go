package response

import (
	"time"
	"gorm.io/gorm"
	"github.com/yazidalg/lidm_backend/internal/app/models"
)

// PrequizWithStatus represents a prequiz with additional status information
type PrequizWithStatus struct {
	ID              uint          `json:"ID"`
	CreatedAt       time.Time     `json:"CreatedAt"`
	UpdatedAt       time.Time     `json:"UpdatedAt"`
	DeletedAt       gorm.DeletedAt `json:"DeletedAt"`
	SubMaterialID   uint          `json:"SubMaterialID"`
	Question        string        `json:"Question"`
	Options         models.Options `json:"Options"`
	CorrectAnswer   string        `json:"CorrectAnswer"`
	Explanation     string        `json:"Explanation"`
	IsAlreadyAnswered bool        `json:"isAlreadyAnswered"`
}

// PrequizStatus represents the overall status of prequizzes
type PrequizStatus struct {
	Answered  int  `json:"answered"`
	Total     int  `json:"total"`
	Completed bool `json:"completed"`
}

// PrequizzesBySubMaterialResponse represents the response for getting prequizzes by submaterial
type PrequizzesBySubMaterialResponse struct {
	Success       bool                  `json:"success"`
	Message       string                `json:"message"`
	PrequizStatus PrequizStatus         `json:"prequiz_status"`
	Prequizzes    []PrequizWithStatus   `json:"prequizzes"`
}
