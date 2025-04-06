package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/pkg/middleware"
)

func NewRoute(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the API"})
	})

	userGroupHandler := router.Group("user")
	userGroupHandler.Use(middleware.AuthRequire)
	userGroupHandler.GET("/:id", userHandler.GetUserById)

	authGroupHandler := router.Group("auth")
	authGroupHandler.POST("/register", authHandler.RegisterUser)
	authGroupHandler.POST("/login", authHandler.LoginUser)
	authGroupHandler.GET("/verify/:verificationToken", authHandler.VerifyEmail)

	router.Run(":3000")
}
