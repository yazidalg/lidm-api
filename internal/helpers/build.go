package helpers

import (
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
	"github.com/yazidalg/lidm_backend/internal/app/socket/prequiz"
	"github.com/yazidalg/lidm_backend/internal/app/socket/quiz"
	"github.com/yazidalg/lidm_backend/internal/middleware"
	"gorm.io/gorm"
)

func NewBuildUser(db *gorm.DB) *handlers.UserHandler {
	userRepository := repositories.NewUserRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	userService := services.NewUserService(userRepository, roleService)
	userHandler := handlers.NewUserHandler(userService)
	return userHandler
}

// NewUserServiceOnly - untuk kebutuhan internal (misal socket) tanpa handler
func NewUserServiceOnly(db *gorm.DB) services.UserServiceInterface {
	userRepository := repositories.NewUserRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	return services.NewUserService(userRepository, roleService)
}

func NewBuildAuth(db *gorm.DB) *handlers.AuthHandler {
	authRepository := repositories.NewAuthRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	authService := services.NewAuthService(authRepository, roleService)
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
	participantRepository := repositories.NewParticipantRepository(db)
	participantService := services.NewParticipantService(participantRepository)
	quizHandler := handlers.NewQuizHandler(participantService, quizService)
	return quizHandler, quizService
}

func NewBuildQuizSession(db *gorm.DB) (*handlers.QuizSessionHandler, services.QuizSessionServiceInterface) {
	quizSessionRepository := repositories.NewQuizSessionRepository(db)
	quizRepository := repositories.NewQuizRepository(db)
	questionRepository := repositories.NewQuestionRepository(db)
	participantRepository := repositories.NewParticipantRepository(db)
	moduleRepository := repositories.NewModuleRepository(db)
	userRepository := repositories.NewUserRepository(db)

	quizSessionService := services.NewQuizSessionService(
		quizSessionRepository,
		quizRepository,
		questionRepository,
		participantRepository,
		moduleRepository,
		userRepository,
	)
	quizSessionHandler := handlers.NewQuizSessionHandler(quizSessionService)
	return quizSessionHandler, quizSessionService
}

func NewBuildSocket(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
	prequizService services.PrequizServiceInterface,
	quizSessionService services.QuizSessionServiceInterface,
	// Tambahan: user service untuk lives & xp
	userService services.UserServiceInterface,
) *handlers.SocketHandler {
	// Factory functions
	quizSessionFactory := func(hub *common.Hub, roomName string, players []*common.Client, questions []models.Question, participants []*models.Participant, quizID uint) common.QuizSessionInterface {
		return quiz.NewQuizSession(hub, roomName, players, questions, participants, quizID)
	}

	prequizSessionFactory := func(hub *common.Hub, roomName string, player *common.Client, questions []models.Prequiz) common.PrequizSessionInterface {
		return prequiz.NewPrequizSession(hub, roomName, player, questions)
	}

	hub := common.NewHub(questionService, quizService, participantService, prequizService, quizSessionFactory, prequizSessionFactory, userService)
	go hub.Run() // Jalankan Hub sebagai goroutine

	// Daftarkan handler WebSocket
	socketHandler := handlers.NewSocketHandler(hub)

	return socketHandler
}

func NewBuildModule(db *gorm.DB) *handlers.ModuleHandler {
	moduleRepository := repositories.NewModuleRepository(db)
	progressRepository := repositories.NewProgressRepository(db)
	videoQuizRepository := repositories.NewVideoQuizRepository(db)
	prequizRepository := repositories.NewPrequizRepository(db)
	
	// Add module progress repository and service
	moduleProgressRepository := repositories.NewModuleProgressRepository(db)
	userRepository := repositories.NewUserRepository(db)
	moduleProgressService := services.NewModuleProgressService(moduleProgressRepository, userRepository, moduleRepository)
	
	moduleService := services.NewModuleService(moduleRepository, progressRepository, videoQuizRepository, prequizRepository, moduleProgressRepository, moduleProgressService)
	moduleHandler := handlers.NewModuleHandler(moduleService)
	return moduleHandler
}

// NewModuleProgressServiceOnly - untuk kebutuhan internal
func NewModuleProgressServiceOnly(db *gorm.DB) services.ModuleProgressServiceInterface {
	moduleProgressRepository := repositories.NewModuleProgressRepository(db)
	userRepository := repositories.NewUserRepository(db)
	moduleRepository := repositories.NewModuleRepository(db)
	return services.NewModuleProgressService(moduleProgressRepository, userRepository, moduleRepository)
}

func NewBuildProgress(db *gorm.DB) *handlers.ProgressHandler {
	moduleRepository := repositories.NewModuleRepository(db)

	userRepository := repositories.NewUserRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	userService := services.NewUserService(userRepository, roleService)

	// Extra repos/services needed for aggregation
	videoQuizRepository := repositories.NewVideoQuizRepository(db)
	videoQuizService := services.NewVideoQuizService(videoQuizRepository, userRepository)

	prequizRepository := repositories.NewPrequizRepository(db)
	prequizService := services.NewPrequizService(prequizRepository, userRepository)

	progressRepository := repositories.NewProgressRepository(db)
	
	// Add module progress dependencies for this context too
	moduleProgressRepository := repositories.NewModuleProgressRepository(db)
	moduleProgressService := services.NewModuleProgressService(moduleProgressRepository, userRepository, moduleRepository)
	
	moduleService := services.NewModuleService(moduleRepository, progressRepository, videoQuizRepository, prequizRepository, moduleProgressRepository, moduleProgressService)

	progressService := services.NewProgressService(progressRepository, userRepository)                                            // Removed lessonRepository
	progressHandler := handlers.NewProgressHandler(progressService, userService, moduleService, videoQuizService, prequizService) // Removed lessonService

	return progressHandler
}

func NewBuildPrequiz(db *gorm.DB) (*handlers.PrequizHandler, services.PrequizServiceInterface) {
	prequizRepository := repositories.NewPrequizRepository(db)

	userRepository := repositories.NewUserRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	userService := services.NewUserService(userRepository, roleService)

	prequizService := services.NewPrequizService(prequizRepository, userRepository)
	
	// Set module progress service for auto-progress updates
	moduleProgressService := NewModuleProgressServiceOnly(db)
	prequizService.SetModuleProgressService(moduleProgressService)
	
	prequizHandler := handlers.NewPrequizHandler(prequizService, userService)

	return prequizHandler, prequizService
}

func NewBuildAuthMiddleware(db *gorm.DB) *middleware.AuthMiddleware {
	authRepository := repositories.NewAuthRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	authService := services.NewAuthService(authRepository, roleService)
	// user service for life reset
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository, roleService)
	authMiddleware := middleware.NewAuthMiddleware(authService, userService)
	return authMiddleware
}

func NewBuildRole(db *gorm.DB) *handlers.RoleHandler {
	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	roleHandler := handlers.NewRoleHandler(roleService)
	return roleHandler
}

func NewBuildUserActivity(db *gorm.DB) (*handlers.UserActivityHandler, services.UserActivityServiceInterface) {
	activityRepository := repositories.NewUserActivityRepository(db)
	userRepository := repositories.NewUserRepository(db)
	activityService := services.NewUserActivityService(activityRepository, userRepository)

	// Add module service for enhanced RAG data
	moduleRepository := repositories.NewModuleRepository(db)
	progressRepository := repositories.NewProgressRepository(db)
	videoQuizRepository := repositories.NewVideoQuizRepository(db)
	prequizRepository := repositories.NewPrequizRepository(db)
	
	// Add module progress dependencies
	moduleProgressRepository := repositories.NewModuleProgressRepository(db)
	moduleProgressService := services.NewModuleProgressService(moduleProgressRepository, userRepository, moduleRepository)
	
	moduleService := services.NewModuleService(moduleRepository, progressRepository, videoQuizRepository, prequizRepository, moduleProgressRepository, moduleProgressService)

	activityHandler := handlers.NewUserActivityHandler(activityService, moduleService) // Removed lessonService
	return activityHandler, activityService
}

func NewBuildDashboard(db *gorm.DB) *handlers.DashboardHandler {
	activityRepository := repositories.NewUserActivityRepository(db)
	userRepository := repositories.NewUserRepository(db)
	activityService := services.NewUserActivityService(activityRepository, userRepository)

	roleRepository := repositories.NewRoleRepository(db)
	roleService := services.NewRoleService(roleRepository)
	userService := services.NewUserService(userRepository, roleService)

	dashboardHandler := handlers.NewDashboardHandler(activityService, userService)
	return dashboardHandler
}

func NewBuildVideoQuiz(db *gorm.DB) (*handlers.VideoQuizHandler, services.VideoQuizServiceInterface) {
	videoQuizRepository := repositories.NewVideoQuizRepository(db)
	userRepository := repositories.NewUserRepository(db)
	videoQuizService := services.NewVideoQuizService(videoQuizRepository, userRepository)
	
	// Set module progress service for auto-progress updates
	moduleProgressService := NewModuleProgressServiceOnly(db)
	videoQuizService.SetModuleProgressService(moduleProgressService)
	
	videoQuizHandler := handlers.NewVideoQuizHandler(videoQuizService)

	return videoQuizHandler, videoQuizService
}
