package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
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

func (h *ProgressHandler) UpdateProgress(c *gin.Context) {
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
			"message": "Failed to update progress",
		})
		return
	}

	// Create a progress request object
	request.UserID = uint32(userResult.ID)
	request.LessonID = uint32(lessonResult.ID)

	progress, err := h.progressService.UpdateProgress(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update progress",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress updated successfully",
		"data":    progress,
	})
}

func (h *ProgressHandler) GetAllProgress(c *gin.Context) {
	progress, err := h.progressService.GetAllProgress()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve progress",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Progress retrieved successfully",
		"data":    progress,
	})
}
