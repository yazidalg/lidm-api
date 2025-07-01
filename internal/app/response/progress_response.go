package response

type ProgressResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	LessonID    uint   `json:"lesson_id"`
	Completed   bool   `json:"completed"`
	CompletedAt string `json:"completed_at"` // ISO 8601 format
}
