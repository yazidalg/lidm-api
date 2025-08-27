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

	fmt.Printf("Checking Module 2 progress for User 2...\n")
	
	// Get current progress
	progress, err := moduleProgressRepo.GetByUserAndModule(2, 2)
	if err != nil {
		log.Fatal("Failed to get progress:", err)
	}
	
	fmt.Printf("✅ Module 2 Progress: %.2f%%\n", progress.Progress)
	fmt.Printf("✅ Is Complete: %t\n", progress.IsComplete)
	fmt.Printf("✅ Is Unlocked: %t\n", progress.IsUnlocked)
	
	if progress.CompletedAt != nil {
		fmt.Printf("✅ Completed At: %s\n", progress.CompletedAt.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("❌ Completed At: nil\n")
	}
}
