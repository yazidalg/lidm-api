package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	fmt.Printf("Fixing Module 3 progress...\n")

	// Module 3 has 3 prequizzes, 0 video quizzes
	// User answered all 3 prequizzes = 100% complete
	now := time.Now()

	err := db.Exec(`
		UPDATE module_progresses 
		SET progress = 100.0, 
			is_complete = true, 
			completed_at = ?, 
			updated_at = ? 
		WHERE user_id = 2 AND module_id = 3
	`, now, now).Error

	if err != nil {
		log.Fatal("Failed to update Module 3 progress:", err)
	}

	fmt.Printf("✅ Module 3 updated to 100%% and marked complete\n")

	// Unlock Module 4 since Module 3 is now complete
	var count int64
	db.Model(&struct{ ID uint }{}).Table("module_progresses").Where("user_id = 2 AND module_id = 4").Count(&count)

	if count == 0 {
		// Create Module 4 progress entry as unlocked
		err = db.Exec(`
			INSERT INTO module_progresses (user_id, module_id, is_unlocked, is_complete, progress, created_at, updated_at) 
			VALUES (2, 4, true, false, 0, NOW(), NOW())
		`).Error
		if err != nil {
			log.Fatal("Failed to create Module 4 progress:", err)
		}
		fmt.Printf("✅ Module 4 created and unlocked\n")
	} else {
		// Update existing Module 4 to be unlocked
		err = db.Exec("UPDATE module_progresses SET is_unlocked = true WHERE user_id = 2 AND module_id = 4").Error
		if err != nil {
			log.Fatal("Failed to unlock Module 4:", err)
		}
		fmt.Printf("✅ Module 4 unlocked\n")
	}

	fmt.Printf("\n🎉 Module 3 progress fixed!\n")
}
