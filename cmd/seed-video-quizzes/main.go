package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
)

func main() {
	log.Println("=== SEED VIDEO QUIZZES TOOL ===")
	
	// Load environment variables
	config.LoadEnv()
	
	// Connect to database
	db := config.ConnectDB()
	
	log.Println("Connected to database successfully!")
	
	// Check command line arguments
	if len(os.Args) < 2 {
		showUsage()
		return
	}
	
	command := os.Args[1]
	
	switch command {
	case "all":
		// Seed video quizzes for all modules
		database.SeedVideoQuizzes(db)
		
	case "module":
		if len(os.Args) < 3 {
			log.Println("Error: module name required")
			log.Println("Usage: go run main.go module \"Module Name\"")
			return
		}
		moduleTitle := os.Args[2]
		database.SeedVideoQuizzesForModule(db, moduleTitle)
		
	case "clear":
		// Clear all video quizzes
		database.ClearVideoQuizzes(db)
		
	case "summary":
		// Show summary of current video quizzes
		database.GetVideoQuizzesSummary(db)
		
	case "fotosintesis":
		// Seed all fotosintesis modules
		modules := []string{
			"Fotosintesis - Dasar",
			"Fotosintesis - Lanjutan", 
			"Fotosintesis - Eksperimen",
		}
		for _, module := range modules {
			database.SeedVideoQuizzesForModule(db, module)
		}
		
	default:
		showUsage()
	}
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go all                              - Seed video quizzes for all modules")
	fmt.Println("  go run main.go module \"Module Name\"             - Seed video quizzes for specific module")
	fmt.Println("  go run main.go fotosintesis                     - Seed video quizzes for all fotosintesis modules")
	fmt.Println("  go run main.go clear                            - Clear all video quizzes")
	fmt.Println("  go run main.go summary                          - Show summary of current video quizzes")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run main.go module \"Fotosintesis - Dasar\"")
	fmt.Println("  go run main.go module \"Fotosintesis - Lanjutan\"")
	fmt.Println("  go run main.go module \"Fotosintesis - Eksperimen\"")
}
