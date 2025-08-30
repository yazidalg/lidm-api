package main

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	fmt.Printf("Unlocking Module 2 and setting up Module 3...\n")

	// 1. Unlock Module 2 for User 2
	err := db.Exec("UPDATE module_progresses SET is_unlocked = true WHERE user_id = 2 AND module_id = 2").Error
	if err != nil {
		log.Fatal("Failed to unlock Module 2:", err)
	}
	fmt.Printf("✅ Module 2 unlocked for User 2\n")

	// 2. Check if Module 3 progress exists, if not create it
	var count int64
	db.Model(&struct{ ID uint }{}).Table("module_progresses").Where("user_id = 2 AND module_id = 3").Count(&count)

	if count == 0 {
		// Create Module 3 progress entry as unlocked (since Module 2 is complete)
		err = db.Exec(`
			INSERT INTO module_progresses (user_id, module_id, is_unlocked, is_complete, progress, created_at, updated_at) 
			VALUES (2, 3, true, false, 0, NOW(), NOW())
		`).Error
		if err != nil {
			log.Fatal("Failed to create Module 3 progress:", err)
		}
		fmt.Printf("✅ Module 3 created and unlocked for User 2\n")
	} else {
		// Update existing Module 3 to be unlocked
		err = db.Exec("UPDATE module_progresses SET is_unlocked = true WHERE user_id = 2 AND module_id = 3").Error
		if err != nil {
			log.Fatal("Failed to unlock Module 3:", err)
		}
		fmt.Printf("✅ Module 3 unlocked for User 2\n")
	}

	fmt.Printf("\n🎉 All done! Modules 2 and 3 are now properly set up.\n")
}
