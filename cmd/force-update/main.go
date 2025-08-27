package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Initialize repositories
	moduleProgressRepo := repositories.NewModuleProgressRepository(db)
	userRepo := repositories.NewUserRepository(db)
	moduleRepo := repositories.NewModuleRepository(db)
	prequizRepo := repositories.NewPrequizRepository(db)
	videoQuizRepo := repositories.NewVideoQuizRepository(db)

	// Initialize service
	moduleProgressService := services.NewModuleProgressService(
		moduleProgressRepo,
		userRepo,
		moduleRepo,
		prequizRepo,
		videoQuizRepo,
	)

	userID := uint(2)
	moduleID := uint(1)

	fmt.Printf("Force recalculating progress for User ID: %d, Module ID: %d...\n", userID, moduleID)

	// Force update progress
	updatedProgress, err := moduleProgressService.UpdateUserProgress(userID, moduleID)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS! Updated progress:\n")
	fmt.Printf("Progress: %.2f%%\n", updatedProgress.Progress)
	fmt.Printf("IsComplete: %t\n", updatedProgress.IsComplete)
	fmt.Printf("IsUnlocked: %t\n", updatedProgress.IsUnlocked)
}
