package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type PrequizHandler struct {
	prequizService services.PrequizServiceInterface
	userService    services.UserServiceInterface
}

func NewPrequizHandler(
	prequizService services.PrequizServiceInterface,
	userService services.UserServiceInterface,
) *PrequizHandler {
	return &PrequizHandler{
		prequizService: prequizService,
		userService:    userService,
	}
}

func (h *PrequizHandler) CreatePrequiz(c *gin.Context) {
	var request request.PrequizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate ModuleID exists (simple validation)
	if request.ModuleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Module ID is required",
			"error":   "Failed to create prequiz",
		})
		return
	}

	result, err := h.prequizService.CreatePrequiz(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create prequiz",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Prequiz created successfully",
		"data":    result,
	})
}

func (h *PrequizHandler) GetPrequizByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prequiz ID"})
		return
	}

	prequiz, err := h.prequizService.GetPrequizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Prequiz not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequiz retrieved successfully",
		"data":    prequiz,
	})
}

func (h *PrequizHandler) GetAllPrequizzes(c *gin.Context) {
	prequizzes, err := h.prequizService.GetAllPrequizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve prequizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequizzes retrieved successfully",
		"data":    prequizzes,
	})
}

func (h *PrequizHandler) UpdatePrequiz(c *gin.Context) {
	var request request.PrequizRequest
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prequiz ID"})
		return
	}

	// Validate ModuleID exists (simple validation)
	if request.ModuleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Module ID is required",
			"error":   "Failed to update prequiz",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	prequiz, err := h.prequizService.UpdatePrequiz(uint(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update prequiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequiz updated successfully",
		"data":    prequiz,
	})
}

func (h *PrequizHandler) GetUserPrequizAnswers(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
			"error":   "Missing user_id parameter",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
		return
	}

	answers, err := h.prequizService.GetUserPrequizAnswers(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user prequiz answers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User prequiz answers retrieved successfully",
		"data":    answers,
	})
}

func (h *PrequizHandler) SubmitPrequizAnswer(c *gin.Context) {
	var request request.PrequizAnswerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
			"error":   "Failed to submit prequiz answer",
		})
		return
	}

	// Convert float64 to uint (JWT stores numbers as float64)
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user ID format",
			"error":   "Failed to submit prequiz answer",
		})
		return
	}
	userID := uint(userIDFloat)

	result, err := h.prequizService.SubmitPrequizAnswer(userID, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to submit prequiz answer",
			"error":   err.Error(),
		})
		return
	}

	// Fetch the prequiz to include explanation and correct answer in the response
	prequiz, _ := h.prequizService.GetPrequizByID(request.PrequizID)

	// Build enriched response payload (keep result fields and add explanation details)
	responsePayload := gin.H{
		"prequiz_id":      result.PrequizID,
		"selected_answer": result.Answer,
		"is_correct":      result.IsCorrect,
	}
	if prequiz != nil {
		responsePayload["question"] = prequiz.Question
		responsePayload["correct_answer"] = prequiz.CorrectAnswer
		responsePayload["explanation"] = prequiz.Explanation
		responsePayload["options"] = prequiz.Options
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Prequiz answer submitted successfully",
		"data":    responsePayload,
	})
}

// GetPrequizzesByModule retrieves all prequizzes for a specific module
func (h *PrequizHandler) GetPrequizzesByModule(c *gin.Context) {
	moduleIDParam := c.Param("module_id")

	moduleID, err := strconv.Atoi(moduleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid module ID",
			"error":   err.Error(),
		})
		return
	}

	// Get user ID from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	// Convert user_id to uint (from float64)
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid user ID format",
		})
		return
	}
	userID := uint(userIDFloat)

	prequizzes, err := h.prequizService.GetPrequizzesByModule(uint(moduleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve prequizzes",
			"error":   err.Error(),
		})
		return
	}

	// Get user's answers to check which prequizzes are already answered
	userAnswers, err := h.prequizService.GetUserPrequizAnswers(userID)
	if err != nil {
		// If error getting user answers, just continue without status
		userAnswers = []models.PrequizUserAnswer{}
	}

	// Create a map for quick lookup of answered prequizzes
	answeredMap := make(map[uint]bool)
	for _, answer := range userAnswers {
		answeredMap[answer.PrequizID] = true
	}

	// Convert prequizzes to response format with status
	prequizzesWithStatus := make([]response.PrequizWithStatus, len(prequizzes))
	answeredCount := 0

	for i, prequiz := range prequizzes {
		isAnswered := answeredMap[prequiz.ID]
		if isAnswered {
			answeredCount++
		}

		prequizzesWithStatus[i] = response.PrequizWithStatus{
			ID:                prequiz.ID,
			CreatedAt:         prequiz.CreatedAt,
			UpdatedAt:         prequiz.UpdatedAt,
			DeletedAt:         prequiz.DeletedAt,
			ModuleID:          prequiz.ModuleID,
			Question:          prequiz.Question,
			Options:           prequiz.Options,
			CorrectAnswer:     prequiz.CorrectAnswer,
			Explanation:       prequiz.Explanation,
			IsAlreadyAnswered: isAnswered,
		}
	}

	// Calculate status
	totalPrequizzes := len(prequizzes)
	isCompleted := answeredCount == totalPrequizzes && totalPrequizzes > 0

	responseData := response.PrequizzesByModuleResponse{
		Success: true,
		Message: "Prequizzes retrieved successfully",
		PrequizStatus: response.PrequizStatus{
			Answered:  answeredCount,
			Total:     totalPrequizzes,
			Completed: isCompleted,
		},
		Prequizzes: prequizzesWithStatus,
	}

	c.JSON(http.StatusOK, responseData)
}
