package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type ParticipantHandler struct {
	participantService services.ParticipantServiceInterface
}

func NewParticipantHandler(participantService services.ParticipantServiceInterface) *ParticipantHandler {
	return &ParticipantHandler{participantService}
}

func (h *ParticipantHandler) CreateParticipant(c *gin.Context) {
	var request request.CreateParticipantRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.participantService.CreateParticipant(request)
	if err != nil {
		// Check if this is a "record not found" error
		errMsg := err.Error()
		if errMsg == "record not found" || errMsg == "sql: no rows in result set" ||
			strings.Contains(errMsg, "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   errMsg,
				"message": "The quiz or user you're trying to add does not exist. Please check the quiz_id and user_id values.",
				"details": "Make sure to create the quiz first before adding participants to it.",
			})
			return
		}

		// For other types of errors
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   errMsg,
			"message": "Failed to create participant",
		})
		return
	}

	resultData := response.ParticipantResponse{
		ID:         result.ID,
		UserID:     result.UserID,
		QuizID:     result.QuizID, // Always use the participant's QuizID
		TotalScore: result.TotalScore,
		CreatedAt:  result.CreatedAt,
		User: response.User{
			ID:    result.User.ID,
			Name:  result.User.Name,
			Email: result.User.Email,
			Class: result.User.Class,
		},
		Quiz: response.Quiz{
			ID:     result.QuizID,
			Status: result.Quiz.Status,
		},
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Participant created successfully",
		"data":    resultData,
	})
}

func (h *ParticipantHandler) GetParticipantByID(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to get participant",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to get participant",
		})
		return
	}

	participant, err := h.participantService.GetParticipantByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get participant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant retrieved successfully",
		"data":    participant,
	})
}

func (h *ParticipantHandler) GetAllParticipants(c *gin.Context) {
	participants, err := h.participantService.GetAllParticipants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get participants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participants retrieved successfully",
		"data":    participants,
	})
}

func (h *ParticipantHandler) GetParticipantsByQuizID(c *gin.Context) {
	quizIDParam := c.Param("quiz_id")

	if quizIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Quiz ID is required",
			"message": "Failed to get participants",
		})
		return
	}

	quizID, err := strconv.Atoi(quizIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Quiz ID format",
			"message": "Failed to get participants",
		})
		return
	}

	participants, err := h.participantService.GetParticipantsByQuizID(uint(quizID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get participants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participants retrieved successfully",
		"data":    participants,
	})
}

func (h *ParticipantHandler) GetParticipantsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")

	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID is required",
			"message": "Failed to get participants",
		})
		return
	}

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid User ID format",
			"message": "Failed to get participants",
		})
		return
	}

	participants, err := h.participantService.GetParticipantsByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get participants",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participants retrieved successfully",
		"data":    participants,
	})
}

func (h *ParticipantHandler) UpdateParticipant(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to update participant",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to update participant",
		})
		return
	}

	// Check if participant exists
	_, err = h.participantService.GetParticipantByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update participant",
		})
		return
	}

	// Parse request body
	var request request.UpdateParticipantRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.participantService.UpdateParticipant(int32(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update participant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant updated successfully",
		"data":    result,
	})
}

func (h *ParticipantHandler) DeleteParticipant(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to delete participant",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to delete participant",
		})
		return
	}

	// Check if participant exists
	participant, err := h.participantService.GetParticipantByID(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete participant",
		})
		return
	}

	if participant == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Participant not found",
			"message": "Record not found",
		})
		return
	}

	err = h.participantService.DeleteParticipant(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete participant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Participant deleted successfully",
		"data":    participant,
	})
}
