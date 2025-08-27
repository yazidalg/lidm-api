package main

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Initialize repositories
	moduleProgressRepo := repositories.NewModuleProgressRepository(db)

	fmt.Printf("Manually updating Module 2 progress for User 2...\n")
	
	// Get existing progress
	progress, err := moduleProgressRepo.GetByUserAndModule(2, 2)
	if err != nil {
		log.Fatal("Failed to get progress:", err)
	}
	
	fmt.Printf("Current progress: %.2f%%\n", progress.Progress)
	
	// Manually update to 100%
	progress.Progress = 100.0
	progress.IsComplete = true
	progress.MarkAsCompleted()
	
	// Update in database
	_, err = moduleProgressRepo.UpdateProgress(progress)
	if err != nil {
		log.Fatal("Failed to update progress:", err)
	}
	
	fmt.Printf("SUCCESS! Updated progress to: %.2f%%\n", progress.Progress)
	fmt.Printf("Is Complete: %t\n", progress.IsComplete)
}
