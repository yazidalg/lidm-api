package database

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Migrate the database schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Module{},
		&models.Quiz{},
		&models.Participant{},
		&models.Question{},
		&models.Answer{},
		&models.Leaderboard{},
		&models.Course{},
		&models.Module{},
		&models.Lesson{},
		&models.Resource{},
		&models.Progress{},
		&models.Prequiz{},
	); err != nil {
		return err
	}

	return nil
}
