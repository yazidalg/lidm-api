package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type UserActivityHandler struct {
	activityService services.UserActivityServiceInterface
	moduleService   services.ModuleServiceInterface
}

func NewUserActivityHandler(
	activityService services.UserActivityServiceInterface,
	moduleService services.ModuleServiceInterface,
) *UserActivityHandler {
	return &UserActivityHandler{
		activityService: activityService,
		moduleService:   moduleService,
	}
}

// GetUserActivities - Get activities for a specific user
func (h *UserActivityHandler) GetUserActivities(c *gin.Context) {
	userIDParam := c.Param("user_id")
	limitParam := c.DefaultQuery("limit", "20")

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 20
	}

	activities, err := h.activityService.GetUserActivities(uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user activities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User activities retrieved successfully",
		"data":    activities,
	})
}

// GetLastActivity - Get last activity for a specific user
func (h *UserActivityHandler) GetLastActivity(c *gin.Context) {
	userIDParam := c.Param("user_id")

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	activity, err := h.activityService.GetLastActivity(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No activity found for this user",
		})
		return
	}

	timeSince, _ := h.activityService.GetTimeSinceLastActivity(uint(userID))

	c.JSON(http.StatusOK, gin.H{
		"message":    "Last activity retrieved successfully",
		"data":       activity,
		"time_since": timeSince,
	})
}

// GetRecentActivities - Get recent activities across all users
func (h *UserActivityHandler) GetRecentActivities(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "50")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 50
	}

	activities, err := h.activityService.GetRecentActivities(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch recent activities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recent activities retrieved successfully",
		"data":    activities,
	})
}

// GetMostActiveUsers - Get most active users
func (h *UserActivityHandler) GetMostActiveUsers(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	users, err := h.activityService.GetMostActiveUsers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch most active users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Most active users retrieved successfully",
		"data":    users,
	})
}

// GetMostActiveUsersDetailed - Get most active users with detailed last activities
func (h *UserActivityHandler) GetMostActiveUsersDetailed(c *gin.Context) {
	// Always get only 1 most active user
	mostActiveUsers, err := h.activityService.GetMostActiveUsers(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch most active user",
		})
		return
	}

	if len(mostActiveUsers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Tidak ada pengguna aktif ditemukan",
		})
		return
	}

	mostActiveUser := mostActiveUsers[0]

	// Get all activities for the most active user to calculate learning time
	allActivities, err := h.activityService.GetUserActivities(mostActiveUser.UserID, 1000) // Get many activities
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user activities",
		})
		return
	}

	// Calculate learning statistics
	var totalMinutes int64
	activityBreakdown := make(map[string]int)
	var lastActivity *models.UserActivity

	if len(allActivities) > 0 {
		lastActivity = &allActivities[0] // Most recent activity

		// Calculate total learning time (rough estimation: 5 minutes per lesson/module activity)
		for _, activity := range allActivities {
			activityBreakdown[activity.ActivityType]++

			// Estimate learning time based on activity type
			switch activity.ActivityType {
			case models.ActivityTypeModuleView, models.ActivityTypeModuleComplete:
				totalMinutes += 10 // 10 minutes per module activity
			case models.ActivityTypeQuizJoin, models.ActivityTypeQuizComplete:
				totalMinutes += 3 // 3 minutes per quiz activity
			case models.ActivityTypeQuizAnswer:
				totalMinutes += 1 // 1 minute per quiz answer
			}
		}
	}

	// Get time since last activity
	timeSince, _ := h.activityService.GetTimeSinceLastActivity(mostActiveUser.UserID)

	result := gin.H{
		"user_id":                  mostActiveUser.UserID,
		"username":                 mostActiveUser.Username,
		"total_activities":         mostActiveUser.TotalCount,
		"total_learning_minutes":   totalMinutes,
		"total_learning_hours":     float64(totalMinutes) / 60,
		"activity_breakdown":       activityBreakdown,
		"last_activity":            lastActivity,
		"time_since_last_activity": timeSince,
	}

	// Add streak information
	currentStreak, maxStreak, err := h.activityService.GetUserStreak(mostActiveUser.UserID)
	if err == nil {
		result["current_streak"] = currentStreak
		result["max_streak"] = maxStreak
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengguna paling aktif dengan detail berhasil diambil",
		"data":    result,
	})
}

// LogActivity - Manually log an activity (for admin use)
func (h *UserActivityHandler) LogActivity(c *gin.Context) {
	var request struct {
		UserID       uint                   `json:"user_id" binding:"required"`
		ActivityType string                 `json:"activity_type" binding:"required"`
		Description  string                 `json:"description"`
		Metadata     map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get client IP and User Agent
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err := h.activityService.LogActivity(
		request.UserID,
		request.ActivityType,
		request.Description,
		request.Metadata,
		clientIP,
		userAgent,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to log activity",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Activity logged successfully",
	})
}

// GetMyActivities - Get activities for the current authenticated user
func (h *UserActivityHandler) GetMyActivities(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	limitParam := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 20
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

	activities, err := h.activityService.GetUserActivities(uid, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user activities",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User activities retrieved successfully",
		"data":    activities,
	})
}

// GetMyLastActivity - Get last activity for the current authenticated user
func (h *UserActivityHandler) GetMyLastActivity(c *gin.Context) {
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

	activity, err := h.activityService.GetLastActivity(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch last activity",
		})
		return
	}

	// Get time since last activity
	timeSince, _ := h.activityService.GetTimeSinceLastActivity(uid)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Last activity retrieved successfully",
		"data":       activity,
		"time_since": timeSince,
	})
}

// GetActivityStats - Get activity statistics
func (h *UserActivityHandler) GetActivityStats(c *gin.Context) {
	// This is a basic implementation - you can extend with more complex stats
	recentActivities, err := h.activityService.GetRecentActivities(100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch activity statistics",
		})
		return
	}

	// Calculate basic stats
	activityTypeCount := make(map[string]int)
	for _, activity := range recentActivities {
		activityTypeCount[activity.ActivityType]++
	}

	mostActiveUsers, err := h.activityService.GetMostActiveUsers(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch most active users",
		})
		return
	}

	stats := map[string]interface{}{
		"total_recent_activities": len(recentActivities),
		"activity_type_breakdown": activityTypeCount,
		"most_active_users":       mostActiveUsers,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Activity statistics retrieved successfully",
		"data":    stats,
	})
}

// GetMyStreak - Get streak information for the current authenticated user
func (h *UserActivityHandler) GetMyStreak(c *gin.Context) {
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

	currentStreak, maxStreak, err := h.activityService.GetUserStreak(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch streak information",
		})
		return
	}

	// Get some additional streak stats
	streakData := gin.H{
		"current_streak": currentStreak,
		"max_streak":     maxStreak,
		"streak_status":  getStreakStatus(currentStreak),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Informasi streak berhasil diambil",
		"data":    streakData,
	})
}

// GetActivitiesForRAG - Get enriched activities data for RAG/AI knowledge system
func (h *UserActivityHandler) GetActivitiesForRAG(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "100")
	userIDParam := c.Query("user_id") // Optional: filter by specific user

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 100
	}

	var activities []models.UserActivity

	if userIDParam != "" {
		// Get activities for specific user
		userID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
		activities, err = h.activityService.GetUserActivities(uint(userID), limit)
	} else {
		// Get recent activities from all users
		activities, err = h.activityService.GetRecentActivities(limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch activities for RAG",
		})
		return
	}

	// Enrich activities data for RAG
	enrichedActivities := make([]gin.H, 0, len(activities))

	for _, activity := range activities {
		// Parse metadata
		var metadata map[string]interface{}
		if activity.MetaData != "" {
			json.Unmarshal([]byte(activity.MetaData), &metadata)
		}

		// Create enriched activity object for RAG
		enrichedActivity := gin.H{
			"id":            activity.ID,
			"user_id":       activity.UserID,
			"activity_type": activity.ActivityType,
			"description":   activity.Description,
			"timestamp":     activity.CreatedAt.Format(time.RFC3339),
			"date":          activity.CreatedAt.Format("2006-01-02"),
			"time":          activity.CreatedAt.Format("15:04:05"),
			"metadata":      metadata,
		}

		// Add learning context for RAG
		if isLearningActivity, _ := metadata["learning_activity"].(bool); isLearningActivity {
			enrichedActivity["is_learning_activity"] = true

			// Extract learning insights
			if contentType, ok := metadata["content_type"].(string); ok {
				enrichedActivity["content_type"] = contentType
			}

			if userIntent, ok := metadata["user_intent"].(string); ok {
				enrichedActivity["user_intent"] = userIntent
			}

			if sessionContext, ok := metadata["session_context"].(string); ok {
				enrichedActivity["session_context"] = sessionContext
			}

			if engagementType, ok := metadata["engagement_type"].(string); ok {
				enrichedActivity["engagement_type"] = engagementType
			}

			// Add detailed content information for RAG
			h.addContentDetails(enrichedActivity, metadata, activity.ActivityType)

			// Learning categorization for RAG
			switch activity.ActivityType {
			case models.ActivityTypeModuleView:
				enrichedActivity["learning_category"] = "curriculum_exploration"
				enrichedActivity["structured_learning"] = true
			case models.ActivityTypeModuleComplete:
				enrichedActivity["learning_category"] = "curriculum_completion"
				enrichedActivity["major_achievement"] = true
			case models.ActivityTypeQuizJoin, models.ActivityTypeQuizComplete, models.ActivityTypeQuizAnswer:
				enrichedActivity["learning_category"] = "knowledge_assessment"
				enrichedActivity["skill_evaluation"] = true
			}
		}

		// Add temporal context for learning patterns
		hour := activity.CreatedAt.Hour()
		switch {
		case hour >= 6 && hour < 12:
			enrichedActivity["time_period"] = "morning"
		case hour >= 12 && hour < 18:
			enrichedActivity["time_period"] = "afternoon"
		case hour >= 18 && hour < 22:
			enrichedActivity["time_period"] = "evening"
		default:
			enrichedActivity["time_period"] = "night"
		}

		// Day of week for learning pattern analysis
		enrichedActivity["day_of_week"] = activity.CreatedAt.Weekday().String()
		enrichedActivity["is_weekend"] = activity.CreatedAt.Weekday() == time.Saturday || activity.CreatedAt.Weekday() == time.Sunday

		enrichedActivities = append(enrichedActivities, enrichedActivity)
	}

	// Add summary statistics for RAG context
	stats := gin.H{
		"total_activities": len(enrichedActivities),
		"time_range": gin.H{
			"start": func() string {
				if len(activities) > 0 {
					return activities[len(activities)-1].CreatedAt.Format(time.RFC3339)
				}
				return ""
			}(),
			"end": func() string {
				if len(activities) > 0 {
					return activities[0].CreatedAt.Format(time.RFC3339)
				}
				return ""
			}(),
		},
	}

	// Activity type breakdown for RAG insights
	activityBreakdown := make(map[string]int)
	learningActivityCount := 0

	for _, activity := range activities {
		activityBreakdown[activity.ActivityType]++

		// Check if it's a learning activity
		var metadata map[string]interface{}
		if activity.MetaData != "" {
			json.Unmarshal([]byte(activity.MetaData), &metadata)
			if isLearning, _ := metadata["learning_activity"].(bool); isLearning {
				learningActivityCount++
			}
		}
	}

	stats["activity_breakdown"] = activityBreakdown
	stats["learning_activity_count"] = learningActivityCount
	stats["learning_percentage"] = func() float64 {
		if len(activities) > 0 {
			return float64(learningActivityCount) / float64(len(activities)) * 100
		}
		return 0
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Data aktivitas untuk RAG berhasil diambil",
		"data": gin.H{
			"activities": enrichedActivities,
			"statistics": stats,
			"rag_context": gin.H{
				"purpose":                    "learning_behavior_analysis",
				"data_enrichment":            "enhanced_metadata_for_ai",
				"temporal_analysis":          true,
				"learning_pattern_detection": true,
			},
		},
	})
}

// addContentDetails adds detailed content information for RAG based on activity type and metadata
func (h *UserActivityHandler) addContentDetails(enrichedActivity gin.H, metadata map[string]interface{}, activityType string) {
	switch activityType {

	case models.ActivityTypeModuleView, models.ActivityTypeModuleComplete:
		// Try to get module details
		if moduleIDStr, ok := metadata["module_id"].(string); ok {
			if moduleID, err := strconv.ParseUint(moduleIDStr, 10, 32); err == nil {
				if module, err := h.moduleService.GetModuleByID(uint32(moduleID)); err == nil {
					enrichedActivity["module_details"] = gin.H{
						"id":          module.ID,
						"title":       module.Title,
						"description": module.Description,
						"icon":        module.Icon,
						"thumbnail":   module.Thumbnail,
					}
				}
			}
		} else if activityType == models.ActivityTypeModuleView && metadata["content_type"] == "module_list" {
			// For /module/all endpoint, add summary of all modules
			if modules, err := h.moduleService.GetAllModules(); err == nil {
				moduleSummaries := make([]gin.H, 0, len(modules))
				for _, module := range modules {
					moduleSummaries = append(moduleSummaries, gin.H{
						"id":          module.ID,
						"title":       module.Title,
						"description": module.Description,
					})
				}
				enrichedActivity["available_modules"] = moduleSummaries
				enrichedActivity["total_modules_available"] = len(modules)
			}
		}
	}
}

// Helper function to get streak status message
func getStreakStatus(streak int) string {
	switch {
	case streak == 0:
		return "Belum ada streak, ayo mulai belajar!"
	case streak == 1:
		return "Streak baru dimulai, pertahankan!"
	case streak < 7:
		return fmt.Sprintf("Streak %d hari, tetap semangat!", streak)
	case streak < 30:
		return fmt.Sprintf("Streak %d hari, luar biasa!", streak)
	case streak < 100:
		return fmt.Sprintf("Streak %d hari, amazing!", streak)
	default:
		return fmt.Sprintf("Streak %d hari, legendary!", streak)
	}
}
