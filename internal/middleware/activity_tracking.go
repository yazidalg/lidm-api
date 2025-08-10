package middleware

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type ActivityTrackingMiddleware struct {
	activityService services.UserActivityServiceInterface
}

func NewActivityTrackingMiddleware(activityService services.UserActivityServiceInterface) *ActivityTrackingMiddleware {
	return &ActivityTrackingMiddleware{
		activityService: activityService,
	}
}

// TrackActivity middleware that automatically tracks user activities
func (m *ActivityTrackingMiddleware) TrackActivity() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Process the request first
		c.Next()

		// Only track for authenticated users and successful responses
		userID, exists := c.Get("user_id")
		if !exists || c.Writer.Status() >= 400 {
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
				return
			}
		default:
			return
		}

		// Determine activity type based on path and method
		activityType, description, metadata := m.determineActivityType(c)
		if activityType == "" {
			return // Skip tracking for non-relevant endpoints
		}

		// Create activity record
		activity := &models.UserActivity{
			UserID:       uid,
			ActivityType: activityType,
			Description:  description,
			MetaData:     metadata,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.GetHeader("User-Agent"),
			CreatedAt:    time.Now(),
		}

		// Track activity asynchronously to avoid blocking the response
		go func() {
			// Parse metadata back to map
			var metadataMap map[string]interface{}
			if activity.MetaData != "" {
				json.Unmarshal([]byte(activity.MetaData), &metadataMap)
			}

			// Log the activity
			m.activityService.LogActivity(
				activity.UserID,
				activity.ActivityType,
				activity.Description,
				metadataMap,
				activity.IPAddress,
				activity.UserAgent,
			)

			// Update streak for learning activities
			isLearningActivity := activityType == models.ActivityTypeLessonView ||
				activityType == models.ActivityTypeLessonComplete ||
				activityType == models.ActivityTypeModuleView ||
				activityType == models.ActivityTypeModuleComplete ||
				activityType == models.ActivityTypeQuizJoin ||
				activityType == models.ActivityTypeQuizComplete ||
				activityType == models.ActivityTypeQuizAnswer

			if isLearningActivity {
				m.activityService.UpdateUserStreak(activity.UserID)
			}
		}()
	})
}

// determineActivityType determines the activity type based on the request
func (m *ActivityTrackingMiddleware) determineActivityType(c *gin.Context) (string, string, string) {
	path := c.Request.URL.Path
	method := c.Request.Method

	var activityType, description string
	metadata := make(map[string]interface{})

	// Add common metadata
	metadata["path"] = path
	metadata["method"] = method
	metadata["response_status"] = c.Writer.Status()

	switch {
	// Auth activities
	case strings.Contains(path, "/auth/login") && method == "POST":
		activityType = models.ActivityTypeLogin
		description = "Pengguna berhasil masuk"

	case strings.Contains(path, "/auth/google") && method == "POST":
		activityType = models.ActivityTypeLogin
		description = "Pengguna masuk dengan Google"
		metadata["login_method"] = "google"

	case strings.Contains(path, "/auth/belajar-login") && method == "POST":
		activityType = models.ActivityTypeLogin
		description = "Pengguna masuk dengan akun Belajar"
		metadata["login_method"] = "belajar"

	case strings.Contains(path, "/auth/logout") && method == "POST":
		activityType = models.ActivityTypeLogout
		description = "Pengguna keluar"

	// Enhanced Lesson activities with more metadata for RAG
	case strings.Contains(path, "/lesson") && method == "GET":
		activityType = models.ActivityTypeLessonView
		
		if strings.Contains(path, "/lesson/all") {
			description = "Melihat daftar semua pelajaran"
			metadata["action"] = "view_all_lessons"
			metadata["content_type"] = "lesson_list"
			metadata["learning_activity"] = true
			metadata["knowledge_area"] = "lessons"
			
			// For RAG: Add context about browsing lessons
			metadata["user_intent"] = "browse_available_lessons"
			metadata["session_context"] = "lesson_discovery"
		} else if lessonID := c.Param("lesson_id"); lessonID != "" {
			description = "Melihat pelajaran spesifik"
			metadata["lesson_id"] = lessonID
			metadata["action"] = "view_specific_lesson"
			metadata["content_type"] = "lesson_detail"
			metadata["learning_activity"] = true
			metadata["knowledge_area"] = "lesson_content"
			
			// For RAG: Add learning context
			metadata["user_intent"] = "study_lesson"
			metadata["session_context"] = "active_learning"
			metadata["engagement_type"] = "content_consumption"
		}

	case strings.Contains(path, "/lesson") && strings.Contains(path, "/complete") && method == "POST":
		activityType = models.ActivityTypeLessonComplete
		description = "Menyelesaikan pelajaran"
		if lessonID := c.Param("lesson_id"); lessonID != "" {
			metadata["lesson_id"] = lessonID
		}
		metadata["action"] = "complete_lesson"
		metadata["content_type"] = "lesson_completion"
		metadata["learning_activity"] = true
		metadata["achievement"] = true
		metadata["completion_timestamp"] = time.Now().Format(time.RFC3339)
		
		// For RAG: Add completion context
		metadata["user_intent"] = "complete_learning_objective"
		metadata["session_context"] = "lesson_completion"
		metadata["engagement_type"] = "achievement"
		metadata["learning_milestone"] = true

	// Enhanced Module activities with more metadata for RAG
	case strings.Contains(path, "/module") && method == "GET" && !strings.Contains(path, "/modules"):
		activityType = models.ActivityTypeModuleView
		
		if strings.Contains(path, "/module/all") {
			description = "Melihat daftar semua modul"
			metadata["action"] = "view_all_modules"
			metadata["content_type"] = "module_list"
			metadata["learning_activity"] = true
			metadata["knowledge_area"] = "modules"
			
			// For RAG: Add context about browsing modules
			metadata["user_intent"] = "explore_learning_paths"
			metadata["session_context"] = "module_discovery"
			metadata["curriculum_browsing"] = true
			metadata["includes_sub_materials"] = true // Modules include sub materials
		} else if moduleID := c.Param("id"); moduleID != "" {
			description = "Melihat modul spesifik"
			metadata["module_id"] = moduleID
			metadata["action"] = "view_specific_module"
			metadata["content_type"] = "module_detail"
			metadata["learning_activity"] = true
			metadata["knowledge_area"] = "module_content"
			
			// For RAG: Add learning path context
			metadata["user_intent"] = "study_module"
			metadata["session_context"] = "structured_learning"
			metadata["engagement_type"] = "curriculum_following"
			metadata["includes_sub_materials"] = true // Module details include sub materials
			metadata["content_structure"] = "module_with_sub_materials"
		}

	case strings.Contains(path, "/module") && strings.Contains(path, "/complete") && method == "POST":
		activityType = models.ActivityTypeModuleComplete
		description = "Menyelesaikan modul"
		if moduleID := c.Param("module_id"); moduleID != "" {
			metadata["module_id"] = moduleID
		}
		metadata["action"] = "complete_module"
		metadata["content_type"] = "module_completion"
		metadata["learning_activity"] = true
		metadata["achievement"] = true
		metadata["completion_timestamp"] = time.Now().Format(time.RFC3339)
		
		// For RAG: Add major completion context
		metadata["user_intent"] = "complete_learning_module"
		metadata["session_context"] = "module_completion"
		metadata["engagement_type"] = "major_achievement"
		metadata["learning_milestone"] = true
		metadata["curriculum_progress"] = true

	// Enhanced Quiz activities
	case strings.Contains(path, "/quiz") && strings.Contains(path, "/join") && method == "POST":
		activityType = models.ActivityTypeQuizJoin
		description = "Bergabung dengan sesi kuis"
		if quizID := c.Param("quiz_id"); quizID != "" {
			metadata["quiz_id"] = quizID
		}
		metadata["action"] = "join_quiz"
		metadata["content_type"] = "quiz_participation"
		metadata["learning_activity"] = true
		metadata["assessment"] = true
		
		// For RAG: Add assessment context
		metadata["user_intent"] = "test_knowledge"
		metadata["session_context"] = "assessment"
		metadata["engagement_type"] = "knowledge_evaluation"

	case strings.Contains(path, "/quiz") && strings.Contains(path, "/submit") && method == "POST":
		activityType = models.ActivityTypeQuizComplete
		description = "Menyelesaikan kuis"
		if quizID := c.Param("quiz_id"); quizID != "" {
			metadata["quiz_id"] = quizID
		}
		metadata["action"] = "complete_quiz"
		metadata["content_type"] = "quiz_completion"
		metadata["learning_activity"] = true
		metadata["assessment"] = true
		metadata["achievement"] = true
		metadata["completion_timestamp"] = time.Now().Format(time.RFC3339)
		
		// For RAG: Add quiz completion context
		metadata["user_intent"] = "complete_assessment"
		metadata["session_context"] = "quiz_completion"
		metadata["engagement_type"] = "assessment_achievement"
		metadata["learning_milestone"] = true

	case strings.Contains(path, "/quiz") && strings.Contains(path, "/answer") && method == "POST":
		activityType = models.ActivityTypeQuizAnswer
		description = "Menjawab pertanyaan kuis"
		if quizID := c.Param("quiz_id"); quizID != "" {
			metadata["quiz_id"] = quizID
		}
		if questionID := c.Param("question_id"); questionID != "" {
			metadata["question_id"] = questionID
		}
		metadata["action"] = "answer_question"
		metadata["content_type"] = "quiz_interaction"
		metadata["learning_activity"] = true
		metadata["assessment"] = true
		
		// For RAG: Add question answering context
		metadata["user_intent"] = "answer_assessment_question"
		metadata["session_context"] = "active_assessment"
		metadata["engagement_type"] = "knowledge_application"

	// Profile activities
	case strings.Contains(path, "/user") && method == "PUT":
		activityType = models.ActivityTypeProfileUpdate
		description = "Memperbarui informasi profil"
		metadata["action"] = "update_profile"
		metadata["content_type"] = "profile_management"
		
		// For RAG: Add profile context
		metadata["user_intent"] = "manage_profile"
		metadata["session_context"] = "profile_update"

	default:
		// Skip tracking for non-relevant endpoints
		return "", "", ""
	}

	// Convert metadata to JSON string
	metadataJSON, _ := json.Marshal(metadata)

	return activityType, description, string(metadataJSON)
}
