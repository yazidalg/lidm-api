package services

import (
	"errors"

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
	userRepository repositories.UserRepositoryInterface
	quizRepository repositories.QuizRepositoryInterface
}

func NewGameService(
	gameRepository repositories.GameRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
	quizRepository repositories.QuizRepositoryInterface,
) *gameService {
	return &gameService{
		gameRepository: gameRepository,
		userRepository: userRepository,
		quizRepository: quizRepository,
	}
}

func (s *gameService) CreateGame(game *models.Game) (*models.Game, error) {
	// Validate that user exists
	_, err := s.userRepository.GetUserById(int(game.UserID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate that quiz exists
	_, err = s.quizRepository.GetQuizByID(int(game.QuizID))
	if err != nil {
		return nil, errors.New("quiz not found")
	}

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
