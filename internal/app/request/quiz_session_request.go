package request

// CreateQuizSessionRequest untuk membuat quiz session baru
type CreateQuizSessionRequest struct {
	Mode     string `json:"mode" binding:"required,oneof=single_player multiplayer"` // Mode permainan
	ModuleID uint   `json:"module_id" binding:"required"`                            // Materi yang dipilih
}

// JoinQuizRequest untuk bergabung ke quiz melalui invite code
type JoinQuizRequest struct {
	InviteCode string `json:"invite_code" binding:"required,len=8"` // Kode undangan
}

// AnswerQuestionRequest untuk menjawab pertanyaan
type AnswerQuestionRequest struct {
	QuestionID   uint   `json:"question_id" binding:"required"`
	UserAnswer   string `json:"user_answer" binding:"required,oneof=A B C D"`
	ResponseTime int32  `json:"response_time" binding:"required,min=0"` // Waktu response dalam milidetik
}

// UpdateQuizSessionRequest untuk update status quiz session
type UpdateQuizSessionRequest struct {
	QuizID        uint   `json:"quiz_id" binding:"required"`
	ParticipantID uint   `json:"participant_id" binding:"required"`
	QuestionID    uint   `json:"question_id" binding:"required"`
	UserAnswer    string `json:"user_answer" binding:"required,oneof=A B C D"`
	ResponseTime  int32  `json:"response_time" binding:"required,min=0"`
}
