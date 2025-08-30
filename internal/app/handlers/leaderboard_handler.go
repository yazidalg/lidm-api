package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type LeaderboardHandler struct {
	LeaderboardService *services.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService *services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		LeaderboardService: leaderboardService,
	}
}

// GetLeaderboard godoc
// @Summary Get quiz leaderboard
// @Description Get leaderboard for solo quiz or matchmaking quiz with top 3 winners and remaining players
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param module_id query int false "Module ID to filter leaderboard"
// @Param quiz_type query string false "Quiz type: solo or matchmaking"
// @Success 200 {object} map[string]interface{} "juara1, juara2, juara3, leaderboard"
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard [get]
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	moduleIDStr := c.Query("module_id")
	quizType := c.Query("quiz_type") // "solo" or "matchmaking"

	var moduleID *uint
	if moduleIDStr != "" {
		if id, err := strconv.ParseUint(moduleIDStr, 10, 32); err == nil {
			modID := uint(id)
			moduleID = &modID
		}
	}

	// Get current user ID from JWT token
	var currentUserID *uint = nil
	if userIDVal, exists := c.Get("user_id"); exists {
		// user_id is stored as float64 in JWT middleware
		if userIDFloat, ok := userIDVal.(float64); ok {
			userID := uint(userIDFloat)
			currentUserID = &userID
		}
	}

	leaderboard, err := h.LeaderboardService.GetLeaderboard(moduleID, quizType, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get leaderboard",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

// GetUserRank godoc
// @Summary Get user rank in leaderboard
// @Description Get specific user's rank and score in the leaderboard
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param module_id query int false "Module ID to filter leaderboard"
// @Param quiz_type query string false "Quiz type: solo or matchmaking"
// @Success 200 {object} map[string]interface{} "rank, score, user_info"
// @Failure 500 {object} map[string]interface{}
// @Router /leaderboard/user/{user_id} [get]
func (h *LeaderboardHandler) GetUserRank(c *gin.Context) {
	userIDStr := c.Param("user_id")
	moduleIDStr := c.Query("module_id")
	quizType := c.Query("quiz_type")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var moduleID *uint
	if moduleIDStr != "" {
		if id, err := strconv.ParseUint(moduleIDStr, 10, 32); err == nil {
			modID := uint(id)
			moduleID = &modID
		}
	}

	// Get current user ID from JWT token
	var currentUserID *uint = nil
	if userIDVal, exists := c.Get("user_id"); exists {
		// user_id is stored as float64 in JWT middleware
		if userIDFloat, ok := userIDVal.(float64); ok {
			currentUID := uint(userIDFloat)
			currentUserID = &currentUID
		}
	}

	userRank, err := h.LeaderboardService.GetUserRank(uint(userID), moduleID, quizType, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user rank",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, userRank)
}
