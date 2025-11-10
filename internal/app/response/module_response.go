package response

import (
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
)

type CustomModuleResponse struct {
	ID            uint32                `json:"ID"`
	Title         string                `json:"Title"`
	Description   string                `json:"Description"`
	Thumbnail     string                `json:"Thumbnail"`
	Icon          string                `json:"Icon"`
	OffsetX       int                   `json:"OffsetX"`
	OffsetY       int                   `json:"OffsetY"`
	CreatedAt     time.Time             `json:"CreatedAt"`
	UpdatedAt     time.Time             `json:"UpdatedAt"`
	VideoMaterial models.VideoMaterial  `json:"video_material"`
	ARExperiments []models.ARExperiment `json:"ar_experiments"`
}
