package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
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
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Quiz ID is required",
		})
		return
	}

	idParam, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid quiz ID",
			"details": err.Error(),
		})
		return
	}

	quiz, err := h.quizService.GetQuizByID(idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Quiz not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz retrieved successfully",
		"data":    quiz,
	})
}

func (h *QuizHandler) GetAllQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetAllQuizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve quizzes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quizzes retrieved successfully",
		"data":    quizzes,
	})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	var req request.QuizRequest

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Quiz ID is required",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid quiz ID",
			"details": err.Error(),
		})
		return
	}

	_, err = h.quizService.GetQuizByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Quiz not found",
			"details": err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	quizModel := models.Quiz{
		Question:      req.Question,
		AnswerTime:    req.AnswerTime,
		ReadTime:      req.ReadTime,
		Options:       models.Options(req.Options),
		CorrectAnswer: req.CorrectAnswer,
		Explanation:   req.Explanation,
		UpdatedAt:     time.Now(),
	}

	quiz, err := h.quizService.UpdateQuiz(int32(id), quizModel)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to update quiz",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz updated successfully",
		"data":    quiz,
	})

}
