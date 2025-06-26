package request

type CourseRequest struct {
	Title       string                `json:"title" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Thumbnail   string                `json:"thumbnail" binding:"required"`
	Modules     []CreateModuleRequest `json:"modules"`
}
