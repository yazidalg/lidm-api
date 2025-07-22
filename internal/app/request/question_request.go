package request

type CreateQuestionRequest struct {
	ModuleID      *uint           `json:"module_id,omitempty"`
	Question      string          `json:"question" binding:"required"`
	AnswerTime    int32           `json:"answer_time" binding:"required"`
	ReadTime      int32           `json:"read_time" binding:"required"`
	CorrectAnswer string          `json:"correct_answer" binding:"required"`
	Options       QuestionOptions `json:"options" binding:"required"`
	Explanation   string          `json:"explanation" binding:"required"`
}

type QuestionOptions struct {
	OptionA string `json:"option_a" binding:"required"`
	OptionB string `json:"option_b" binding:"required"`
	OptionC string `json:"option_c" binding:"required"`
	OptionD string `json:"option_d" binding:"required"`
}

type UpdateQuestionRequest struct {
	ModuleID      *uint           `json:"module_id,omitempty"`
	Question      string          `json:"question"`
	AnswerTime    int32           `json:"answer_time"`
	ReadTime      int32           `json:"read_time"`
	CorrectAnswer string          `json:"correct_answer"`
	Options       QuestionOptions `json:"options" binding:"required"`
	Explanation   string          `json:"explanation"`
}
