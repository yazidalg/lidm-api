package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type QuizSessionHandler struct {
	quizSessionService services.QuizSessionServiceInterface
}

func NewQuizSessionHandler(quizSessionService services.QuizSessionServiceInterface) *QuizSessionHandler {
	return &QuizSessionHandler{quizSessionService}
}

// CreateQuizSession creates a new quiz session (single player or multiplayer)
func (h *QuizSessionHandler) CreateQuizSession(c *gin.Context) {
	var req request.CreateQuizSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	// Get user from context (from middleware)
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	result, err := h.quizSessionService.CreateQuizSession(user.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create quiz session",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz session created successfully",
		"data":    result,
	})
}

// JoinQuiz allows a user to join a quiz using invite code
func (h *QuizSessionHandler) JoinQuiz(c *gin.Context) {
	var req request.JoinQuizRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	// Get user from context
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	result, err := h.quizSessionService.JoinQuiz(user.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to join quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined quiz",
		"data":    result,
	})
}

// AnswerQuestion allows a user to submit an answer to a question
func (h *QuizSessionHandler) AnswerQuestion(c *gin.Context) {
	var req request.AnswerQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request format",
		})
		return
	}

	// Get user from context
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	result, err := h.quizSessionService.AnswerQuestion(user.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to submit answer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer submitted successfully",
		"data":    result,
	})
}

// GetQuizSession retrieves current quiz session state
func (h *QuizSessionHandler) GetQuizSession(c *gin.Context) {
	quizIDParam := c.Param("quiz_id")

	if quizIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Quiz ID is required",
			"message": "Failed to get quiz session",
		})
		return
	}

	quizID, err := strconv.ParseUint(quizIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid quiz ID format",
			"message": "Failed to get quiz session",
		})
		return
	}

	result, err := h.quizSessionService.GetQuizSession(uint(quizID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get quiz session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz session retrieved successfully",
		"data":    result,
	})
}

// GetQuizResult retrieves quiz results
func (h *QuizSessionHandler) GetQuizResult(c *gin.Context) {
	quizIDParam := c.Param("quiz_id")

	if quizIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Quiz ID is required",
			"message": "Failed to get quiz result",
		})
		return
	}

	quizID, err := strconv.ParseUint(quizIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid quiz ID format",
			"message": "Failed to get quiz result",
		})
		return
	}

	result, err := h.quizSessionService.GetQuizResult(uint(quizID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get quiz result",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz result retrieved successfully",
		"data":    result,
	})
}

// FinishQuiz marks a participant as finished
func (h *QuizSessionHandler) FinishQuiz(c *gin.Context) {
	quizIDParam := c.Param("quiz_id")

	if quizIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Quiz ID is required",
			"message": "Failed to finish quiz",
		})
		return
	}

	quizID, err := strconv.ParseUint(quizIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid quiz ID format",
			"message": "Failed to finish quiz",
		})
		return
	}

	// Get user from context
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	err = h.quizSessionService.FinishQuiz(uint(quizID), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to finish quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz finished successfully",
	})
}
