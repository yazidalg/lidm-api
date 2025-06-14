package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/middleware"
)

func NewRoute(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	forgotPasswordHandler *handlers.ForgotPasswordHandler,
	questionHandler *handlers.QuestionHandler,
	answerHandler *handlers.AnswerHandler,
	participantHandler *handlers.ParticipantHandler,
	quizHandler *handlers.QuizHandler,
	socketHandler *handlers.SocketHandler,
) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	userGroupHandler := router.Group("user")
	userGroupHandler.Use(middleware.AuthRequire)
	userGroupHandler.GET("/profile", userHandler.GetUserById)

	authGroupHandler := router.Group("auth")
	authGroupHandler.POST("/register", authHandler.RegisterUser)
	authGroupHandler.POST("/login", authHandler.LoginUser)
	authGroupHandler.GET("/verify/:verificationToken", authHandler.VerifyEmail)

	// Forgot password routes
	forgotPasswordGroup := router.Group("password")
	forgotPasswordGroup.POST("/forgot", forgotPasswordHandler.RequestPasswordReset)
	forgotPasswordGroup.POST("/verify-otp", forgotPasswordHandler.VerifyOTP)
	forgotPasswordGroup.POST("/reset", forgotPasswordHandler.ResetPassword)

	questionGroupHandler := router.Group("question")
	questionGroupHandler.Use(middleware.AuthRequire)
	questionGroupHandler.POST("/create", questionHandler.CreateQuestion)
	questionGroupHandler.GET("/:id", questionHandler.GetQuestionByID)
	questionGroupHandler.GET("/all", questionHandler.GetAllQuestions)
	questionGroupHandler.PUT("/:id", questionHandler.UpdateQuestion)
	questionGroupHandler.DELETE("/:id", questionHandler.DeleteQuestion)

	answerGroupHandler := router.Group("answer")
	answerGroupHandler.Use(middleware.AuthRequire)
	answerGroupHandler.POST("/create", answerHandler.CreateAnswer)
	answerGroupHandler.GET("/:id", answerHandler.GetAnswerByID)
	answerGroupHandler.GET("/all", answerHandler.GetAllAnswers)
	answerGroupHandler.PUT("/:id", answerHandler.UpdateAnswer)
	answerGroupHandler.DELETE("/:id", answerHandler.DeleteAnswer)

	participantGroupHandler := router.Group("participant")
	participantGroupHandler.Use(middleware.AuthRequire)
	participantGroupHandler.POST("/create", participantHandler.CreateParticipant)
	participantGroupHandler.GET("/:id", participantHandler.GetParticipantByID)
	participantGroupHandler.GET("/all", participantHandler.GetAllParticipants)
	participantGroupHandler.GET("/quiz/:quiz_id", participantHandler.GetParticipantsByQuizID)
	participantGroupHandler.GET("/user/:user_id", participantHandler.GetParticipantsByUserID)
	participantGroupHandler.PUT("/:id", participantHandler.UpdateParticipant)
	participantGroupHandler.DELETE("/:id", participantHandler.DeleteParticipant)

	quizGroupHandler := router.Group("quiz")
	quizGroupHandler.Use(middleware.AuthRequire)
	quizGroupHandler.POST("/create", quizHandler.CreateQuiz)
	quizGroupHandler.GET("/:id", quizHandler.GetQuizByID)
	quizGroupHandler.GET("/all", quizHandler.GetAllQuizzes)
	quizGroupHandler.PUT("/:id", quizHandler.UpdateQuiz)
	quizGroupHandler.DELETE("/:id", quizHandler.DeleteQuiz)

	socketGroupHandler := router.Group("ws")
	socketGroupHandler.Use(middleware.AuthRequire)
	socketGroupHandler.GET("/:roomName", socketHandler.ServeWs)
	socketGroupHandler.GET("/matchmaking", socketHandler.MatchMaking)

	router.Run(":3000") // Use PORT from environment variable, default to 8080 if not set
}
