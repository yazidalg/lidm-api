package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/helpers"
)

func main() {
	log.Println("Starting module progress initialization...")

	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Get module progress service
	moduleProgressService := helpers.NewModuleProgressServiceOnly(db)

	// Get all users
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	log.Printf("Found %d users to initialize", len(users))

	// Initialize module progress for each user
	successCount := 0
	for _, user := range users {
		log.Printf("Initializing module progress for user ID %d (%s)", user.ID, user.Name)

		// Check if user already has any module progress
		var existingProgress []models.ModuleProgress
		if err := db.Where("user_id = ?", user.ID).Find(&existingProgress).Error; err != nil {
			log.Printf("Error checking existing progress for user %d: %v", user.ID, err)
			continue
		}

		if len(existingProgress) > 0 {
			log.Printf("User %d already has module progress, skipping", user.ID)
			continue
		}

		// Initialize user progress (this will create progress for first module)
		err := moduleProgressService.InitializeUserProgress(user.ID)
		if err != nil {
			log.Printf("Failed to initialize module progress for user %d: %v", user.ID, err)
			continue
		}

		successCount++
		log.Printf("Successfully initialized module progress for user %d", user.ID)
	}

	log.Printf("Module progress initialization completed. Success: %d/%d users", successCount, len(users))
}
