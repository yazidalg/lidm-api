package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type GameHandler struct {
	gameService services.GameServiceInterface
}

func NewGameHandler(gameService services.GameServiceInterface) *GameHandler {
	return &GameHandler{gameService}
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	var gameReq request.GameRequest
	if err := c.ShouldBindJSON(&gameReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Convert string IDs to uint
	userID, err := strconv.ParseUint(gameReq.UserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID format",
			"error":   err.Error(),
		})
		return
	}

	quizID, err := strconv.ParseUint(gameReq.QuizID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid quiz ID format",
			"error":   err.Error(),
		})
		return
	}

	// Create the game model with proper ID types
	game := models.Game{
		UserID: uint(userID),
		QuizID: uint(quizID),
		Score:  gameReq.Score,
		Answer: gameReq.Answer,
	}

	createdGame, err := h.gameService.CreateGame(&game)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to create game"

		// Check for specific error types
		switch err.Error() {
		case "user not found":
			status = http.StatusBadRequest
			message = "User not found - please use a valid user ID"
		case "quiz not found":
			status = http.StatusBadRequest
			message = "Quiz not found - please use a valid quiz ID"
		}

		c.JSON(status, gin.H{
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Game created successfully",
		"data":    createdGame,
	})
}

func (h *GameHandler) GetGameById(c *gin.Context) {
	id := c.Param("id")
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid game ID",
			"error":   err.Error(),
		})
		return
	}

	game, err := h.gameService.GetGameById(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Game not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game retrieved successfully",
		"data":    game,
	})
}

func (h *GameHandler) GetAllGames(c *gin.Context) {
	games, err := h.gameService.GetAllGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve games",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Games retrieved successfully",
		"data":    games,
	})
}

func (h *GameHandler) UpdateGame(c *gin.Context) {
	id := c.Param("id")
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid game ID",
			"error":   err.Error(),
		})
		return
	}

	var game models.Game
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	updatedGame, err := h.gameService.UpdateGame(gameID, &game)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update game",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game updated successfully",
		"data":    updatedGame,
	})
}

func (h *GameHandler) DeleteGame(c *gin.Context) {
	id := c.Param("id")
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid game ID",
			"error":   err.Error(),
		})
		return
	}

	deletedGame, err := h.gameService.DeleteGame(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete game",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Game deleted successfully",
		"data":    deletedGame,
	})
}
