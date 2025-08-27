package main

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
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

	// Get all module progress records
	var allProgress []models.ModuleProgress
	if err := db.Find(&allProgress).Error; err != nil {
		log.Fatal("Failed to get module progress records:", err)
	}

	fmt.Printf("Found %d module progress records to recalculate\n", len(allProgress))

	// Recalculate progress for each record
	successCount := 0
	for _, progress := range allProgress {
		fmt.Printf("Recalculating progress for User ID: %d, Module ID: %d... ", progress.UserID, progress.ModuleID)
		
		// Update progress using the new logic
		_, err := moduleProgressService.UpdateUserProgress(progress.UserID, progress.ModuleID)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			continue
		}
		
		// Get updated progress to show new value
		updatedProgress, err := moduleProgressService.GetUserModuleProgress(progress.UserID, progress.ModuleID)
		if err != nil {
			fmt.Printf("FAILED to get updated progress: %v\n", err)
			continue
		}
		
		fmt.Printf("SUCCESS: %.2f%%\n", updatedProgress.Progress)
		successCount++
	}

	fmt.Printf("\nRecalculation complete! Successfully updated %d out of %d records.\n", successCount, len(allProgress))
}
