package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type VideoQuizHandler struct {
	videoQuizService services.VideoQuizServiceInterface
}

func NewVideoQuizHandler(videoQuizService services.VideoQuizServiceInterface) *VideoQuizHandler {
	return &VideoQuizHandler{
		videoQuizService: videoQuizService,
	}
}

func (h *VideoQuizHandler) CreateVideoQuiz(c *gin.Context) {
	var request request.VideoQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.videoQuizService.CreateVideoQuiz(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create video quiz",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Video quiz created successfully",
		"data":    result,
	})
}

func (h *VideoQuizHandler) GetVideoQuizByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video quiz ID"})
		return
	}

	videoQuiz, err := h.videoQuizService.GetVideoQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Video quiz not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video quiz retrieved successfully",
		"data":    videoQuiz,
	})
}

func (h *VideoQuizHandler) GetVideoQuizzesByVideoMaterial(c *gin.Context) {
	videoMaterialIDParam := c.Param("video_material_id")
	videoMaterialID, err := strconv.Atoi(videoMaterialIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video material ID"})
		return
	}

	videoQuizzes, err := h.videoQuizService.GetVideoQuizzesByVideoMaterialID(uint(videoMaterialID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve video quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video quizzes retrieved successfully",
		"data":    videoQuizzes,
	})
}

func (h *VideoQuizHandler) UpdateVideoQuiz(c *gin.Context) {
	var request request.VideoQuizRequest
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video quiz ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	videoQuiz, err := h.videoQuizService.UpdateVideoQuiz(uint(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update video quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video quiz updated successfully",
		"data":    videoQuiz,
	})
}

func (h *VideoQuizHandler) DeleteVideoQuiz(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video quiz ID"})
		return
	}

	err = h.videoQuizService.DeleteVideoQuiz(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete video quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video quiz deleted successfully",
	})
}

func (h *VideoQuizHandler) SubmitVideoQuizAnswer(c *gin.Context) {
	var request request.VideoQuizAnswerRequest

	// Get user ID from context (assuming it's set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	// Convert float64 to uint (JWT stores numbers as float64)
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid user ID",
		})
		return
	}
	userID := uint(userIDFloat)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.videoQuizService.SubmitVideoQuizAnswer(userID, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to submit video quiz answer",
			"error":   err.Error(),
		})
		return
	}

	// Fetch quiz details to include explanation in the response
	videoQuiz, _ := h.videoQuizService.GetVideoQuizByID(request.VideoQuizID)

	responsePayload := gin.H{
		"video_quiz_id":   result.VideoQuizID,
		"selected_answer": result.SelectedAnswer,
		"is_correct":      result.IsCorrect,
		"answered_at":     result.AnsweredAt,
		"response_time":   result.ResponseTime,
	}
	if videoQuiz != nil {
		responsePayload["question"] = videoQuiz.Question
		responsePayload["correct_answer"] = videoQuiz.CorrectAnswer
		responsePayload["explanation"] = videoQuiz.Explanation
		responsePayload["options"] = videoQuiz.Options
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Video quiz answer submitted successfully",
		"data":    responsePayload,
	})
}

func (h *VideoQuizHandler) GetUserVideoQuizAnswers(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	videoMaterialIDParam := c.Param("video_material_id")
	videoMaterialID, err := strconv.Atoi(videoMaterialIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video material ID"})
		return
	}

	userAnswers, err := h.videoQuizService.GetUserVideoQuizAnswers(userID, uint(videoMaterialID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve user video quiz answers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User video quiz answers retrieved successfully",
		"data":    userAnswers,
	})
}

func (h *VideoQuizHandler) GetAllUserVideoQuizAnswers(c *gin.Context) {
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

	answers, err := h.videoQuizService.GetAllUserVideoQuizAnswers(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user video quiz answers",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User video quiz answers retrieved successfully",
		"data":    answers,
	})
}
