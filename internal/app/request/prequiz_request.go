package request

type PrequizRequest struct {
	UserID        uint            `json:"user_id" binding:"required"`
	LessonID      uint            `json:"lesson_id" binding:"required"`
	Options       QuestionOptions `json:"question" binding:"required"`
	Question      string          `json:"question" binding:"required"`
	CorrectAnswer string          `json:"correct_answer" binding:"required"`
	Explanation   string          `json:"explanation" binding:"required"`
}
