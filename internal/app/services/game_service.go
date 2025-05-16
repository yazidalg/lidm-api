package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type GameServiceInterface interface {
	CreateGame(game *models.Game) (*models.Game, error)
	GetGameById(id int) (*models.Game, error)
	GetAllGames() (*[]models.Game, error)
	UpdateGame(id int, game *models.Game) (*models.Game, error)
	DeleteGame(id int) (*models.Game, error)
}

type gameService struct {
	gameRepository repositories.GameRepositoryInterface
}

func NewGameService(gameRepository repositories.GameRepositoryInterface) *gameService {
	return &gameService{gameRepository}
}

func (s *gameService) CreateGame(game *models.Game) (*models.Game, error) {
	return s.gameRepository.CreateGame(game)
}

func (s *gameService) GetGameById(id int) (*models.Game, error) {
	return s.gameRepository.GetGameById(id)
}

func (s *gameService) GetAllGames() (*[]models.Game, error) {
	return s.gameRepository.GetAllGames()
}

func (s *gameService) UpdateGame(id int, game *models.Game) (*models.Game, error) {
	return s.gameRepository.UpdateGame(id, game)
}

func (s *gameService) DeleteGame(id int) (*models.Game, error) {
	return s.gameRepository.DeleteGame(id)
}
