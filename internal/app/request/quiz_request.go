package request

type QuizRequest struct {
	Question      string             `json:"question" binding:"required"`
	AnswerTime    int32              `json:"answer_time" binding:"required"`
	ReadTime      int32              `json:"read_time" binding:"required"`
	Options       QuizOptionsRequest `json:"options" binding:"required"`
	CorrectAnswer string             `json:"correct_answer" binding:"required"`
	Explanation   string             `json:"explanation" binding:"required"`
}

type QuizOptionsRequest struct {
	OptionA string `json:"option_a" binding:"required"`
	OptionB string `json:"option_b" binding:"required"`
	OptionC string `json:"option_c" binding:"required"`
	OptionD string `json:"option_d" binding:"required"`
}
