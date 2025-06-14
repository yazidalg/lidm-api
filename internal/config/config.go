package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func LoadEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // default to dev if ENV is not set
	}

	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, skipping...")
		} else {
			log.Println(".env file loaded")
		}
	} else {
		log.Println("Production environment, skipping .env load")
	}
}

func MigrateDb(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Quiz{}, &models.Participant{}, &models.Question{}, &models.Answer{}, &models.Leaderboard{})
	db.Migrator().DropColumn(&models.Participant{}, "score")
	// db.Migrator().DropTable(&models.User{}, &models.Quiz{}, &models.Participant{}, &models.Question{}, &models.Answer{}, &models.Leaderboard{})
}

func ConnectDB() *gorm.DB {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Format: user:password@tcp(host:port)/dbname?parseTime=true&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to MySQL database: %v", err))
	}

	return db
}
