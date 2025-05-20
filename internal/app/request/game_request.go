package request

type GameRequest struct {
	UserID string `json:"user_id" binding:"required"`
	QuizID string `json:"quiz_id" binding:"required"`
	Score  int    `json:"score"`
	Answer string `json:"answer"`
}
