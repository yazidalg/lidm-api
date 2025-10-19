package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

func ConnectDB() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	instanceConnectionName := os.Getenv("DB_HOST") // this is your instance connection name
	env := os.Getenv("ENV")

	var dsn string
	if env == "production" {
		// Use Cloud SQL Unix socket
		dsn = fmt.Sprintf(
			"%s:%s@unix(/cloudsql/%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			dbUser, dbPassword, instanceConnectionName, dbName,
		)
	} else {
		// Local development connection
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}

	return db
}

