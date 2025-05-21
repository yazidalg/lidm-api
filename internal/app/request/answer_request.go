package request

type CreateAnswerRequest struct {
	QuestionID     uint   `json:"question_id" binding:"required"`
	UserID         uint   `json:"user_id" binding:"required"`
	QuizID         []uint `json:"quiz_id" binding:"required"`
	OptionSelected string `json:"option_selected" binding:"required"`
	IsCorrect      bool   `json:"is_correct" binding:"required"`
	Score          int    `json:"score" binding:"required"`
}
