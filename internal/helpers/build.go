package helpers

import (
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/app/socket"
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

func NewBuildQuestion(db *gorm.DB) (*handlers.QuestionHandler, services.QuestionServiceInterface) {
	questionRepository := repositories.NewQuestionRepository(db)
	questionService := services.NewQuestionService(questionRepository)
	questionHandler := handlers.NewQuestionHandler(questionService)
	return questionHandler, questionService
}

func NewBuildAnswer(db *gorm.DB) *handlers.AnswerHandler {
	answerRepository := repositories.NewAnswerRepository(db)
	answerService := services.NewAnswerService(answerRepository)
	answerHandler := handlers.NewAnswerHandler(answerService)
	return answerHandler
}

func NewBuildParticipant(db *gorm.DB) (*handlers.ParticipantHandler, services.ParticipantServiceInterface) {
	participantRepository := repositories.NewParticipantRepository(db)
	participantService := services.NewParticipantService(participantRepository)
	participantHandler := handlers.NewParticipantHandler(participantService)
	return participantHandler, participantService
}

func NewBuildQuiz(db *gorm.DB) (*handlers.QuizHandler, services.QuizServiceInterface) {
	quizRepository := repositories.NewQuizRepository(db)
	quizService := services.NewQuizService(quizRepository)
	quizHandler := handlers.NewQuizHandler(quizService)
	return quizHandler, quizService
}

func NewBuildSocket(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
	prequizService services.PrequizServiceInterface,
) *handlers.SocketHandler {
	hub := socket.NewHub(questionService, quizService, participantService, prequizService)
	go hub.Run() // Jalankan Hub sebagai goroutine

	// Daftarkan handler WebSocket
	socketHandler := handlers.NewSocketHandler(hub)

	return socketHandler
}

func NewBuildModule(db *gorm.DB) *handlers.ModuleHandler {
	moduleRepository := repositories.NewModuleRepository(db)
	moduleService := services.NewModuleService(moduleRepository)
	moduleHandler := handlers.NewModuleHandler(moduleService)
	return moduleHandler
}

func NewBuildLesson(db *gorm.DB) *handlers.LessonHandler {
	moduleRepository := repositories.NewModuleRepository(db)
	moduleService := services.NewModuleService(moduleRepository)
	lessonRepository := repositories.NewLessonRepository(db)
	lessonService := services.NewLessonService(lessonRepository, moduleRepository)
	lessonHandler := handlers.NewLessonHandler(lessonService, moduleService)
	return lessonHandler
}

func NewBuildProgress(db *gorm.DB) *handlers.ProgressHandler {
	moduleRepository := repositories.NewModuleRepository(db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	lessonRepository := repositories.NewLessonRepository(db)
	lessonService := services.NewLessonService(lessonRepository, moduleRepository)

	progressRepository := repositories.NewProgressRepository(db)
	progressService := services.NewProgressService(progressRepository, userRepository, lessonRepository)
	progressHandler := handlers.NewProgressHandler(progressService, userService, lessonService)

	return progressHandler
}

func NewBuildPrequiz(db *gorm.DB) (*handlers.PrequizHandler, services.PrequizServiceInterface) {
	moduleRepository := repositories.NewModuleRepository(db)

	prequizRepository := repositories.NewPrequizRepository(db)

	lessonRepository := repositories.NewLessonRepository(db)
	lessonService := services.NewLessonService(lessonRepository, moduleRepository)

	userRepository := repositories.NewUserRepository(db)

	prequizService := services.NewPrequizService(prequizRepository, lessonRepository, userRepository)
	prequizHandler := handlers.NewPrequizHandler(prequizService, lessonService, userRepository)

	return prequizHandler, prequizService
}
