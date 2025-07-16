package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
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
	socketHandler *handlers.SocketHandler,
	moduleHandler *handlers.ModuleHandler,
	lessonHandler *handlers.LessonHandler,
	progressHandler *handlers.ProgressHandler,
	prequizHandler *handlers.PrequizHandler,
	authMiddleware *middleware.AuthMiddleware,
	roleHandler *handlers.RoleHandler,
) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	// Public routes (tidak perlu auth)
	authGroupHandler := router.Group("auth")
	{
		authGroupHandler.POST("/register", authHandler.RegisterUser)
		authGroupHandler.POST("/login", authHandler.LoginUser)
		authGroupHandler.GET("/verify/:verificationToken", authHandler.VerifyEmail)
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
	{
		userGroupHandler.GET("/profile", userHandler.GetUserById)
	}

	// Question routes - Admin only untuk CUD, User bisa Read
	questionGroupHandler := router.Group("question")
	questionGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin only routes
		questionAdminGroup := questionGroupHandler.Group("")
		questionAdminGroup.Use(authMiddleware.RequireAdmin)
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
		// Admin only routes
		answerAdminGroup := answerGroupHandler.Group("")
		answerAdminGroup.Use(authMiddleware.RequireAdmin)
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
		// Admin only routes
		participantAdminGroup := participantGroupHandler.Group("")
		participantAdminGroup.Use(authMiddleware.RequireAdmin)
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
	{
		// Admin only routes
		quizAdminGroup := quizGroupHandler.Group("")
		quizAdminGroup.Use(authMiddleware.RequireAdmin)
		{
			quizAdminGroup.POST("/create", quizHandler.CreateQuiz)
			quizAdminGroup.PUT("/:id", quizHandler.UpdateQuiz)
			quizAdminGroup.DELETE("/:id", quizHandler.DeleteQuiz)
		}

		// User accessible routes
		quizGroupHandler.GET("/:id", quizHandler.GetQuizByID)
		quizGroupHandler.GET("/all", quizHandler.GetAllQuizzes)
	}

	// Socket routes - Semua authenticated user bisa akses
	socketGroupHandler := router.Group("ws")
	socketGroupHandler.Use(authMiddleware.RequireAuth)
	{
		socketGroupHandler.GET("/:roomName", socketHandler.ServeWs)
		socketGroupHandler.GET("/matchmaking", socketHandler.MatchMaking)
		socketGroupHandler.GET("/prequiz", socketHandler.PreQuiz)
	}

	// Module routes - Admin only untuk CUD, User bisa Read
	moduleGroupHandler := router.Group("module")
	moduleGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin only routes
		moduleAdminGroup := moduleGroupHandler.Group("")
		moduleAdminGroup.Use(authMiddleware.RequireAdmin)
		{
			moduleAdminGroup.POST("/create", moduleHandler.CreateModule)
			moduleAdminGroup.PUT("/:id", moduleHandler.UpdateModule)
			moduleAdminGroup.DELETE("/:id", moduleHandler.DeleteModule)
		}

		// User accessible routes
		moduleGroupHandler.GET("/:id", moduleHandler.GetModuleByID)
		moduleGroupHandler.GET("/all", moduleHandler.GetAllModules)
	}

	// Lesson routes - Admin only untuk CUD, User bisa Read
	lessonGroupHandler := router.Group("lesson")
	lessonGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin only routes
		lessonAdminGroup := lessonGroupHandler.Group("")
		lessonAdminGroup.Use(authMiddleware.RequireAdmin)
		{
			lessonAdminGroup.POST("/create", lessonHandler.CreateLesson)
			lessonAdminGroup.PUT("/:id", lessonHandler.UpdateLesson)
			lessonAdminGroup.DELETE("/:id", lessonHandler.DeleteLesson)
		}

		// User accessible routes
		lessonGroupHandler.GET("/:id", lessonHandler.GetLessonByID)
		lessonGroupHandler.GET("/all", lessonHandler.GetAllLessons)
	}

	// Progress routes - User untuk data sendiri, Admin untuk semua
	progressGroupHandler := router.Group("progress")
	progressGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin only routes
		progressAdminGroup := progressGroupHandler.Group("")
		progressAdminGroup.Use(authMiddleware.RequireAdmin)
		{
			progressAdminGroup.GET("/all", progressHandler.GetAllProgress)
		}

		// User accessible routes
		progressGroupHandler.POST("/create", progressHandler.CreateProgress)
		progressGroupHandler.PUT("/:id", progressHandler.UpdateProgress)
	}

	// Prequiz routes - Admin only untuk CUD, User bisa Read
	prequizGroupHandler := router.Group("prequiz")
	prequizGroupHandler.Use(authMiddleware.RequireAuth)
	{
		// Admin only routes
		prequizAdminGroup := prequizGroupHandler.Group("")
		prequizAdminGroup.Use(authMiddleware.RequireAdmin)
		{
			prequizAdminGroup.POST("/create", prequizHandler.CreatePrequiz)
			prequizAdminGroup.PUT("/:id", prequizHandler.UpdatePrequiz)
		}

		// User accessible routes
		prequizGroupHandler.GET("/:id", prequizHandler.GetPrequizByID)
		prequizGroupHandler.GET("/all", prequizHandler.GetAllPrequizzes)
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

	router.Run(":3000")
}
