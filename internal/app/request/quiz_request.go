package request

type CreateQuizRequest struct {
	ParticipantID uint   `json:"participant_id" binding:"required"`
	QuestionsIDs  []uint `json:"questions_ids" binding:"required"`
	AnswersIDs    []uint `json:"answers_ids"`
	Status        string `json:"status"`
}

type UpdateQuizRequest struct {
	ParticipantID uint   `json:"participant_id"`
	QuestionsIDs  []uint `json:"questions_ids"`
	AnswersIDs    []uint `json:"answers_ids"`
	Status        string `json:"status"`
}
