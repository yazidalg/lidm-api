package request

type CreateModuleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	SortOrder   uint   `json:"sort_order" binding:"required"`
	CourseID    uint   `json:"course_id" binding:"required"`
}

type UpdateModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   uint   `json:"sort_order"`
}
