package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type DashboardHandler struct {
	activityService services.UserActivityServiceInterface
	userService     services.UserServiceInterface
}

func NewDashboardHandler(activityService services.UserActivityServiceInterface, userService services.UserServiceInterface) *DashboardHandler {
	return &DashboardHandler{
		activityService: activityService,
		userService:     userService,
	}
}

// GetDashboard - Get dashboard data including most active users and recent activities
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	// Get most active users
	mostActiveUsers, err := h.activityService.GetMostActiveUsers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch most active users",
		})
		return
	}

	// Get recent activities
	recentActivities, err := h.activityService.GetRecentActivities(20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch recent activities",
		})
		return
	}

	// Calculate activity type breakdown
	activityTypeCount := make(map[string]int)
	for _, activity := range recentActivities {
		activityTypeCount[activity.ActivityType]++
	}

	// Prepare dashboard data similar to your image
	dashboardData := map[string]interface{}{
		"most_active_users": mostActiveUsers,
		"recent_activities": recentActivities,
		"activity_stats": map[string]interface{}{
			"total_activities":        len(recentActivities),
			"activity_type_breakdown": activityTypeCount,
		},
		"summary": map[string]interface{}{
			"total_active_users": len(mostActiveUsers),
			"last_updated":       "Just now",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dashboard data retrieved successfully",
		"data":    dashboardData,
	})
}

// GetUserDashboard - Get dashboard data for a specific user
func (h *DashboardHandler) GetUserDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Convert userID to uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get user's last activity
	lastActivity, err := h.activityService.GetLastActivity(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch last activity",
		})
		return
	}

	// Get time since last activity
	timeSince, _ := h.activityService.GetTimeSinceLastActivity(uid)

	// Get user's recent activities
	userActivities, err := h.activityService.GetUserActivities(uid, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user activities",
		})
		return
	}

	// Calculate user activity stats
	activityTypeCount := make(map[string]int)
	for _, activity := range userActivities {
		activityTypeCount[activity.ActivityType]++
	}

	userDashboard := map[string]interface{}{
		"last_activity": map[string]interface{}{
			"activity":   lastActivity,
			"time_since": timeSince,
		},
		"recent_activities": userActivities,
		"activity_stats": map[string]interface{}{
			"total_activities":        len(userActivities),
			"activity_type_breakdown": activityTypeCount,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User dashboard data retrieved successfully",
		"data":    userDashboard,
	})
}
