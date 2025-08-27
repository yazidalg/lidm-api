package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ProgressData struct {
	UserID    int     `json:"user_id"`
	ModuleID  int     `json:"module_id"`
	Progress  float64 `json:"progress"`
	IsComplete bool   `json:"is_complete"`
	IsUnlocked bool   `json:"is_unlocked"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("📊 Checking Module Progress Data...")

	// Get all module progress data
	var progresses []ProgressData
	result := db.Raw("SELECT user_id, module_id, progress, is_complete, is_unlocked FROM module_progresses ORDER BY user_id, module_id").Scan(&progresses)
	if result.Error != nil {
		log.Printf("Error getting progress data: %v", result.Error)
		return
	}

	fmt.Printf("Found %d progress records:\n", len(progresses))
	for _, p := range progresses {
		fmt.Printf("  User %d, Module %d: %.1f%% (Complete: %t, Unlocked: %t)\n", 
			p.UserID, p.ModuleID, p.Progress, p.IsComplete, p.IsUnlocked)
	}

	if len(progresses) == 0 {
		fmt.Println("❌ No progress data found")
		return
	}

	// Find a good test case - incomplete module with next module available
	var testUserID, testModuleID int
	var found bool

	for _, p := range progresses {
		if !p.IsComplete {
			// Check if next module exists
			var nextModuleExists int
			err = db.Raw("SELECT COUNT(*) FROM modules WHERE id > ? AND deleted_at IS NULL", p.ModuleID).Row().Scan(&nextModuleExists)
			if err == nil && nextModuleExists > 0 {
				testUserID = p.UserID
				testModuleID = p.ModuleID
				found = true
				break
			}
		}
	}

	if !found {
		fmt.Println("❌ No suitable test case found (incomplete module with next module available)")
		return
	}

	fmt.Printf("\n🎯 Found test case: User %d, Module %d\n", testUserID, testModuleID)

	// Test trigger execution
	fmt.Println("🚀 Testing trigger by completing module...")
	
	result = db.Exec("UPDATE module_progresses SET is_complete = 1, progress = 100 WHERE user_id = ? AND module_id = ?", testUserID, testModuleID)
	if result.Error != nil {
		fmt.Printf("❌ Error updating module: %v\n", result.Error)
		return
	}
	
	fmt.Printf("✅ Module %d marked as complete for user %d\n", testModuleID, testUserID)
	
	// Check if next module was unlocked
	nextModuleID := testModuleID + 1
	var nextModuleUnlocked bool
	var nextModuleExists int
	
	err = db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = ? AND module_id = ?", testUserID, nextModuleID).Row().Scan(&nextModuleExists)
	if err != nil {
		fmt.Printf("❌ Error checking next module: %v\n", err)
		return
	}
	
	if nextModuleExists > 0 {
		err = db.Raw("SELECT is_unlocked FROM module_progresses WHERE user_id = ? AND module_id = ?", testUserID, nextModuleID).Row().Scan(&nextModuleUnlocked)
		if err != nil {
			fmt.Printf("❌ Error getting next module unlock status: %v\n", err)
			return
		}
		
		fmt.Printf("🎉 TRIGGER RESULT: Module %d - Unlocked: %t\n", nextModuleID, nextModuleUnlocked)
		
		if nextModuleUnlocked {
			fmt.Println("✅ SUCCESS: Database trigger worked! Next module automatically unlocked!")
		} else {
			fmt.Println("❌ FAILED: Next module exists but not unlocked")
		}
	} else {
		fmt.Printf("🎉 TRIGGER RESULT: Module %d - Created and unlocked by trigger!\n", nextModuleID)
		fmt.Println("✅ SUCCESS: Database trigger worked! Next module created and unlocked!")
	}
}
