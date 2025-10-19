package main

import (
	"fmt" // Import the fmt package
	"os"
	"strings"

	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
	"github.com/yazidalg/lidm_backend/internal/helpers"
	"github.com/yazidalg/lidm_backend/internal/realtime/socketio"
	"github.com/yazidalg/lidm_backend/internal/routes"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	database.Migrate(db)

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
	)

	// Start Socket.IO server (mounts /socket.io endpoints)
	socketio.StartSocketIOServer(router, questionService, quizService, userServiceForSocket, participantService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("🚀 Starting server on port %s\n", port)
	err := router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}