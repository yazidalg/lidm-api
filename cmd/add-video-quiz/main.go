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

	// Find video material for module 2
	var videoMaterial models.VideoMaterial
	result := db.Where("module_id = ?", 2).First(&videoMaterial)
	if result.Error != nil {
		log.Fatalf("Video material for module 2 not found: %v", result.Error)
	}

	log.Printf("Found video material for module 2: %s (ID: %d)", videoMaterial.Title, videoMaterial.ID)

	// Convert timestamp 02:07 to seconds (2 minutes 7 seconds = 127 seconds)
	timestampStart := 127
	timestampEnd := timestampStart + 30 // Give 30 seconds to answer

	// Create video quiz
	videoQuiz := models.VideoQuiz{
		VideoMaterialID: videoMaterial.ID,
		Question:        "Apa fungsi utama akar pada tumbuhan?",
		TimestampStart:  timestampStart,
		TimestampEnd:    timestampEnd,
		Options: models.Options{
			OptionA: "Menyerap air dan nutrisi dari tanah",
			OptionB: "Menghasilkan cahaya",
			OptionC: "Menyimpan oksigen",
			OptionD: "Menghindari hama",
		},
		CorrectAnswer: "A",
		Explanation:   "Akar memiliki fungsi utama untuk menyerap air dan nutrisi dari tanah yang diperlukan untuk pertumbuhan dan proses fotosintesis tumbuhan.",
		Order:         1,
	}

	// Save to database
	err := db.Create(&videoQuiz).Error
	if err != nil {
		log.Fatalf("Failed to create video quiz: %v", err)
	}

	log.Printf("Successfully created video quiz for module 2")
	log.Printf("Quiz ID: %d", videoQuiz.ID)
	log.Printf("Question: %s", videoQuiz.Question)
	log.Printf("Timestamp: %d seconds (02:07)", videoQuiz.TimestampStart)
	log.Printf("Video quiz setup completed!")
}
