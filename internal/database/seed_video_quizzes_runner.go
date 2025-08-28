package database

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func RunVideoQuizzesSeeding() {
	log.Println("🎥 Starting Video Quizzes Seeding...")

	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	log.Println("Connected to database successfully!")

	// Clear existing video quizzes first (optional)
	log.Println("Clearing existing video quizzes...")
	ClearVideoQuizzes(db)

	// Seed all video quizzes for all modules
	log.Println("Starting video quizzes seeding for all modules...")
	SeedVideoQuizzes(db)

	// Show summary
	log.Println("Showing video quizzes summary...")
	GetVideoQuizzesSummary(db)

	log.Println("✅ Video Quizzes Seeding completed!")
}
