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

	userHandler := helpers.NewBuildUser(db)
	authHandler := helpers.NewBuildAuth(db)
	forgotPasswordHandler := helpers.NewBuildForgotPassword(db)
	questionHandler, questionService := helpers.NewBuildQuestion(db)
	answerHandler := helpers.NewBuildAnswer(db)
	participantHandler, participantService := helpers.NewBuildParticipant(db)
	quizHandler, quizService := helpers.NewBuildQuiz(db)
	socketHandler := helpers.NewBuildSocket(questionService, quizService, participantService)
	moduleHandler := helpers.NewBuildModule(db)
	lessonHandler := helpers.NewBuildLesson(db)
	progressHandler := helpers.NewBuildProgress(db)

	routes.NewRoute(
		authHandler,
		userHandler,
		forgotPasswordHandler,
		questionHandler,
		answerHandler,
		participantHandler,
		quizHandler,
		socketHandler,
		moduleHandler,
		lessonHandler,
		progressHandler,
	)
}
