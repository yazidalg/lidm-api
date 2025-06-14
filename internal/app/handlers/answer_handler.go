package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type AnswerHandler struct {
	answerService services.AnswerServiceInterface
}

func NewAnswerHandler(answerService services.AnswerServiceInterface) *AnswerHandler {
	return &AnswerHandler{answerService}
}

func (h *AnswerHandler) CreateAnswer(c *gin.Context) {
	var request request.CreateAnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.answerService.CreateAnswer(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create answer",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Answer created successfully",
		"data":    result,
	})
}

func (h *AnswerHandler) GetAnswerByID(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to get answer",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to get answer",
		})
		return
	}

	answer, err := h.answerService.GetAnswerByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get answer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer retrieved successfully",
		"data":    answer,
	})
}

func (h *AnswerHandler) GetAllAnswers(c *gin.Context) {
	answers, err := h.answerService.GetAllAnswers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get answers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answers retrieved successfully",
		"data":    answers,
	})
}

func (h *AnswerHandler) UpdateAnswer(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to update answer",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to update answer",
		})
		return
	}

	// Check if answer exists
	existingAnswer, err := h.answerService.GetAnswerByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update answer",
		})
		return
	}

	// Parse request body
	var request request.CreateAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to model
	existingAnswer.QuestionID = request.QuestionID
	existingAnswer.OptionSelected = request.OptionSelected
	existingAnswer.IsCorrect = request.IsCorrect
	existingAnswer.Score = request.Score

	result, err := h.answerService.UpdateAnswer(int32(id), existingAnswer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update answer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer updated successfully",
		"data":    result,
	})
}

func (h *AnswerHandler) DeleteAnswer(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to delete answer",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to delete answer",
		})
		return
	}

	// Check if answer exists
	answer, err := h.answerService.GetAnswerByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete answer",
		})
		return
	}

	if answer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Answer not found",
			"message": "Record not found",
		})
		return
	}

	err = h.answerService.DeleteAnswer(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete answer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Answer deleted successfully",
		"data":    answer,
	})
}
