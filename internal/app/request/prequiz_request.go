package request

type PrequizRequest struct {
	LessonID      uint            `json:"lesson_id" binding:"required"`
	Options       QuestionOptions `json:"options" binding:"required"`
	Question      string          `json:"question" binding:"required"`
	CorrectAnswer string          `json:"correct_answer" binding:"required"`
	Explanation   string          `json:"explanation" binding:"required"`
}
