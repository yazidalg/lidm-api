package main

import (
	"fmt"
	"time"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	fmt.Printf("Creating proper progress records for User 2...\n")

	now := time.Now()

	// Create Module 2 progress (3/3 prequizzes + 1/1 video quiz = 100%)
	err := db.Exec(`
		INSERT INTO module_progresses (user_id, module_id, is_unlocked, is_complete, progress, completed_at, created_at, updated_at) 
		VALUES (2, 2, true, true, 100.0, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			is_unlocked = true,
			is_complete = true,
			progress = 100.0,
			completed_at = ?,
			updated_at = NOW()
	`, now, now).Error

	if err != nil {
		fmt.Printf("Error creating Module 2: %v\n", err)
	} else {
		fmt.Printf("✅ Module 2: 100%% complete, unlocked\n")
	}

	// Create Module 3 progress (0/3 prequizzes = 0% but unlocked since Module 2 is complete)
	err = db.Exec(`
		INSERT INTO module_progresses (user_id, module_id, is_unlocked, is_complete, progress, created_at, updated_at) 
		VALUES (2, 3, true, false, 0.0, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			is_unlocked = true,
			updated_at = NOW()
	`).Error

	if err != nil {
		fmt.Printf("Error creating Module 3: %v\n", err)
	} else {
		fmt.Printf("✅ Module 3: 0%% progress, unlocked\n")
	}

	fmt.Printf("\n🎉 Progress records created successfully!\n")
}
