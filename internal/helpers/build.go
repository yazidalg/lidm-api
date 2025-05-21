package helpers

import (
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"gorm.io/gorm"
)

func NewBuildUser(db *gorm.DB) *handlers.UserHandler {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	return userHandler
}

func NewBuildAuth(db *gorm.DB) *handlers.AuthHandler {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	return authHandler
}

func NewBuildForgotPassword(db *gorm.DB) *handlers.ForgotPasswordHandler {
	forgotPasswordRepository := repositories.NewForgotPasswordRepository(db)
	authRepository := repositories.NewAuthRepository(db)
	forgotPasswordService := services.NewForgotPasswordService(forgotPasswordRepository, authRepository)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(forgotPasswordService)
	return forgotPasswordHandler
}

func NewBuildQuestion(db *gorm.DB) *handlers.QuestionHandler {
	questionRepository := repositories.NewQuestionRepository(db)
	questionService := services.NewQuestionService(questionRepository)
	questionHandler := handlers.NewQuestionHandler(questionService)
	return questionHandler
}

func NewBuildAnswer(db *gorm.DB) *handlers.AnswerHandler {
	answerRepository := repositories.NewAnswerRepository(db)
	answerService := services.NewAnswerService(answerRepository)
	answerHandler := handlers.NewAnswerHandler(answerService)
	return answerHandler
}

func NewBuildParticipant(db *gorm.DB) *handlers.ParticipantHandler {
	participantRepository := repositories.NewParticipantRepository(db)
	participantService := services.NewParticipantService(participantRepository)
	participantHandler := handlers.NewParticipantHandler(participantService)
	return participantHandler
}

func NewBuildQuiz(db *gorm.DB) *handlers.QuizHandler {
	quizRepository := repositories.NewQuizRepository(db)
	quizService := services.NewQuizService(quizRepository)
	quizHandler := handlers.NewQuizHandler(quizService)
	return quizHandler
}

func NewBuildSocket(db *gorm.DB) *handlers.SocketHandler {
	questionRepository := repositories.NewQuestionRepository(db)
	questionService := services.NewQuestionService(questionRepository)

	quizRepository := repositories.NewQuizRepository(db)
	quizService := services.NewQuizService(quizRepository)

	participantRepository := repositories.NewParticipantRepository(db)
	participantService := services.NewParticipantService(participantRepository)

	answerRepository := repositories.NewAnswerRepository(db)
	answerService := services.NewAnswerService(answerRepository)

	socketHandler := handlers.NewSocketHandler(questionService, quizService, participantService, answerService)
	return socketHandler
}
