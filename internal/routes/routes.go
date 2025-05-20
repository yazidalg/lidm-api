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

	router.Run(":3000")
}
