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

	fmt.Printf("📊 Final Status Check for User 2:\n")
	fmt.Printf("=====================================\n")

	// Check Module 1
	if progress1, err := moduleProgressRepo.GetByUserAndModule(2, 1); err == nil {
		fmt.Printf("Module 1: %.1f%% | Complete: %t | Unlocked: %t\n",
			progress1.Progress, progress1.IsComplete, progress1.IsUnlocked)
	}

	// Check Module 2
	if progress2, err := moduleProgressRepo.GetByUserAndModule(2, 2); err == nil {
		fmt.Printf("Module 2: %.1f%% | Complete: %t | Unlocked: %t\n",
			progress2.Progress, progress2.IsComplete, progress2.IsUnlocked)
	}

	// Check Module 3
	if progress3, err := moduleProgressRepo.GetByUserAndModule(2, 3); err == nil {
		fmt.Printf("Module 3: %.1f%% | Complete: %t | Unlocked: %t\n",
			progress3.Progress, progress3.IsComplete, progress3.IsUnlocked)
	} else {
		fmt.Printf("Module 3: No progress record\n")
	}

	fmt.Printf("=====================================\n")
	fmt.Printf("🎯 Module 2 should now show 100%% and allow progression to Module 3!\n")
}
