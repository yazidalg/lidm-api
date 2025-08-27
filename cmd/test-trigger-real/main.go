package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

	fmt.Println("🧪 Testing Real Trigger Execution...")

	// Find user with module 1 complete but module 2 not unlocked
	userID := 1 // Let's test with user 1
	
	// Check module 1 status
	var module1Progress float64
	var module1Complete bool
	err = db.Raw("SELECT progress, is_complete FROM module_progresses WHERE user_id = ? AND module_id = 1", userID).Row().Scan(&module1Progress, &module1Complete)
	if err != nil {
		fmt.Printf("❌ Error getting module 1 progress: %v\n", err)
		return
	}

	fmt.Printf("📊 Module 1 - Progress: %.2f%%, Complete: %t\n", module1Progress, module1Complete)

	// Check module 2 status
	var module2Unlocked bool
	var module2Exists int
	err = db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = ? AND module_id = 2", userID).Row().Scan(&module2Exists)
	if err != nil {
		fmt.Printf("❌ Error checking module 2: %v\n", err)
		return
	}

	if module2Exists > 0 {
		err = db.Raw("SELECT is_unlocked FROM module_progresses WHERE user_id = ? AND module_id = 2", userID).Row().Scan(&module2Unlocked)
		if err != nil {
			fmt.Printf("❌ Error getting module 2 unlock status: %v\n", err)
			return
		}
		fmt.Printf("🔓 Module 2 - Exists: Yes, Unlocked: %t\n", module2Unlocked)
	} else {
		fmt.Println("🔓 Module 2 - Exists: No")
	}

	// If module 1 is not complete, let's test the trigger
	if !module1Complete {
		fmt.Println("\n🚀 Testing trigger by completing module 1...")
		
		// Update module 1 to complete (this should trigger auto-unlock)
		result := db.Exec("UPDATE module_progresses SET is_complete = 1, progress = 100 WHERE user_id = ? AND module_id = 1", userID)
		if result.Error != nil {
			fmt.Printf("❌ Error updating module 1: %v\n", result.Error)
			return
		}
		
		fmt.Println("✅ Module 1 marked as complete")
		
		// Check if module 2 was automatically created/unlocked
		var module2ExistsAfter int
		err = db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = ? AND module_id = 2", userID).Row().Scan(&module2ExistsAfter)
		if err != nil {
			fmt.Printf("❌ Error checking module 2 after trigger: %v\n", err)
			return
		}
		
		if module2ExistsAfter > 0 {
			var module2UnlockedAfter bool
			err = db.Raw("SELECT is_unlocked FROM module_progresses WHERE user_id = ? AND module_id = 2", userID).Row().Scan(&module2UnlockedAfter)
			if err != nil {
				fmt.Printf("❌ Error getting module 2 status after trigger: %v\n", err)
				return
			}
			
			fmt.Printf("🎉 TRIGGER RESULT: Module 2 - Exists: Yes, Unlocked: %t\n", module2UnlockedAfter)
			
			if module2UnlockedAfter {
				fmt.Println("✅ SUCCESS: Trigger worked! Module 2 automatically unlocked!")
			} else {
				fmt.Println("❌ FAILED: Module 2 exists but not unlocked")
			}
		} else {
			fmt.Println("❌ FAILED: Module 2 was not created by trigger")
		}
	} else {
		fmt.Println("⚠️  Module 1 already complete. Trigger won't fire again.")
		fmt.Println("💡 You can test by manually setting module 1 incomplete then running this again")
	}
}
