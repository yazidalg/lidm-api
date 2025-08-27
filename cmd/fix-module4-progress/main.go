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

	fmt.Printf("Fixing Module 4 progress for user who answered all quizzes...\n")
	
	// From API response: User answered 3/3 prequizzes + 3/3 video quizzes = 6/6 = 100%
	now := time.Now()
	
	err := db.Exec(`
		UPDATE module_progresses 
		SET progress = 100.0, 
			is_complete = true, 
			completed_at = ?, 
			updated_at = ? 
		WHERE user_id = 2 AND module_id = 4
	`, now, now).Error
	
	if err != nil {
		fmt.Printf("Error updating Module 4: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Module 4 updated to 100%% and marked complete\n")
	
	// Unlock Module 5 since Module 4 is now complete
	var count int64
	db.Model(&struct{ ID uint }{}).Table("module_progresses").Where("user_id = 2 AND module_id = 5").Count(&count)
	
	if count == 0 {
		// Create Module 5 progress entry as unlocked
		err = db.Exec(`
			INSERT INTO module_progresses (user_id, module_id, is_unlocked, is_complete, progress, created_at, updated_at) 
			VALUES (2, 5, true, false, 0, NOW(), NOW())
		`).Error
		if err != nil {
			fmt.Printf("Error creating Module 5: %v\n", err)
			return
		}
		fmt.Printf("✅ Module 5 created and unlocked\n")
	} else {
		// Update existing Module 5 to be unlocked
		err = db.Exec("UPDATE module_progresses SET is_unlocked = true WHERE user_id = 2 AND module_id = 5").Error
		if err != nil {
			fmt.Printf("Error unlocking Module 5: %v\n", err)
			return
		}
		fmt.Printf("✅ Module 5 unlocked\n")
	}
	
	fmt.Printf("\n🎉 Module 4 progress fixed from 50%% to 100%%!\n")
}
