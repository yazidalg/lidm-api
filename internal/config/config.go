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

// LoadEnv loads .env only in non-production environments
func LoadEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️  No .env file found, skipping...")
		} else {
			log.Println("✅ .env file loaded")
		}
	} else {
		log.Println("🌍 Production environment detected — skipping .env load")
	}
}

// ConnectDB opens a MySQL connection using GORM
func ConnectDB() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	env := os.Getenv("ENV")

	var dsn string

	if env == "production" {
		// -----------------------------
		// Cloud SQL connection (Unix socket)
		// -----------------------------
		instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME") // e.g. "project:region:instance"
		if instanceConnectionName == "" {
			log.Fatal("❌ INSTANCE_CONNECTION_NAME is not set")
		}

		dsn = fmt.Sprintf(
			"%s:%s@unix(/cloudsql/%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			dbUser, dbPassword, instanceConnectionName, dbName,
		)
		log.Printf("🔗 Connecting to Cloud SQL via /cloudsql/%s", instanceConnectionName)

	} else {
		// -----------------------------
		// Local dev (TCP connection)
		// -----------------------------
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		if dbHost == "" {
			dbHost = "127.0.0.1"
		}
		if dbPort == "" {
			dbPort = "3306"
		}

		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		)
		log.Printf("🔗 Connecting to local DB at %s:%s", dbHost, dbPort)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Printf("❌ Database connection failed: %v", err)
		log.Printf("🔍 DSN used: %s", dsn)
		panic(fmt.Sprintf("Failed to connect to MySQL: %v", err))
	}

	log.Println("✅ Database connected successfully")
	return db
}
