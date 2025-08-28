package main

import (
	"fmt"
	"log"
	"time"

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

	// Manually update to 100% without triggering auto-unlock
	now := time.Now()
	progress.Progress = 100.0
	progress.IsComplete = true
	progress.CompletedAt = &now

	// Direct database update to avoid model triggers
	err = db.Model(progress).Updates(map[string]interface{}{
		"progress":     100.0,
		"is_complete":  true,
		"completed_at": now,
	}).Error

	if err != nil {
		log.Fatal("Failed to update progress:", err)
	}

	fmt.Printf("SUCCESS! Updated progress to: 100.00%%\n")
	fmt.Printf("Is Complete: true\n")
}
