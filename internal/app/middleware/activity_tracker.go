package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

// ActivityTracker middleware to automatically log user activities
func ActivityTracker(activityService services.UserActivityServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continue with the request
		c.Next()

		// Log activity after the request is processed
		userID, exists := c.Get("user_id")
		if !exists {
			return // Skip if no user authenticated
		}

		uid, ok := userID.(uint)
		if !ok {
			return // Skip if user_id is not uint
		}

		// Determine activity type based on the route
		activityType := determineActivityType(c.Request.Method, c.FullPath())
		if activityType == "" {
			return // Skip if not a trackable activity
		}

		// Get client information
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Create metadata
		metadata := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.FullPath(),
			"status":     c.Writer.Status(),
			"params":     c.Params,
			"user_agent": userAgent,
		}

		// Log the activity (don't block on errors)
		go func() {
			_ = activityService.LogActivity(
				uid,
				activityType,
				"",
				metadata,
				clientIP,
				userAgent,
			)
		}()
	}
}

// Helper function to determine activity type based on route
func determineActivityType(method, path string) string {
	switch {
	case method == "POST" && path == "/auth/login":
		return models.ActivityTypeLogin
	case method == "POST" && path == "/auth/logout":
		return models.ActivityTypeLogout
	case method == "POST" && path == "/quiz/:id/join":
		return models.ActivityTypeQuizJoin
	case method == "POST" && path == "/quiz/:id/finish":
		return models.ActivityTypeQuizComplete
	case method == "POST" && path == "/quiz/answer":
		return models.ActivityTypeQuizAnswer
	case method == "GET" && path == "/lesson/:id":
		return models.ActivityTypeLessonView
	case method == "POST" && path == "/progress":
		return models.ActivityTypeLessonComplete
	case method == "GET" && path == "/module/:id":
		return models.ActivityTypeModuleView
	case method == "PUT" && path == "/user/profile":
		return models.ActivityTypeProfileUpdate
	default:
		return "" // No tracking for this route
	}
}
