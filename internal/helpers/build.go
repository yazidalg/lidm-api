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

func NewBuildQuiz(db *gorm.DB) *handlers.QuizHandler {
	quizRepository := repositories.NewQuizRepository(db)
	quizService := services.NewQuizService(quizRepository)
	quizHandler := handlers.NewQuizHandler(quizService)
	return quizHandler
}
