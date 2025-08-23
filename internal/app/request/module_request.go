package request

import "github.com/yazidalg/lidm_backend/internal/app/models"

type ModuleRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	OffsetX     *float64 `json:"offset_x" binding:"omitempty"`  // Optional field for X position
	OffsetY     *float64 `json:"offset_y" binding:"omitempty"`  // Optional field for Y position
	Icon        *string  `json:"icon" binding:"omitempty"`      // Optional field for icon path
	Thumbnail   *string  `json:"thumbnail" binding:"omitempty"` // Optional field for thumbnail URL
}

type UpdateModuleRequest struct {
	Title         string                `json:"title" binding:"required"`
	Description   string                `json:"description" binding:"required"`
	OffsetX       *float64              `json:"offset_x" binding:"omitempty"`  // Optional field for X position
	OffsetY       *float64              `json:"offset_y" binding:"omitempty"`  // Optional field for Y position
	Icon          *string               `json:"icon" binding:"omitempty"`      // Optional field for icon path
	Thumbnail     *string               `json:"thumbnail" binding:"omitempty"` // Optional field for thumbnail URL
	VideoMaterial *models.VideoMaterial `json:"video_material" binding:"omitempty"`
}

type CreateModuleWithVideoRequest struct {
	Title         string                 `json:"title" binding:"required"`
	Description   string                 `json:"description" binding:"required"`
	OffsetX       *float64               `json:"offset_x" binding:"omitempty"`
	OffsetY       *float64               `json:"offset_y" binding:"omitempty"`
	Icon          *string                `json:"icon" binding:"omitempty"`
	Thumbnail     *string                `json:"thumbnail" binding:"omitempty"`
	VideoMaterial []models.VideoMaterial `json:"video_material" binding:"omitempty"`
}
