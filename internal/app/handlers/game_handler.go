package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type GameHandler struct {
	gameService services.GameServiceInterface
}

func NewGameHandler(gameService services.GameServiceInterface) *GameHandler {
	return &GameHandler{gameService}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
}
