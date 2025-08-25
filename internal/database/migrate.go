package database

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Migrate Role first since User depends on it
	if err := db.AutoMigrate(&models.Role{}); err != nil {
		return err
	}

	// Seed default roles first
	if err := models.SeedRoles(db); err != nil {
		return err
	}

	// Then migrate other models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Module{},
		&models.VideoMaterial{},
		&models.VideoQuiz{},
		&models.VideoQuizUserAnswer{},
		&models.ARExperiment{},
		&models.QuizQuestion{},
		&models.Flashcard{},
		&models.UserFlashcardProgress{},
		&models.Quiz{},
		&models.Participant{},
		&models.Question{},
		&models.Answer{},
		&models.Leaderboard{},
		&models.Resource{},
		&models.Progress{},
		&models.Prequiz{},
		&models.PrequizUserAnswer{},
		&models.QuizSession{},
		&models.UserActivity{},
		&models.ModuleProgress{}, // Add new ModuleProgress model
	); err != nil {
		return err
	}

	// Seed default admin user
	if err := SeedAdminUser(db); err != nil {
		return err
	}

	return nil
}
