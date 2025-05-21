package handlers

import (
	"net/http"
	"strconv"

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
	var request request.CreateQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.quizService.CreateQuiz(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create quiz",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz created successfully",
		"data":    result,
	})
}

func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to get quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to get quiz",
		})
		return
	}

	quiz, err := h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get quiz",
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
			"error":   err.Error(),
			"message": "Failed to get quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quizzes retrieved successfully",
		"data":    quizzes,
	})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	var request request.UpdateQuizRequest
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to update quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to update quiz",
		})
		return
	}

	// Check if quiz exists
	_, err = h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update quiz",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.quizService.UpdateQuiz(uint(id), request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz updated successfully",
		"data":    result,
	})
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to delete quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to delete quiz",
		})
		return
	}

	// Check if quiz exists
	_, err = h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete quiz",
		})
		return
	}

	err = h.quizService.DeleteQuiz(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz deleted successfully",
	})
}
