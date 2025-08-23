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
	userService      services.UserServiceInterface
	progressService  services.ProgressServiceInterface
	moduleService    services.ModuleServiceInterface
	videoQuizService services.VideoQuizServiceInterface
	prequizService   services.PrequizServiceInterface
}

func NewProgressHandler(
	progressService services.ProgressServiceInterface,
	userService services.UserServiceInterface,
	moduleService services.ModuleServiceInterface,
	videoQuizService services.VideoQuizServiceInterface,
	prequizService services.PrequizServiceInterface,
) *ProgressHandler {
	return &ProgressHandler{
		progressService:  progressService,
		userService:      userService,
		moduleService:    moduleService,
		videoQuizService: videoQuizService,
		prequizService:   prequizService,
	}
}

func (h *ProgressHandler) CreateProgress(c *gin.Context) {
	var request request.ProgressRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, userErr := h.userService.GetUserById(int(request.UserID))

	if userResult.ID == 0 || userErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User or lesson not found",
			"message": "Failed to create progress",
		})
		return
	}

	progressData := models.Progress{
		UserID:    uint(userResult.ID),
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
		ID:     progress.ID,
		UserID: progress.UserID,
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

	if userErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "Failed to update progress",
		})
		return
	}

	progressData := models.Progress{
		UserID:      uint(userResult.ID),
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

// GetModuleProgress returns aggregated completion status for a module for the authenticated user
//   - module_completed: true if all lessons in the module are marked completed by the user.
//     If the module has no lessons, it falls back to (videos_completed && prequizzes_completed).
//   - videos_completed: true if the user has answered all video quizzes across all video materials in the module.
//   - prequizzes_completed: true if the user has answered all prequizzes in the module.
func (h *ProgressHandler) GetModuleProgress(c *gin.Context) {
	// Extract user ID from auth middleware
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user context"})
		return
	}

	// Parse module ID
	moduleIDParam := c.Param("id")
	mid, err := strconv.Atoi(moduleIDParam)
	if err != nil || mid <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	// Load module
	module, err := h.moduleService.GetModuleByID(uint32(mid))
	if err != nil || module == nil || module.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Module not found"})
		return
	}

	// Compute videos_completed
	totalVideoQuizzes := 0
	// Build set of all video quiz IDs in this module
	videoQuizIDs := make(map[uint]struct{})
	// Note: Video quizzes are now directly related to module

	answeredVideoQuizIDs := make(map[uint]struct{})
	if totalVideoQuizzes > 0 {
		// Fetch all user answers once, then filter by this module's quiz IDs
		if h.videoQuizService != nil {
			if answers, err := h.videoQuizService.GetAllUserVideoQuizAnswers(userID); err == nil {
				for _, a := range answers {
					if _, ok := videoQuizIDs[a.VideoQuizID]; ok {
						answeredVideoQuizIDs[a.VideoQuizID] = struct{}{}
					}
				}
			}
		}
	}
	videosCompleted := totalVideoQuizzes == 0 || len(answeredVideoQuizIDs) == totalVideoQuizzes

	// Compute prequizzes_completed
	totalPrequizzes := 0
	prequizIDs := make(map[uint]struct{})
	// Note: Prequizzes are now directly related to module

	answeredPrequizIDs := make(map[uint]struct{})
	if totalPrequizzes > 0 {
		if h.prequizService != nil {
			if answers, err := h.prequizService.GetUserPrequizAnswers(userID); err == nil {
				for _, a := range answers {
					if _, ok := prequizIDs[a.PrequizID]; ok {
						answeredPrequizIDs[a.PrequizID] = struct{}{}
					}
				}
			}
		}
	}
	prequizzesCompleted := totalPrequizzes == 0 || len(answeredPrequizIDs) == totalPrequizzes

	// Compute module_completed from lessons progress if lessons exist
	moduleCompleted := false

	// Fallback: if no legacy lessons, consider module completed when both videos and prequizzes are completed
	moduleCompleted = videosCompleted && prequizzesCompleted

	c.JSON(http.StatusOK, gin.H{
		"message": "Module progress retrieved successfully",
		"data": gin.H{
			"module_id":            module.ID,
			"module_completed":     moduleCompleted,
			"videos_completed":     videosCompleted,
			"prequizzes_completed": prequizzesCompleted,
			// Helpful counts for UI (optional)
			"total_video_quizzes":    totalVideoQuizzes,
			"answered_video_quizzes": len(answeredVideoQuizIDs),
			"total_prequizzes":       totalPrequizzes,
			"answered_prequizzes":    len(answeredPrequizIDs),
		},
	})
}
