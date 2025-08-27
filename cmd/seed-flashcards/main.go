package main

import (
	"fmt"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
	"gorm.io/gorm"
)

type FlashcardData struct {
	ModuleID  uint
	FrontText string
	BackText  string
}

func main() {
	// Load environment variables
	config.LoadEnv()
	
	// Initialize database connection
	db := config.ConnectDB()

	fmt.Println("🚀 Starting flashcard seeding...")

	// Define all flashcards data
	flashcards := []FlashcardData{
		// Module 1 Flashcards
		{ModuleID: 1, FrontText: "Fotosintesis", BackText: "Proses tumbuhan membuat makanan sendiri dengan cahaya matahari, air, dan karbon dioksida."},
		{ModuleID: 1, FrontText: "Bahasa Yunani dari fotosintesis", BackText: "\"photo\" = cahaya, \"synthesis\" = menyusun atau membuat."},
		{ModuleID: 1, FrontText: "Bagian tumbuhan membuat makanannya", BackText: "Di daun, tepatnya di bagian yang disebut kloroplas."},
		{ModuleID: 1, FrontText: "Tujuan tumbuhan membuat makanan sendiri", BackText: "Agar bisa tumbuh, membuat bunga dan buah, serta bertahan hidup."},
		{ModuleID: 1, FrontText: "Bahan yang dibutuhkan tumbuhan untuk fotosintesis", BackText: "Cahaya matahari, air, karbon dioksida, dan klorofil."},
		{ModuleID: 1, FrontText: "Akibat jika salah satu bahan fotosintesis tidak ada", BackText: "Tumbuhan tidak bisa membuat makanan dan bisa layu atau berhenti tumbuh."},

		// Module 2 Flashcards
		{ModuleID: 2, FrontText: "Bagian tumbuhan yang disebut \"dapur makanan tumbuhan\"", BackText: "Daun"},
		{ModuleID: 2, FrontText: "Zat hijau pada daun yang membantu fotosintesis", BackText: "Klorofil"},
		{ModuleID: 2, FrontText: "Fungsi klorofil", BackText: "Menangkap cahaya matahari untuk membuat energi makanan."},
		{ModuleID: 2, FrontText: "Lubang kecil di daun yang berguna untuk keluar masuk gas", BackText: "Stomata"},
		{ModuleID: 2, FrontText: "Gas apa yang masuk ke daun melalui stomata", BackText: "Karbon dioksida (CO₂)."},
		{ModuleID: 2, FrontText: "Gas yang keluar dari daun setelah fotosintesis", BackText: "Oksigen (O₂)."},

		// Module 3 Flashcards
		{ModuleID: 3, FrontText: "Akar", BackText: "Bagian tumbuhan untuk menyerap air dan mineral"},
		{ModuleID: 3, FrontText: "Gas dari udara yang dibutuhkan tumbuhan untuk fotosintesis", BackText: "Karbon dioksida (CO₂)."},
		{ModuleID: 3, FrontText: "Tempat keluar masuknya gas O₂ dan CO₂", BackText: "Stomata daun."},
		{ModuleID: 3, FrontText: "Suhu yang menghambat fotosintesis", BackText: "Terlalu dingin atau panas."},
		{ModuleID: 3, FrontText: "Zat hijau daun", BackText: "Klorofil."},
		{ModuleID: 3, FrontText: "Fungsi klorofil", BackText: "Menangkap cahaya matahari."},

		// Module 4 Flashcards
		{ModuleID: 4, FrontText: "Yang diserap akar di tanah", BackText: "Air dan mineral"},
		{ModuleID: 4, FrontText: "Kloroplas", BackText: "Tempat di dalam daun tempat terjadinya fotosintesis."},
		{ModuleID: 4, FrontText: "Hasil utama fotosintesis", BackText: "Glukosa dan oksigen."},
		{ModuleID: 4, FrontText: "Fungsi Jaringan Floem", BackText: "Mengedarkan hasil fotosintesis ke seluruh bagian tumbuhan"},
		{ModuleID: 4, FrontText: "Glukosa berlebih disimpan jadi?", BackText: "Pati (biji, buah, umbi)."},
		{ModuleID: 4, FrontText: "Pati", BackText: "Bentuk cadangan makanan yang disimpan di biji, buah, atau umbi."},

		// Module 6 Flashcards
		{ModuleID: 6, FrontText: "Rumus kimia glukosa", BackText: "C₆H₁₂O₆"},
		{ModuleID: 6, FrontText: "Rumus kimia oksigen", BackText: "O₂"},
		{ModuleID: 6, FrontText: "Fungsi glukosa untuk tumbuhan", BackText: "Memperbesar akar, batang, dan daun."},
		{ModuleID: 6, FrontText: "Hasil utama fotosintesis", BackText: "Glukosa dan oksigen."},
		{ModuleID: 6, FrontText: "Alasan oksigen penting untuk manusia dan hewan", BackText: "Supaya manusia dan hewan bisa bernapas."},
		{ModuleID: 6, FrontText: "Alasan oksigen penting untuk lingkungan", BackText: "Menjaga udara tetap bersih dan segar"},
	}

	// Check and create flashcards
	for i, flashcardData := range flashcards {
		fmt.Printf("📚 Creating flashcard %d/%d for Module %d...\n", i+1, len(flashcards), flashcardData.ModuleID)

		// Check if module exists
		var module models.Module
		if err := db.First(&module, flashcardData.ModuleID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("⚠️  Module %d not found, skipping flashcard: %s\n", flashcardData.ModuleID, flashcardData.FrontText)
				continue
			}
			log.Printf("Error checking module %d: %v", flashcardData.ModuleID, err)
			continue
		}

		// Check if flashcard already exists (avoid duplicates)
		var existingFlashcard models.Flashcard
		err := db.Where("module_id = ? AND front_text = ?", flashcardData.ModuleID, flashcardData.FrontText).
			First(&existingFlashcard).Error

		if err == nil {
			fmt.Printf("⚠️  Flashcard already exists: %s\n", flashcardData.FrontText)
			continue
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("Error checking existing flashcard: %v", err)
			continue
		}

		// Create new flashcard with auto-increment order
		var lastOrder int
		db.Model(&models.Flashcard{}).Where("module_id = ?", flashcardData.ModuleID).
			Select("COALESCE(MAX(`order`), 0)").Scan(&lastOrder)

		newFlashcard := models.Flashcard{
			ModuleID:  flashcardData.ModuleID,
			FrontText: flashcardData.FrontText,
			BackText:  flashcardData.BackText,
			Order:     lastOrder + 1,
		}

		if err := db.Create(&newFlashcard).Error; err != nil {
			log.Printf("❌ Failed to create flashcard for Module %d: %v", flashcardData.ModuleID, err)
			continue
		}

		fmt.Printf("✅ Created flashcard: %s\n", flashcardData.FrontText)
	}

	fmt.Println("🎉 Flashcard seeding completed!")

	// Print summary
	fmt.Println("\n📊 Summary:")
	for moduleID := uint(1); moduleID <= 6; moduleID++ {
		var count int64
		db.Model(&models.Flashcard{}).Where("module_id = ?", moduleID).Count(&count)
		if count > 0 {
			fmt.Printf("Module %d: %d flashcards\n", moduleID, count)
		}
	}

	fmt.Println("\n🔥 All flashcards have been successfully added to the database!")
	fmt.Println("You can now use them in your frontend application.")
}
