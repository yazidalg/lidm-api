package main

import (
	"fmt"

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

	fmt.Printf("📊 Checking All Module Progress for User 2:\n")
	fmt.Printf("==========================================\n")
	
	for moduleID := 1; moduleID <= 6; moduleID++ {
		progress, err := moduleProgressRepo.GetByUserAndModule(2, uint(moduleID))
		if err != nil {
			fmt.Printf("Module %d: No progress record\n", moduleID)
		} else {
			fmt.Printf("Module %d: %.1f%% | Complete: %t | Unlocked: %t\n", 
				moduleID, progress.Progress, progress.IsComplete, progress.IsUnlocked)
		}
	}
	
	fmt.Printf("==========================================\n")
}
