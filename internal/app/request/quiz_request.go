package request

type CreateQuizRequest struct {
	QuestionsIDs    []uint `json:"questions_ids" binding:"required"`
	ParticipantsIDs []uint `json:"participants_ids"`
	Status          string `json:"status"`
}

type UpdateQuizRequest struct {
	QuestionsIDs []uint `json:"questions_ids"`
	Status       string `json:"status"`
	WinnerID     *uint  `json:"winner_id"`
}
