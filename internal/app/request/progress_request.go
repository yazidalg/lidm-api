package request

type ProgressRequest struct {
	UserID   uint32 `json:"user_id" validate:"required"`
	LessonID uint32 `json:"lesson_id" validate:"required"`
}
