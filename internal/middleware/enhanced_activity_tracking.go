package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

// ContentEnricher interface untuk mengambil data detail dari content
type ContentEnricher interface {
	GetModuleDetails(moduleID uint) (*ModuleContent, error)
	GetQuizDetails(quizID uint) (*QuizContent, error)
}

type ModuleContent struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	Difficulty    string   `json:"difficulty"`
	TotalDuration int      `json:"total_duration_minutes"`
	Tags          []string `json:"tags"`
	LearningPath  string   `json:"learning_path"`
}

type QuizContent struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	ModuleID      uint     `json:"module_id"`
	ModuleName    string   `json:"module_name"`
	QuestionCount int      `json:"question_count"`
	Duration      int      `json:"duration_minutes"`
	Difficulty    string   `json:"difficulty"`
	Topics        []string `json:"topics"`
}

type EnhancedActivityTrackingMiddleware struct {
	activityService services.UserActivityServiceInterface
	contentEnricher ContentEnricher
}

func NewEnhancedActivityTrackingMiddleware(activityService services.UserActivityServiceInterface, contentEnricher ContentEnricher) *EnhancedActivityTrackingMiddleware {
	return &EnhancedActivityTrackingMiddleware{
		activityService: activityService,
		contentEnricher: contentEnricher,
	}
}

// Enhanced activity tracking with detailed content metadata
func (m *EnhancedActivityTrackingMiddleware) TrackActivityWithEnhancedMetadata() gin.HandlerFunc {
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

		// Determine activity type and enrich with enhanced metadata
		activityType, description, metadata := m.determineActivityWithContent(c)
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
			isLearningActivity := activityType == models.ActivityTypeModuleView ||
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

// determineActivityWithContent - Enhanced version with content details
func (m *EnhancedActivityTrackingMiddleware) determineActivityWithContent(c *gin.Context) (string, string, string) {
	path := c.Request.URL.Path
	method := c.Request.Method

	var activityType, description string
	metadata := make(map[string]interface{})

	// Add common metadata
	metadata["path"] = path
	metadata["method"] = method
	metadata["response_status"] = c.Writer.Status()
	metadata["timestamp"] = time.Now().Format(time.RFC3339)

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

	// Enhanced Module activities with content details
	case strings.Contains(path, "/module") && method == "GET" && !strings.Contains(path, "/modules"):
		activityType = models.ActivityTypeModuleView

		if strings.Contains(path, "/module/all") {
			description = "Melihat daftar semua modul"
			metadata["action"] = "view_all_modules"
		} else if moduleID := c.Param("module_id"); moduleID != "" {
			if id, err := strconv.ParseUint(moduleID, 10, 32); err == nil {
				if moduleContent, err := m.contentEnricher.GetModuleDetails(uint(id)); err == nil {
					description = fmt.Sprintf("Melihat modul: %s", moduleContent.Title)
					metadata["module_id"] = moduleID
					metadata["module_content"] = moduleContent
					metadata["action"] = "view_specific_module"
				}
			}
		}

	case strings.Contains(path, "/module") && strings.Contains(path, "/complete") && method == "POST":
		activityType = models.ActivityTypeModuleComplete
		if moduleID := c.Param("module_id"); moduleID != "" {
			if id, err := strconv.ParseUint(moduleID, 10, 32); err == nil {
				if moduleContent, err := m.contentEnricher.GetModuleDetails(uint(id)); err == nil {
					description = fmt.Sprintf("Menyelesaikan modul: %s", moduleContent.Title)
					metadata["module_id"] = moduleID
					metadata["module_content"] = moduleContent
					metadata["action"] = "complete_module"
					metadata["completion_timestamp"] = time.Now().Format(time.RFC3339)
				}
			}
		}

	// Enhanced Quiz activities
	case strings.Contains(path, "/quiz") && strings.Contains(path, "/join") && method == "POST":
		activityType = models.ActivityTypeQuizJoin
		if quizID := c.Param("quiz_id"); quizID != "" {
			if id, err := strconv.ParseUint(quizID, 10, 32); err == nil {
				if quizContent, err := m.contentEnricher.GetQuizDetails(uint(id)); err == nil {
					description = fmt.Sprintf("Bergabung dengan kuis: %s", quizContent.Title)
					metadata["quiz_id"] = quizID
					metadata["quiz_content"] = quizContent
					metadata["action"] = "join_quiz"
				}
			}
		}

	case strings.Contains(path, "/quiz") && strings.Contains(path, "/submit") && method == "POST":
		activityType = models.ActivityTypeQuizComplete
		if quizID := c.Param("quiz_id"); quizID != "" {
			if id, err := strconv.ParseUint(quizID, 10, 32); err == nil {
				if quizContent, err := m.contentEnricher.GetQuizDetails(uint(id)); err == nil {
					description = fmt.Sprintf("Menyelesaikan kuis: %s", quizContent.Title)
					metadata["quiz_id"] = quizID
					metadata["quiz_content"] = quizContent
					metadata["action"] = "complete_quiz"
					metadata["completion_timestamp"] = time.Now().Format(time.RFC3339)
				}
			}
		}

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

	// Profile activities
	case strings.Contains(path, "/user") && method == "PUT":
		activityType = models.ActivityTypeProfileUpdate
		description = "Memperbarui informasi profil"
		metadata["action"] = "update_profile"

	default:
		// Skip tracking for non-relevant endpoints
		return "", "", ""
	}

	// Convert metadata to JSON string
	metadataJSON, _ := json.Marshal(metadata)

	return activityType, description, string(metadataJSON)
}
