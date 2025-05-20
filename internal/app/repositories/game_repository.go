package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type GameRepositoryInterface interface {
	CreateGame(game *models.Game) (*models.Game, error)
	GetGameById(id int) (*models.Game, error)
	GetAllGames() (*[]models.Game, error)
	UpdateGame(id int, game *models.Game) (*models.Game, error)
	DeleteGame(id int) (*models.Game, error)
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) *gameRepository {
	return &gameRepository{db}
}

func (r *gameRepository) CreateGame(game *models.Game) (*models.Game, error) {
	err := r.db.Create(game).Error
	return game, err
}

func (r *gameRepository) GetGameById(id int) (*models.Game, error) {
	var game models.Game
	err := r.db.First(&game, "id = ?", id).Error
	return &game, err
}

func (r *gameRepository) GetAllGames() (*[]models.Game, error) {
	var games []models.Game
	err := r.db.Find(&games).Error
	return &games, err
}

func (r *gameRepository) UpdateGame(id int, game *models.Game) (*models.Game, error) {
	err := r.db.Model(&game).Where("id = ?", id).Updates(game).Error
	return game, err
}

func (r *gameRepository) DeleteGame(id int) (*models.Game, error) {
	var game models.Game
	err := r.db.Where("id = ?", id).Delete(&game).Error
	return &game, err
}
