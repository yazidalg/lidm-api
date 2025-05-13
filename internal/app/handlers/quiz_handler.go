package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type QuizHandler struct {
	quizService services.QuizServiceInterface
}

func NewQuizHandler(quizService services.QuizServiceInterface) *QuizHandler {
	return &QuizHandler{quizService}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req request.QuizRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	quiz, err := h.quizService.CreateQuiz(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create quiz",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz created successfully",
		"data":    quiz,
	})
}

func (h *QuizHandler) GetQuizByID(c *gin.Context) {
}
