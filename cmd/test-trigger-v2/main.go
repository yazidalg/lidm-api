package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TriggerInfo struct {
	Trigger string `json:"trigger"`
	Event   string `json:"event"`
	Table   string `json:"table"`
	Timing  string `json:"timing"`
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

	fmt.Println("🧪 Testing Database Trigger (v2)...")

	// Test 1: Show all triggers
	var triggers []TriggerInfo
	result := db.Raw("SHOW TRIGGERS").Scan(&triggers)
	if result.Error != nil {
		log.Printf("Error getting triggers: %v", result.Error)
	} else {
		fmt.Printf("📋 Found %d triggers in database:\n", len(triggers))
		for _, trigger := range triggers {
			fmt.Printf("  - %s (%s %s on %s)\n", trigger.Trigger, trigger.Timing, trigger.Event, trigger.Table)
		}
	}

	// Test 2: Check specific trigger with different approach
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM information_schema.triggers WHERE trigger_name = 'auto_unlock_next_module'").Scan(&count).Error
	if err != nil {
		log.Printf("Error checking trigger via information_schema: %v", err)
	} else {
		fmt.Printf("🔍 Trigger count in information_schema: %d\n", count)
	}

	// Test 3: Try to simulate trigger execution
	fmt.Println("\n🚀 Testing trigger functionality...")
	
	// Get a test user and module
	var userID, moduleID int
	err = db.Raw("SELECT user_id, module_id FROM module_progresses WHERE is_complete = 0 LIMIT 1").Row().Scan(&userID, &moduleID)
	if err != nil {
		fmt.Println("❌ No incomplete module progress found for testing")
		return
	}

	fmt.Printf("📝 Test data: User %d, Module %d\n", userID, moduleID)

	// Check if next module exists
	var nextModuleID int
	err = db.Raw("SELECT id FROM modules WHERE id > ? AND deleted_at IS NULL ORDER BY id ASC LIMIT 1", moduleID).Row().Scan(&nextModuleID)
	if err != nil {
		fmt.Println("❌ No next module found for testing")
		return
	}

	fmt.Printf("➡️  Next module would be: %d\n", nextModuleID)

	// Check current state
	var currentProgress float64
	var isComplete bool
	err = db.Raw("SELECT progress, is_complete FROM module_progresses WHERE user_id = ? AND module_id = ?", userID, moduleID).Row().Scan(&currentProgress, &isComplete)
	if err != nil {
		fmt.Printf("❌ Error getting current progress: %v\n", err)
		return
	}

	fmt.Printf("📊 Current state: Progress=%.2f%%, Complete=%t\n", currentProgress, isComplete)

	if isComplete {
		fmt.Println("⚠️  Module already complete, trigger won't fire")
		return
	}

	// Check if next module progress exists
	var nextProgressExists int
	err = db.Raw("SELECT COUNT(*) FROM module_progresses WHERE user_id = ? AND module_id = ?", userID, nextModuleID).Row().Scan(&nextProgressExists)
	if err != nil {
		fmt.Printf("❌ Error checking next module progress: %v\n", err)
		return
	}

	fmt.Printf("🔗 Next module progress exists: %s\n", map[bool]string{true: "Yes", false: "No"}[nextProgressExists > 0])

	fmt.Println("\n✅ Trigger setup verification complete!")
	fmt.Println("💡 To test trigger, complete a module via the API and check if next module unlocks automatically")
}
