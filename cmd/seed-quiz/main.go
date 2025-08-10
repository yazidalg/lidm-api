package main

import (
	"fmt"
	"os"

	"github.com/yazidalg/lidm_backend/internal/config"
	"github.com/yazidalg/lidm_backend/internal/database"
)

// Usage:
// go run cmd/seed-quiz/main.go "Fotosintesis" "Ekosistem"
// or build: go build -o seed-quiz cmd/seed-quiz/main.go && ./seed-quiz "Fotosintesis"
func main() {
	config.LoadEnv()
	db := config.ConnectDB()

	if len(os.Args) < 2 {
		fmt.Println("Provide at least one module title. Example: go run cmd/seed-quiz/main.go \"Fotosintesis\"")
		os.Exit(1)
	}
	modules := os.Args[1:]
	database.SeedQuizData(db, modules)
	fmt.Println("Quiz seeding finished.")
}
