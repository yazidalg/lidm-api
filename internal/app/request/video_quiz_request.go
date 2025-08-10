package request

type VideoQuizRequest struct {
	VideoMaterialID uint            `json:"video_material_id" binding:"required"`
	Question        string          `json:"question" binding:"required"`
	TimestampStart  int             `json:"timestamp_start" binding:"required"` // Detik ke berapa quiz muncul
	TimestampEnd    int             `json:"timestamp_end" binding:"required"`   // Detik ke berapa quiz berakhir
	Options         QuestionOptions `json:"options" binding:"required"`
	CorrectAnswer   string          `json:"correct_answer" binding:"required"`
	Explanation     string          `json:"explanation" binding:"required"`
	Order           int             `json:"order"`
}

type VideoQuizAnswerRequest struct {
	VideoQuizID    uint   `json:"video_quiz_id" binding:"required"`
	SelectedAnswer string `json:"selected_answer" binding:"required"` // A, B, C, atau D
	ResponseTime   int    `json:"response_time" binding:"required"`   // Waktu respons dalam detik
}
