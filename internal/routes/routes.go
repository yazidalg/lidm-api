package routes

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/middleware"
)

func NewRoute(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	forgotPasswordHandler *handlers.ForgotPasswordHandler,
	questionHandler *handlers.QuestionHandler,
	answerHandler *handlers.AnswerHandler,
	participantHandler *handlers.ParticipantHandler,
	quizHandler *handlers.QuizHandler,
	quizSessionHandler *handlers.QuizSessionHandler,
	socketHandler *handlers.SocketHandler,
	moduleHandler *handlers.ModuleHandler,
	progressHandler *handlers.ProgressHandler,
	prequizHandler *handlers.PrequizHandler,
	videoQuizHandler *handlers.VideoQuizHandler,
	authMiddleware *middleware.AuthMiddleware,
	roleHandler *handlers.RoleHandler,
	activityHandler *handlers.UserActivityHandler,
	activityService services.UserActivityServiceInterface,
	dashboardHandler *handlers.DashboardHandler,
	leaderboardHandler *handlers.LeaderboardHandler,
	flashcardHandler *handlers.FlashcardHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	router := gin.Default()

	// Enable CORS - configurable for different environments
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Build allowed origins list from environment variable and defaults
		allowedOrigins := []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"https://lidm-frontend-629808488591.asia-southeast2.run.app",
		}

		// Add production frontend domain
		if prodFrontend := os.Getenv("FRONTEND_URL"); prodFrontend != "" {
			allowedOrigins = append(allowedOrigins, prodFrontend)
		}

		// Also support comma-separated list of origins
		if allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS"); allowedOriginsEnv != "" {
			envOrigins := strings.Split(allowedOriginsEnv, ",")
			for _, envOrigin := range envOrigins {
				trimmed := strings.TrimSpace(envOrigin)
				if trimmed != "" {
					allowedOrigins = append(allowedOrigins, trimmed)
				}
			}
		}

		// Allow requests from allowed origins or if no origin (like Postman)
		allowOrigin := ""
		if origin == "" {
			allowOrigin = "*"
		} else {
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					allowOrigin = origin
					break
				}
			}
		}

		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Content-Length")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize activity tracking middleware
	activityMiddleware := middleware.NewActivityTrackingMiddleware(activityService)

	// Static file serving for uploads
	router.Static("/uploads", "./uploads")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	// Health check routes (public - no auth required)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/healthy", healthHandler.Healthy)

	// Public routes (tidak perlu auth)
	authGroupHandler := router.Group("auth")
	authGroupHandler.Use(activityMiddleware.TrackActivity()) // Track auth activities
	{
		authGroupHandler.POST("/register", authHandler.RegisterUser)
		authGroupHandler.POST("/login", authHandler.LoginUser)
		authGroupHandler.POST("/google", authHandler.GoogleLogin)
		authGroupHandler.GET("/verify/:verificationToken", authHandler.VerifyEmail)
		authGroupHandler.POST("/belajar-login", authHandler.BelajarLogin)
		authGroupHandler.POST("/logout", authHandler.Logout)
	}

	// Forgot password routes (public)
	forgotPasswordGroup := router.Group("password")
	{
		forgotPasswordGroup.POST("/forgot", forgotPasswordHandler.RequestPasswordReset)
		forgotPasswordGroup.POST("/verify-otp", forgotPasswordHandler.VerifyOTP)
		forgotPasswordGroup.POST("/reset", forgotPasswordHandler.ResetPassword)
	}

	// User routes (authenticated users only)
	userGroupHandler := router.Group("user")
	userGroupHandler.Use(authMiddleware.RequireAuth)
	userGroupHandler.Use(activityMiddleware.TrackActivity()) // Track user activities
	{
		userGroupHandler.GET("/profile", userHandler.GetUserById)
		userGroupHandler.GET("/admin", userHandler.GetUserAdmin)
	}

	// Question routes - Admin only untuk CUD, User bisa Read
	questionGroupHandler := router.Group("question")
	questionGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes
		questionAdminGroup := questionGroupHandler.Group("")
		questionAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			questionAdminGroup.POST("/create", questionHandler.CreateQuestion)
			questionAdminGroup.PUT("/:id", questionHandler.UpdateQuestion)
			questionAdminGroup.DELETE("/:id", questionHandler.DeleteQuestion)
		}

		// User accessible routes
		questionGroupHandler.GET("/:id", questionHandler.GetQuestionByID)
		questionGroupHandler.GET("/all", questionHandler.GetAllQuestions)
		questionGroupHandler.GET("/random", questionHandler.GetRandomQuestion)
	}

	// Answer routes - Admin only untuk CUD, User bisa Read
	answerGroupHandler := router.Group("answer")
	answerGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes
		answerAdminGroup := answerGroupHandler.Group("")
		answerAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			answerAdminGroup.POST("/create", answerHandler.CreateAnswer)
			answerAdminGroup.PUT("/:id", answerHandler.UpdateAnswer)
			answerAdminGroup.DELETE("/:id", answerHandler.DeleteAnswer)
		}

		// User accessible routes
		answerGroupHandler.GET("/:id", answerHandler.GetAnswerByID)
		answerGroupHandler.GET("/all", answerHandler.GetAllAnswers)
	}

	// Participant routes - User bisa create untuk diri sendiri, Admin bisa semua
	participantGroupHandler := router.Group("participant")
	participantGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes
		participantAdminGroup := participantGroupHandler.Group("")
		participantAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			participantAdminGroup.GET("/all", participantHandler.GetAllParticipants)
			participantAdminGroup.PUT("/:id", participantHandler.UpdateParticipant)
			participantAdminGroup.DELETE("/:id", participantHandler.DeleteParticipant)
		}

		// User accessible routes
		participantGroupHandler.POST("/create", participantHandler.CreateParticipant)
		participantGroupHandler.GET("/:id", participantHandler.GetParticipantByID)
		participantGroupHandler.GET("/quiz/:quiz_id", participantHandler.GetParticipantsByQuizID)
		participantGroupHandler.GET("/user/:user_id", participantHandler.GetParticipantsByUserID)
	}

	// Quiz routes - Admin only untuk CUD, User bisa Read
	quizGroupHandler := router.Group("quiz")
	quizGroupHandler.Use(authMiddleware.RequireAuth)
	quizGroupHandler.Use(activityMiddleware.TrackActivity()) // Track quiz activities
	{
		// Admin and Teacher routes
		quizAdminGroup := quizGroupHandler.Group("")
		quizAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			quizAdminGroup.PUT("/:id", quizHandler.UpdateQuiz)
			quizAdminGroup.DELETE("/:id", quizHandler.DeleteQuiz)
		}

		// User accessible routes
		quizGroupHandler.GET("/:id", quizHandler.GetQuizByID)
		quizGroupHandler.GET("/all", quizHandler.GetAllQuizzes)
		quizGroupHandler.GET("/module/:module_id", quizHandler.GetQuizzesByModule)
		quizGroupHandler.POST("/create", quizHandler.CreateQuizLobby)
		quizGroupHandler.POST("/join", quizHandler.JoinQuizLobby)
	}

	// Quiz Session routes - User bisa create dan join quiz session
	quizSessionGroupHandler := router.Group("quiz-sessions")
	quizSessionGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// User accessible routes untuk quiz session
		quizSessionGroupHandler.POST("/", quizSessionHandler.CreateQuizSession)
		quizSessionGroupHandler.POST("/join", quizSessionHandler.JoinQuiz)
		quizSessionGroupHandler.POST("/answer", quizSessionHandler.AnswerQuestion)
		quizSessionGroupHandler.GET("/:quiz_id", quizSessionHandler.GetQuizSession)
		quizSessionGroupHandler.GET("/:quiz_id/results", quizSessionHandler.GetQuizResult)
		quizSessionGroupHandler.POST("/:quiz_id/finish", quizSessionHandler.FinishQuiz)
	}

	// Socket routes - Semua authenticated user bisa akses
	socketGroupHandler := router.Group("ws")
	socketGroupHandler.Use(authMiddleware.RequireAuth)
	{
		socketGroupHandler.GET("/:roomName", socketHandler.ServeWs)
		socketGroupHandler.GET("/matchmaking", socketHandler.MatchMaking)
		socketGroupHandler.GET("/prequiz", socketHandler.PreQuiz)
		socketGroupHandler.GET("/quiz-session/:quiz_id", socketHandler.QuizSession)
	}

	// Module routes - Admin only untuk CUD, User bisa Read
	moduleGroupHandler := router.Group("module")
	moduleGroupHandler.Use(authMiddleware.RequireAuth)
	moduleGroupHandler.Use(activityMiddleware.TrackActivity()) // Track module activities
	{
		// Admin and Teacher routes
		moduleAdminGroup := moduleGroupHandler.Group("admin")
		moduleAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			moduleAdminGroup.GET("/all", moduleHandler.GetAllModulesAdmin)
			moduleAdminGroup.GET("/all-no-pagination", moduleHandler.GetAllModulesAdminAll)
		}

		// Admin and Teacher routes (original admin group)
		moduleAdminGroupOriginal := moduleGroupHandler.Group("")
		moduleAdminGroupOriginal.Use(authMiddleware.RequireAdminOrTeacher)
		{
			moduleAdminGroupOriginal.POST("/create", moduleHandler.CreateModule)
			moduleAdminGroupOriginal.POST("/create-with-video", moduleHandler.CreateModuleWithVideo)
			moduleAdminGroupOriginal.PUT("/:id", moduleHandler.UpdateModule)
			moduleAdminGroupOriginal.PUT("/:id/with-video", moduleHandler.UpdateModuleWithVideo)
			moduleAdminGroupOriginal.POST("/ar-experiment", moduleHandler.AddARExperimentToModule)
			moduleAdminGroupOriginal.DELETE("/:id", moduleHandler.DeleteModule)
			moduleAdminGroupOriginal.POST("/:id/upload-icon", moduleHandler.UploadModuleIcon)
			moduleAdminGroupOriginal.DELETE("/:id/delete-icon", moduleHandler.DeleteModuleIcon)
		}

		// User accessible routes (register more specific pattern before generic one)
		moduleGroupHandler.GET("/:id/progress", progressHandler.GetModuleProgress)
		moduleGroupHandler.GET("/:id", moduleHandler.GetModuleByID)
		moduleGroupHandler.GET("/:id/admin", moduleHandler.GetModuleByIdAdmin) // Use unlock-aware endpoint for authenticated users
		moduleGroupHandler.GET("/all", moduleHandler.GetAllModulesWithUnlockStatus)
		// Keep legacy endpoint for backward compatibility
		moduleGroupHandler.GET("/all-legacy", moduleHandler.GetAllModulesWithProgress)
	}

	// Progress routes - User untuk data sendiri, Admin untuk semua
	progressGroupHandler := router.Group("progress")
	progressGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes
		progressAdminGroup := progressGroupHandler.Group("")
		progressAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			progressAdminGroup.GET("/all", progressHandler.GetAllProgress)
		}

		// User accessible routes
		progressGroupHandler.POST("/create", progressHandler.CreateProgress)
		progressGroupHandler.PUT("/:id", progressHandler.UpdateProgress)
		// Alias endpoint to fetch module progress via progress namespace
		progressGroupHandler.GET("/module/:id", progressHandler.GetModuleProgress)
	}

	// Prequiz routes - Admin only untuk CUD, User bisa Read
	prequizGroupHandler := router.Group("prequiz")
	prequizGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes
		prequizAdminGroup := prequizGroupHandler.Group("")
		prequizAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			prequizAdminGroup.POST("/create", prequizHandler.CreatePrequiz)
			prequizAdminGroup.PUT("/:id", prequizHandler.UpdatePrequiz)
		}

		// User accessible routes
		prequizGroupHandler.GET("/:id", prequizHandler.GetPrequizByID)
		prequizGroupHandler.GET("/all", prequizHandler.GetAllPrequizzes)
		prequizGroupHandler.GET("/module/:module_id", prequizHandler.GetPrequizzesByModule)
		prequizGroupHandler.GET("/user-answers", prequizHandler.GetUserPrequizAnswers)
		prequizGroupHandler.POST("/submit", prequizHandler.SubmitPrequizAnswer)
	}

	// Video Quiz routes - Admin only untuk CUD, User bisa Read
	videoQuizGroupHandler := router.Group("video-quiz")
	videoQuizGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin and Teacher routes for Create, Update, Delete
		videoQuizAdminGroup := videoQuizGroupHandler.Group("")
		videoQuizAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			videoQuizAdminGroup.POST("/create", videoQuizHandler.CreateVideoQuiz)
			videoQuizAdminGroup.PUT("/:id", videoQuizHandler.UpdateVideoQuiz)
			videoQuizAdminGroup.DELETE("/:id", videoQuizHandler.DeleteVideoQuiz)
		}

		// User accessible routes
		videoQuizGroupHandler.GET("/:id", videoQuizHandler.GetVideoQuizByID)
		videoQuizGroupHandler.GET("/video-material/:video_material_id", videoQuizHandler.GetVideoQuizzesByVideoMaterial)
		videoQuizGroupHandler.POST("/submit", videoQuizHandler.SubmitVideoQuizAnswer)
		videoQuizGroupHandler.GET("/user-answers", videoQuizHandler.GetAllUserVideoQuizAnswers)
		videoQuizGroupHandler.GET("/user-answers/:video_material_id", videoQuizHandler.GetUserVideoQuizAnswers)
	}

	// Admin routes - Khusus untuk management
	adminGroup := router.Group("admin")
	adminGroup.Use(authMiddleware.RequireAuth)
	adminGroup.Use(authMiddleware.RequireAdmin)
	{
		// User management untuk admin
		adminGroup.GET("/users", userHandler.GetAllUsers)
		adminGroup.PUT("/users/:id/role", userHandler.UpdateUserRole)
		adminGroup.DELETE("/users/:id", userHandler.DeleteUser)

		// Role management untuk admin
		adminGroup.GET("/roles", roleHandler.GetAllRoles)
		adminGroup.POST("/roles", roleHandler.CreateRole)
		adminGroup.GET("/roles/:id", roleHandler.GetRoleById)
		adminGroup.PUT("/roles/:id", roleHandler.UpdateRole)
		adminGroup.DELETE("/roles/:id", roleHandler.DeleteRole)
	}

	// Admin/Teacher routes - User account management
	adminTeacherGroup := router.Group("admin")
	adminTeacherGroup.Use(authMiddleware.RequireAuth)
	adminTeacherGroup.Use(authMiddleware.RequireAdminOrTeacher)
	{
		// Update user account (name and email only)
		adminTeacherGroup.PUT("/users/:id/account", userHandler.UpdateAccount)
	}

	// Public RAG endpoint (no auth required for AI/knowledge systems)
	router.GET("/user-activity/for-rag", activityHandler.GetActivitiesForRAG) // Enhanced data for RAG/AI

	// User Activity routes - Admin dan User bisa akses
	activityGroup := router.Group("user-activity")
	activityGroup.Use(authMiddleware.RequireAuth)
	{
		// User accessible routes
		activityGroup.GET("/my-activities", activityHandler.GetMyActivities)
		activityGroup.GET("/my-last", activityHandler.GetMyLastActivity)
		activityGroup.GET("/my-streak", activityHandler.GetMyStreak)
		activityGroup.GET("/recent", activityHandler.GetRecentActivities)
		activityGroup.GET("/most-active", activityHandler.GetMostActiveUsers)                  // Summary version
		activityGroup.GET("/most-active-detailed", activityHandler.GetMostActiveUsersDetailed) // Detailed version

		// Admin and Teacher routes
		activityAdminGroup := activityGroup.Group("")
		activityAdminGroup.Use(authMiddleware.RequireAdminOrTeacher)
		{
			activityAdminGroup.GET("/users/:user_id", activityHandler.GetUserActivities)
			activityAdminGroup.GET("/stats", activityHandler.GetActivityStats)
			activityAdminGroup.POST("/log", activityHandler.LogActivity)
		}
	}

	// Dashboard routes - Admin dan User bisa akses
	dashboardGroup := router.Group("dashboard")
	dashboardGroup.Use(authMiddleware.RequireAuth)
	{
		dashboardGroup.GET("/", dashboardHandler.GetDashboard)
	}

	// File upload routes - User perlu auth
	uploadGroup := router.Group("upload")
	uploadGroup.Use(authMiddleware.RequireAuth)
	{
		// General image upload
		uploadGroup.POST("/image", func(c *gin.Context) {
			fileUploadHandler := handlers.NewFileUploadHandler()
			fileUploadHandler.UploadImage(c)
		})

		// Multiple images upload
		uploadGroup.POST("/images", func(c *gin.Context) {
			fileUploadHandler := handlers.NewFileUploadHandler()
			fileUploadHandler.UploadMultipleImages(c)
		})
	}

	// Leaderboard routes - User perlu auth
	leaderboardGroup := router.Group("leaderboard")
	leaderboardGroup.Use(authMiddleware.RequireAuth)
	{
		leaderboardGroup.GET("/", leaderboardHandler.GetLeaderboard)
		leaderboardGroup.GET("/user/:user_id", leaderboardHandler.GetUserRank)
	}

	// Flashcard routes - User perlu auth untuk FSRS algorithm
	flashcardGroup := router.Group("flashcard")
	flashcardGroup.Use(authMiddleware.RequireAuth)
	{
		flashcardGroup.GET("/all", flashcardHandler.GetAllFlashcards)                                     // Get all flashcards
		flashcardGroup.GET("/intervals", flashcardHandler.GetFlashcardIntervals)                          // Get interval options (1m, 5m, 7h, 10h)
		flashcardGroup.POST("/module/:module_id/initialize", flashcardHandler.InitializeModuleFlashcards) // Copy/initialize all flashcards in module
		flashcardGroup.POST("/:flashcard_id/review", flashcardHandler.ReviewFlashcard)                    // Review flashcard with grade
		flashcardGroup.POST("/:flashcard_id/initialize", flashcardHandler.InitializeFlashcard)            // Initialize single flashcard
		flashcardGroup.GET("/due", flashcardHandler.GetDueFlashcards)                                     // Get due flashcards
		flashcardGroup.GET("/stats", flashcardHandler.GetRetentionStats)                                  // Get retention statistics
	}

	return router
}
