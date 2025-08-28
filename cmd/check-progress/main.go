package main

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Query the module progress directly from database
	var progress models.ModuleProgress
	result := db.Where("user_id = ? AND module_id = ?", 2, 1).First(&progress)

	if result.Error != nil {
		log.Fatal("Failed to get module progress:", result.Error)
	}

	fmt.Printf("Current database values for User ID: 2, Module ID: 1:\n")
	fmt.Printf("Progress: %.2f%%\n", progress.Progress)
	fmt.Printf("IsComplete: %t\n", progress.IsComplete)
	fmt.Printf("IsUnlocked: %t\n", progress.IsUnlocked)
	fmt.Printf("StartedAt: %v\n", progress.StartedAt)
	fmt.Printf("CompletedAt: %v\n", progress.CompletedAt)
}
