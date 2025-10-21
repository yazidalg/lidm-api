package main

import (
	"context"
	"fmt" // Import the fmt package
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
	"github.com/yazidalg/lidm_backend/internal/helpers"
	"github.com/yazidalg/lidm_backend/internal/realtime/socketio"
	"github.com/yazidalg/lidm_backend/internal/routes"
)

func main() {
	startTime := time.Now()
	fmt.Println("🚀 Starting LIDM Backend...")

	config.LoadEnv()
	fmt.Println("✅ Environment loaded")

	// Start HTTP server first to satisfy Cloud Run health checks
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Starting HTTP server on port %s\n", port)
	fmt.Printf("📡 Server will listen on 0.0.0.0:%s\n", port)

	// Create a basic router first for health checks
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

	// Start server in a goroutine to allow database connection
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Failed to start initial server: %v\n", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	fmt.Println("🔌 Connecting to database...")
	db := config.ConnectDB()
	if db == nil {
		fmt.Println("⚠️  Database connection failed, starting in limited mode...")
		// Start a basic server without database features
		basicRouter := gin.Default()
		basicRouter.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "ok",
				"timestamp": time.Now().Unix(),
				"uptime":    time.Since(startTime).Seconds(),
				"mode":      "limited",
			})
		})
		basicRouter.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Welcome to the API (Limited Mode)"})
		})

		// Gracefully shutdown the initial server
		server.Shutdown(context.TODO())

		// Start the basic server
		basicServer := &http.Server{
			Addr:    "0.0.0.0:" + port,
			Handler: basicRouter,
		}

		fmt.Printf("✅ Basic server started on port %s (limited mode)\n", port)
		if err := basicServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Failed to start basic server: %v\n", err)
			panic(fmt.Sprintf("Failed to start basic server: %v", err))
		}
		return
	}
	fmt.Println("✅ Database connected")

	fmt.Println("📊 Running database migrations...")
	database.Migrate(db)
	fmt.Println("✅ Database migrations completed")

	// Optional quiz seeding via env: SEED_QUIZ_MODULES="Module Title 1,Module Title 2"
	if modsEnv := os.Getenv("SEED_QUIZ_MODULES"); modsEnv != "" {
		modules := strings.Split(modsEnv, ",")
		for i := range modules {
			modules[i] = strings.TrimSpace(modules[i])
		}
		database.SeedQuizData(db, modules)
	}

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
	// Health handler for monitoring
	healthHandler := helpers.NewBuildHealth(db, startTime)
	// User service khusus untuk socket (lives & xp)
	userServiceForSocket := helpers.NewUserServiceOnly(db)
	socketHandler := helpers.NewBuildSocket(questionService, quizService, participantService, prequizService, quizSessionService, userServiceForSocket)

	// Suppress unused variable warning for videoQuizService if needed
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

	// Start Socket.IO server (mounts /socket.io endpoints)
	fmt.Println("🔌 Starting Socket.IO server...")
	socketio.StartSocketIOServer(fullRouter, questionService, quizService, userServiceForSocket, participantService)
	fmt.Println("✅ Socket.IO server started")

	// Gracefully shutdown the initial server and start the full application
	fmt.Println("🔄 Restarting server with full application...")
	server.Shutdown(context.TODO())

	// Start the full application server
	fullServer := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: fullRouter,
	}

	fmt.Printf("✅ Full application server started on port %s\n", port)

	// Keep the application running
	if err := fullServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("❌ Failed to start full server: %v\n", err)
		panic(fmt.Sprintf("Failed to start full server: %v", err))
	}
}
