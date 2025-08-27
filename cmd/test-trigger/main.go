package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	fmt.Printf("🧪 Testing Database Trigger...\n")
	
	// Check current triggers
	rows, err := db.Raw("SHOW TRIGGERS LIKE 'auto_unlock_next_module'").Rows()
	if err != nil {
		fmt.Printf("Error checking triggers: %v\n", err)
		return
	}
	defer rows.Close()
	
	hasRows := false
	for rows.Next() {
		hasRows = true
		fmt.Printf("✅ Trigger 'auto_unlock_next_module' is installed\n")
	}
	
	if !hasRows {
		fmt.Printf("❌ Trigger 'auto_unlock_next_module' not found\n")
		return
	}
	
	// Test manual trigger simulation
	fmt.Printf("\n🔄 Simulating module completion to test trigger...\n")
	
	// Before state
	var beforeCount int64
	db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = 2 AND module_id = 6").Scan(&beforeCount)
	fmt.Printf("Module 6 progress records before: %d\n", beforeCount)
	
	// Simulate completing Module 5 (should trigger unlock of Module 6)
	err = db.Exec(`
		UPDATE module_progresses 
		SET is_complete = 1, progress = 100, updated_at = NOW()
		WHERE user_id = 2 AND module_id = 5
	`).Error
	
	if err != nil {
		fmt.Printf("Error updating Module 5: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Module 5 marked as complete\n")
	
	// After state
	var afterCount int64
	db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = 2 AND module_id = 6").Scan(&afterCount)
	fmt.Printf("Module 6 progress records after: %d\n", afterCount)
	
	if afterCount > beforeCount {
		fmt.Printf("🎉 Trigger worked! Module 6 was auto-unlocked\n")
		
		// Check the details
		var isUnlocked bool
		db.Raw("SELECT is_unlocked FROM module_progresses WHERE user_id = 2 AND module_id = 6").Scan(&isUnlocked)
		fmt.Printf("Module 6 is_unlocked: %t\n", isUnlocked)
	} else {
		fmt.Printf("❌ Trigger didn't work or Module 6 already existed\n")
	}
}
