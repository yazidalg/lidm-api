package response

import "time"

// QuizResultResponse untuk hasil akhir quiz
type QuizResultResponse struct {
	QuizID             uint                `json:"quiz_id"`
	Mode               string              `json:"mode"`
	Status             string              `json:"status"`
	ModuleName         string              `json:"module_name"`
	TotalQuestions     int                 `json:"total_questions"`
	ParticipantResults []ParticipantResult `json:"participant_results"`
	Winner             *ParticipantResult  `json:"winner,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	CompletedAt        *time.Time          `json:"completed_at,omitempty"`
}

// ParticipantResult untuk hasil individual peserta
type ParticipantResult struct {
	UserID             uint       `json:"user_id"`
	Username           string     `json:"username"`
	TotalScore         int        `json:"total_score"`
	CorrectAnswers     int        `json:"correct_answers"`
	WrongAnswers       int        `json:"wrong_answers"`
	ConsecutiveCorrect int        `json:"consecutive_correct"`
	IsFinished         bool       `json:"is_finished"`
	FinishedAt         *time.Time `json:"finished_at,omitempty"`
}

// QuizSessionResponse untuk real-time quiz session
type QuizSessionResponse struct {
	QuizID          uint                `json:"quiz_id"`
	CurrentQuestion *QuestionDetail     `json:"current_question,omitempty"`
	QuestionNumber  int                 `json:"question_number"`
	TotalQuestions  int                 `json:"total_questions"`
	Participants    []ParticipantStatus `json:"participants"`
	TimeRemaining   int32               `json:"time_remaining,omitempty"` // Waktu tersisa untuk menjawab
	Phase           string              `json:"phase"`                    // "reading", "answering", "result", "finished"
}

// QuestionDetail untuk detail pertanyaan
type QuestionDetail struct {
	ID            uint   `json:"id"`
	Question      string `json:"question"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	ReadTime      int32  `json:"read_time"`
	AnswerTime    int32  `json:"answer_time"`
	Explanation   string `json:"explanation,omitempty"`    // Hanya ditampilkan setelah jawaban
	CorrectAnswer string `json:"correct_answer,omitempty"` // Hanya ditampilkan setelah jawaban
}

// ParticipantStatus untuk status peserta real-time
type ParticipantStatus struct {
	UserID        uint   `json:"user_id"`
	Username      string `json:"username"`
	IsReady       bool   `json:"is_ready"`
	HasAnswered   bool   `json:"has_answered"`
	CurrentScore  int    `json:"current_score"`
	CurrentStreak int    `json:"current_streak"`
	IsFinished    bool   `json:"is_finished"`
}

// JoinQuizResponse untuk response join quiz
type JoinQuizResponse struct {
	QuizID       uint   `json:"quiz_id"`
	InviteCode   string `json:"invite_code"`
	ModuleName   string `json:"module_name"`
	HostUsername string `json:"host_username"`
	Status       string `json:"status"`
	Mode         string `json:"mode"`
	Message      string `json:"message"`
}
