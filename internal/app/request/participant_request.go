package request

type CreateParticipantRequest struct {
	UserID     uint `json:"user_id" binding:"required"`
	QuizID     uint `json:"quiz_id" binding:"required"`
	TotalScore int  `json:"total_score"`
}

type UpdateParticipantRequest struct {
	UserID     uint `json:"user_id"`
	QuizID     uint `json:"quiz_id"`
	TotalScore int  `json:"total_score"`
}
