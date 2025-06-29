package request

type ModuleRequest struct {
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	Thumbnail   string          `json:"thumbnail" binding:"omitempty"`    // Optional field for thumbnail URL
	Lessons     []LessonRequest `json:"lessons" binding:"omitempty,dive"` // Optional field for lessons
}
