package request

type CreateParticipantRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	QuizID     uint   `json:"quiz_id" binding:"required"`
	TotalScore int    `json:"total_score" binding:"required"`
	TotalXP    int    `json:"total_xp" binding:"required"`
	Result     string `json:"result" binding:"required"`
}

type UpdateParticipantRequest struct {
	UserID     uint   `json:"user_id"`
	QuizID     uint   `json:"quiz_id"`
	TotalScore int    `json:"total_score"`
	TotalXP    int    `json:"total_xp"`
	Result     string `json:"result"`
}
