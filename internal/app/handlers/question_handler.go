package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type QuestionHandler struct {
	questionService services.QuestionServiceInterface
}

func NewQuestionHandler(questionService services.QuestionServiceInterface) *QuestionHandler {
	return &QuestionHandler{questionService}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var request request.CreateQuestionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.questionService.CreateQuestion(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create question",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Question created successfully",
		"data":    result,
	})
}

func (h *QuestionHandler) GetQuestionByID(c *gin.Context) {

}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {

}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {

}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {

}
