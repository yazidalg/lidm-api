package helpers

import (
	"github.com/yazidalg/lidm_backend/internal/app/handlers"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"gorm.io/gorm"
)

func NewBuildAuth(db *gorm.DB) *handlers.AuthHandler {
	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	return authHandler
}
