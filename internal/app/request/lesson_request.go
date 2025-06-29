package request

type LessonRequest struct {
	ModuleID  uint   `json:"module_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	SortOrder uint16 `json:"sort_order" binding:"omitempty"` // Optional field for sort order
}
