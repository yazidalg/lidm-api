package request

type ModuleRequest struct {
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	OffsetX     *float64        `json:"offset_x" binding:"omitempty"`     // Optional field for X position
	OffsetY     *float64        `json:"offset_y" binding:"omitempty"`     // Optional field for Y position
	Icon        *string         `json:"icon" binding:"omitempty"`         // Optional field for icon path
	Thumbnail   *string         `json:"thumbnail" binding:"omitempty"`    // Optional field for thumbnail URL
	Lessons     []LessonRequest `json:"lessons" binding:"omitempty,dive"` // Optional field for lessons
}
