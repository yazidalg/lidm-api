package main

import (
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
	"github.com/yazidalg/lidm_backend/internal/helpers"
	"github.com/yazidalg/lidm_backend/internal/routes"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	database.Migrate(db)

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
	lessonHandler := helpers.NewBuildLesson(db)
	progressHandler := helpers.NewBuildProgress(db)
	prequizHandler, prequizService := helpers.NewBuildPrequiz(db)
	videoQuizHandler, videoQuizService := helpers.NewBuildVideoQuiz(db)
	activityHandler, activityService := helpers.NewBuildUserActivity(db)
	dashboardHandler := helpers.NewBuildDashboard(db)
	socketHandler := helpers.NewBuildSocket(questionService, quizService, participantService, prequizService, quizSessionService)

	// Suppress unused variable warning for videoQuizService if needed
	_ = videoQuizService

	routes.NewRoute(
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
		lessonHandler,
		progressHandler,
		prequizHandler,
		videoQuizHandler,
		authMiddleware,
		roleHandler, // Add roleHandler parameter
		activityHandler,
		activityService,
		dashboardHandler,
	)
}
