package main

import (
	"fmt" // Import the fmt package
	"os"
	"strings"
	"time"

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

	fmt.Println("🔌 Connecting to database...")
	db := config.ConnectDB()
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

	router := routes.NewRoute(
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
		roleHandler, // Add roleHandler parameter
		activityHandler,
		activityService,
		dashboardHandler,
		leaderboardHandler,
		flashcardHandler,
		healthHandler,
	)

	// Start Socket.IO server (mounts /socket.io endpoints)
	fmt.Println("🔌 Starting Socket.IO server...")
	socketio.StartSocketIOServer(router, questionService, quizService, userServiceForSocket, participantService)
	fmt.Println("✅ Socket.IO server started")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Starting HTTP server on port %s\n", port)
	fmt.Printf("📡 Server will listen on 0.0.0.0:%s\n", port)

	err := router.Run("0.0.0.0:" + port)
	if err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
