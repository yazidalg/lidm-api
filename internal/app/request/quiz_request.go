package request

type CreateQuizRequest struct {
	Status        string `json:"status,omitempty"`
	Mode          string `json:"mode" binding:"required,oneof=single_player multiplayer"`
	ModuleID      uint   `json:"module_id" binding:"required"`
	QuestionCount int    `json:"question_count,omitempty"` // Default 5 jika tidak diisi
	QuestionsIDs  []uint `json:"questions_ids,omitempty"`
	InviteCode    string `json:"invite_code,omitempty"`
	HostUserID    uint   `json:"host_user_id,omitempty"` // ID user yang membuat kuis (untuk mode multiplayer)
}

type UpdateQuizRequest struct {
	Status       string `json:"status,omitempty"`
	WinnerID     *uint  `json:"winner_id,omitempty"`
	QuestionsIDs []uint `json:"questions_ids,omitempty"`
}

type JoinQuizWithCodeRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}
