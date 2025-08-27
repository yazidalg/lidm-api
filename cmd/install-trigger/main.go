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

	fmt.Println("🔧 Installing Database Trigger...")

	// Drop existing trigger first
	dropSQL := "DROP TRIGGER IF EXISTS auto_unlock_next_module"
	if err := db.Exec(dropSQL).Error; err != nil {
		log.Printf("Error dropping trigger: %v", err)
	}

	// Create trigger
	triggerSQL := `
CREATE TRIGGER auto_unlock_next_module
    AFTER UPDATE ON module_progresses
    FOR EACH ROW
BEGIN
    -- Hanya jalankan jika module baru saja complete (dari false ke true)
    IF NEW.is_complete = 1 AND OLD.is_complete = 0 THEN
        
        -- Cari module selanjutnya berdasarkan ID
        SET @next_module_id = (
            SELECT id 
            FROM modules 
            WHERE id > NEW.module_id 
              AND deleted_at IS NULL 
            ORDER BY id ASC 
            LIMIT 1
        );
        
        -- Jika ada module selanjutnya
        IF @next_module_id IS NOT NULL THEN
            
            -- Check apakah progress record sudah ada
            SET @existing_progress = (
                SELECT id 
                FROM module_progresses 
                WHERE user_id = NEW.user_id 
                  AND module_id = @next_module_id 
                  AND deleted_at IS NULL
            );
            
            -- Jika belum ada, buat record baru
            IF @existing_progress IS NULL THEN
                INSERT INTO module_progresses (
                    user_id, 
                    module_id, 
                    is_unlocked, 
                    is_complete, 
                    progress, 
                    created_at, 
                    updated_at
                ) VALUES (
                    NEW.user_id, 
                    @next_module_id, 
                    1, 
                    0, 
                    0, 
                    NOW(), 
                    NOW()
                );
            ELSE
                -- Jika sudah ada, unlock saja
                UPDATE module_progresses 
                SET is_unlocked = 1, updated_at = NOW() 
                WHERE id = @existing_progress;
            END IF;
            
        END IF;
        
    END IF;
END`

	if err := db.Exec(triggerSQL).Error; err != nil {
		log.Fatal("❌ Failed to create trigger:", err)
	}

	fmt.Println("✅ Trigger installed successfully!")

	// Verify trigger exists
	var triggerName string
	result := db.Raw("SHOW TRIGGERS LIKE 'auto_unlock_next_module'").Scan(&triggerName)
	if result.Error != nil {
		log.Printf("Error checking trigger: %v", result.Error)
	} else if result.RowsAffected > 0 {
		fmt.Println("✅ Trigger verification: Found in database")
	} else {
		fmt.Println("❌ Trigger verification: Not found")
	}
}
