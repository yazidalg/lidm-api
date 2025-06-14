package handlers

import (
	"net/http"
	"strconv"

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
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to get question",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to get question",
		})
		return
	}

	question, err := h.questionService.GetQuestionByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get question",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question retrieved successfully",
		"data":    question,
	})
}

func (h *QuestionHandler) GetAllQuestions(c *gin.Context) {
	questions, err := h.questionService.GetAllQuestions()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get questions",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Questions retrieved successfully",
		"data":    questions,
	})
}

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	var request request.UpdateQuestionRequest
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to update question",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to update question",
		})
		return
	}

	// Check if question exists
	_, err = h.questionService.GetQuestionByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update question",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.questionService.UpdateQuestion(int32(id), request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update question",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question updated successfully",
		"data":    result,
	})
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to delete question",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to delete question",
		})
		return
	}

	question, err := h.questionService.GetQuestionByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete question",
		})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Question not found",
			"message": "Record not found",
		})
		return
	}

	err = h.questionService.DeleteQuestion(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete question",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Question deleted successfully",
		"data":    question,
	})
}

func (h *QuestionHandler) GetRandomQuestion(c *gin.Context) {
	question, err := h.questionService.GetRandomQuestion(3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get random question",
		})
		return
	}

	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No random question found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Random question retrieved successfully",
		"data":    question,
	})
}
