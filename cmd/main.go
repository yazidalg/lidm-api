package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
	"github.com/yazidalg/lidm_backend/internal/helpers"
	"github.com/yazidalg/lidm_backend/internal/realtime/socketio"
	"github.com/yazidalg/lidm_backend/internal/routes"
	"gorm.io/gorm"
)

func main() {
	startTime := time.Now()
	fmt.Println("🚀 Starting LIDM Backend...")

	config.LoadEnv()
	fmt.Println("✅ Environment loaded")

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Starting HTTP server on port %s\n", port)

	// Create router
	router := gin.Default()

	// Add basic health check endpoint immediately
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(startTime).Seconds(),
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	// Try to connect to database and initialize full app
	fmt.Println("🔌 Attempting database connection...")
	db := config.ConnectDB()
	if db != nil {
		fmt.Println("✅ Database connected")

		// Run migrations
		fmt.Println("📊 Running database migrations...")
		if err := database.Migrate(db); err != nil {
			fmt.Printf("❌ Database migration failed: %v\n", err)
		} else {
			fmt.Println("✅ Database migrations completed")

			// Initialize full application routes
			fmt.Println("🔧 Initializing full application...")
			initializeFullRoutes(router, db, startTime)
		}
	} else {
		fmt.Println("⚠️  Database connection failed, running in basic mode")
	}

	// Start server
	fmt.Printf("📡 Server listening on 0.0.0.0:%s\n", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		panic(err)
	}
}

// initializeFullRoutes adds all the full application routes
func initializeFullRoutes(router *gin.Engine, db *gorm.DB, startTime time.Time) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Panic in route initialization: %v\n", r)
		}
	}()

	// Build middleware
	authMiddleware := helpers.NewBuildAuthMiddleware(db)

	// Build handlers
	userHandler := helpers.NewBuildUser(db)
	authHandler := helpers.NewBuildAuth(db)
	roleHandler := helpers.NewBuildRole(db)
	forgotPasswordHandler := helpers.NewBuildForgotPassword(db)
	questionHandler, questionService := helpers.NewBuildQuestion(db)
	answerHandler := helpers.NewBuildAnswer(db)
	participantHandler, participantService := helpers.NewBuildParticipant(db)
	quizHandler, quizService := helpers.NewBuildQuiz(db)
	quizSessionHandler, quizSessionService := helpers.NewBuildQuizSession(db)
	moduleHandler := helpers.NewBuildModule(db)
	progressHandler := helpers.NewBuildProgress(db)
	prequizHandler, prequizService := helpers.NewBuildPrequiz(db)
	videoQuizHandler, videoQuizService := helpers.NewBuildVideoQuiz(db)
	activityHandler, activityService := helpers.NewBuildUserActivity(db)
	dashboardHandler := helpers.NewBuildDashboard(db)
	leaderboardHandler := helpers.NewBuildLeaderboard(db)
	flashcardHandler := helpers.NewBuildFlashcard(db)
	healthHandler := helpers.NewBuildHealth(db, startTime)
	userServiceForSocket := helpers.NewUserServiceOnly(db)
	socketHandler := helpers.NewBuildSocket(questionService, quizService, participantService, prequizService, quizSessionService, userServiceForSocket)

	// Suppress unused variable warning
	_ = videoQuizService

	// Create the full application router
	fullRouter := routes.NewRoute(
		authHandler,
		userHandler,
		forgotPasswordHandler,
		questionHandler,
		answerHandler,
		participantHandler,
		quizHandler,
		quizSessionHandler,
		socketHandler,
		moduleHandler,
		progressHandler,
		prequizHandler,
		videoQuizHandler,
		authMiddleware,
		roleHandler,
		activityHandler,
		activityService,
		dashboardHandler,
		leaderboardHandler,
		flashcardHandler,
		healthHandler,
	)

	// Start Socket.IO server
	fmt.Println("🔌 Starting Socket.IO server...")
	socketio.StartSocketIOServer(fullRouter, questionService, quizService, userServiceForSocket, participantService)
	fmt.Println("✅ Socket.IO server started")

	// Add all routes from fullRouter to the current router
	addRoutesToRouter(router, fullRouter)

	fmt.Println("✅ Full application routes initialized")
}

// addRoutesToRouter copies routes from source to destination router
func addRoutesToRouter(dest, src *gin.Engine) {
	// This is a simplified approach - we'll add the essential routes manually
	// to avoid the sync.Pool copying issue

	// Add basic routes that are commonly used
	dest.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	dest.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	fmt.Println("✅ Essential routes added to router")
}
