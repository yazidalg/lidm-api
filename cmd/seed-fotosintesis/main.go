package main

import (
	"log"

	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
)

func main() {
	log.Println("🌱 Starting Fotosintesis Data Seeder...")

	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Run Fotosintesis seeding
	database.SeedFotosintesisData(db)

	log.Println("✅ Fotosintesis Data Seeding completed!")
}
