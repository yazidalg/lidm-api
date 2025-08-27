package main

import (
	"fmt"
	"log"

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

	// Recalculate progress for Module 2, User 2
	fmt.Printf("Recalculating progress for User ID: 2, Module ID: 2...\n")
	
	// Update progress using the new logic
	updatedProgress, err := moduleProgressService.UpdateUserProgress(2, 2)
	if err != nil {
		log.Fatal("Failed to update progress:", err)
	}
	
	fmt.Printf("SUCCESS! New progress: %.2f%%\n", updatedProgress.Progress)
	fmt.Printf("Is Complete: %t\n", updatedProgress.IsComplete)
	fmt.Printf("Is Unlocked: %t\n", updatedProgress.IsUnlocked)
}
