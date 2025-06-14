package request

type CreateQuizRequest struct {
	QuestionsIDs []uint `json:"questions_ids" binding:"required"`
	Status       string `json:"status"`
}

type UpdateQuizRequest struct {
	QuestionsIDs []uint `json:"questions_ids"`
	Status       string `json:"status"`
}
