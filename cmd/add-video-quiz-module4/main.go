package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
)

func main() {
	// Load environment and connect to database
	config.LoadEnv()
	db := config.ConnectDB()
	database.Migrate(db)

	// Find video material for module 4
	var videoMaterial models.VideoMaterial
	result := db.Where("module_id = ?", 4).First(&videoMaterial)
	if result.Error != nil {
		log.Fatalf("Video material for module 4 not found: %v", result.Error)
	}

	log.Printf("Found video material for module 4: %s (ID: %d)", videoMaterial.Title, videoMaterial.ID)

	// Create video quizzes for module 4
	videoQuizzes := []models.VideoQuiz{
		{
			VideoMaterialID: videoMaterial.ID,
			Question:        "Apa saja yang dibutuhkan tumbuhan untuk melakukan fotosintesis?",
			TimestampStart:  227, // 03:47 = 3*60 + 47 = 227 seconds
			TimestampEnd:    227 + 30,
			Options: models.Options{
				OptionA: "Air, Karbondioksida, Matahari, Klorofil",
				OptionB: "Oksigen, Air, Karbondioksida, Tanah",
				OptionC: "Air, Matahari, Tanah, Oksigen",
				OptionD: "Air, Klorofil, Nitrogen, Cahaya Bulan",
			},
			CorrectAnswer: "A",
			Explanation:   "Tumbuhan membutuhkan air, karbondioksida, sinar matahari, dan klorofil untuk melakukan fotosintesis. Keempat komponen ini bekerja sama untuk menghasilkan makanan (glukosa) dan oksigen.",
			Order:         1,
		},
		{
			VideoMaterialID: videoMaterial.ID,
			Question:        "Apa hasil utama dari proses fotosintesis?",
			TimestampStart:  248, // 04:08 = 4*60 + 8 = 248 seconds
			TimestampEnd:    248 + 30,
			Options: models.Options{
				OptionA: "Nitrogen dan Air",
				OptionB: "Oksigen dan Karbohidrat (makanan/glukosa)",
				OptionC: "Karbondioksida dan Oksigen",
				OptionD: "Cahaya Matahari dan Air",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis menghasilkan dua produk utama: oksigen yang dilepaskan ke udara dan karbohidrat (glukosa) yang menjadi makanan bagi tumbuhan.",
			Order:         2,
		},
		{
			VideoMaterialID: videoMaterial.ID,
			Question:        "Mengapa fotosintesis penting bagi makhluk hidup di Bumi?",
			TimestampStart:  284, // 04:44 = 4*60 + 44 = 284 seconds
			TimestampEnd:    284 + 30,
			Options: models.Options{
				OptionA: "Karena menghasilkan tanah yang subur",
				OptionB: "Karena menghasilkan oksigen dan makanan untuk rantai makanan",
				OptionC: "Karena membuat daun selalu hijau",
				OptionD: "Karena menyerap semua air hujan",
			},
			CorrectAnswer: "B",
			Explanation:   "Fotosintesis sangat penting karena menghasilkan oksigen yang kita hirup dan makanan (glukosa) yang menjadi dasar rantai makanan di Bumi.",
			Order:         3,
		},
	}

	// Save all video quizzes to database
	for i, quiz := range videoQuizzes {
		err := db.Create(&quiz).Error
		if err != nil {
			log.Fatalf("Failed to create video quiz %d: %v", i+1, err)
		}

		log.Printf("Successfully created video quiz %d for module 4", i+1)
		log.Printf("Quiz ID: %d", quiz.ID)
		log.Printf("Question: %s", quiz.Question)
		log.Printf("Timestamp: %d seconds", quiz.TimestampStart)
		log.Printf("Correct Answer: %s", quiz.CorrectAnswer)
		log.Println("---")
	}

	log.Printf("All video quizzes for module 4 setup completed!")
}
