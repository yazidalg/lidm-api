package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/database"
)

func main() {
	log.Println("🎥 Running Video Quizzes Seeding directly...")
	
	// Call the seeding function
	database.RunVideoQuizzesSeeding()
}
