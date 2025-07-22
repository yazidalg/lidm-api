package request

type CreateQuizRequest struct {
	Status        string `json:"status,omitempty"`
	Mode          string `json:"mode" binding:"required,oneof=single_player multiplayer"`
	ModuleID      uint   `json:"module_id" binding:"required"`
	QuestionCount int    `json:"question_count,omitempty"` // Default 5 jika tidak diisi
	QuestionsIDs  []uint `json:"questions_ids,omitempty"`
}

type UpdateQuizRequest struct {
	Status       string `json:"status,omitempty"`
	WinnerID     *uint  `json:"winner_id,omitempty"`
	QuestionsIDs []uint `json:"questions_ids,omitempty"`
}
