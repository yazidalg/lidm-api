package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
)

func NewRoute(
	authHandler *handlers.AuthHandler,
) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API",
		})
	})

	authGroupHandler := router.Group("auth")
	authGroupHandler.POST("/register", authHandler.RegisterUser)
	authGroupHandler.POST("/login", authHandler.LoginUser)
	authGroupHandler.GET("/verify/:verificationToken", authHandler.VerifyEmail)

	router.Run(":3000")
}
