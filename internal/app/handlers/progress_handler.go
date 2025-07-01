package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type ProgressHandler struct {
	userService     services.UserServiceInterface
	lessonService   services.LessonServiceInterface
	progressService services.ProgressServiceInterface
}

func NewProgressHandler(
	progressService services.ProgressServiceInterface,
	userService services.UserServiceInterface,
	lessonService services.LessonServiceInterface,
) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		userService:     userService,
		lessonService:   lessonService,
	}
}

func (h *ProgressHandler) CreateProgress(c *gin.Context) {
	var request request.ProgressRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, userErr := h.userService.GetUserById(int(request.UserID))
	lessonResult, lessonErr := h.lessonService.GetLessonByID(uint32(request.LessonID))

	if userResult.ID == 0 || lessonResult == nil || userErr != nil || lessonErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User or lesson not found",
			"message": "Failed to create progress",
		})
		return
	}

	progressData := models.Progress{
		UserID:    uint(userResult.ID),
		LessonID:  uint(lessonResult.ID),
		Completed: false,
	}

	progress, err := h.progressService.CreateProgress(progressData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create progress",
		})
		return
	}

	progressResponse := &response.ProgressResponse{
		ID:       progress.ID,
		UserID:   progress.UserID,
		LessonID: progress.LessonID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress created successfully",
		"data":    progressResponse,
	})
}

func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
	var request request.ProgressRequest
	now := time.Now()
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid progress ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, userErr := h.userService.GetUserById(int(request.UserID))
	lessonResult, lessonErr := h.lessonService.GetLessonByID(uint32(request.LessonID))

	if userResult.ID == 0 || lessonResult == nil || userErr != nil || lessonErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User or lesson not found",
			"message": "Failed to update progress",
		})
		return
	}

	progressData := models.Progress{
		UserID:      uint(userResult.ID),
		LessonID:    uint(lessonResult.ID),
		Completed:   false,
		CompletedAt: &now, // Set CompletedAt to current time
	}

	progress, err := h.progressService.UpdateProgress(uint(id), progressData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update progress",
		})
		return
	}

	progressResponse := &response.ProgressResponse{
		ID:          progress.ID,
		UserID:      progress.UserID,
		LessonID:    progress.LessonID,
		Completed:   progress.Completed,
		CompletedAt: progress.CompletedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress updated successfully",
		"data":    progressResponse,
	})
}

func (h *ProgressHandler) GetAllProgress(c *gin.Context) {
	progresses, err := h.progressService.GetAllProgress()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve progress",
		})
		return
	}

	progressResponses := make([]response.ProgressResponse, 0)

	// 3. Loop through the database results
	for _, p := range progresses {
		// This loop is a new, inner scope.
		// It can access `progressResponses` because it was declared in the parent scope.
		var completedAtStr string
		if p.CompletedAt != nil {
			completedAtStr = p.CompletedAt.Format(time.RFC3339)
		}

		response := response.ProgressResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			LessonID:    p.LessonID,
			Completed:   p.Completed,
			CompletedAt: completedAtStr,
		}

		// Append the formatted response to the slice
		progressResponses = append(progressResponses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress retrieved successfully",
		"data":    progressResponses,
	})
}
