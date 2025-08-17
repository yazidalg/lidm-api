package database

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

// SeedQuizData seeds sample quizzes (single & multiplayer) plus participants and optional questions per module.
// Safe to run multiple times; skips if quizzes already exist for target module(s).
func SeedQuizData(db *gorm.DB, moduleTitles []string) {
	rand.Seed(time.Now().UnixNano())

	for _, title := range moduleTitles {
		var module models.Module
		if err := db.Where("title = ?", title).First(&module).Error; err != nil {
			log.Printf("SeedQuiz: module '%s' not found, skipping", title)
			continue
		}

		// Check existing quizzes
		var existing int64
		db.Model(&models.Quiz{}).Where("module_id = ?", module.ID).Count(&existing)
		if existing > 0 {
			log.Printf("SeedQuiz: quizzes already exist for module '%s' (id=%d), skip", title, module.ID)
			continue
		}

		// Get any two users (host + opponent) for multiplayer sample
		var users []models.User
		db.Limit(2).Find(&users)
		if len(users) == 0 {
			log.Printf("SeedQuiz: no users to attach, skipping module '%s'", title)
			continue
		}

		hostID := users[0].ID

		// Create single player quiz sample
		spQuiz := models.Quiz{
			Status:        "completed",
			Mode:          "single_player",
			ModuleID:      &module.ID,
			HostUserID:    hostID,
			QuestionCount: 10,
		}
		if err := db.Create(&spQuiz).Error; err != nil {
			log.Printf("SeedQuiz: error creating single_player quiz: %v", err)
		}

		// Create multiplayer quiz lobby (pending)
		mpQuiz := models.Quiz{
			Status:        "pending",
			Mode:          "multiplayer",
			ModuleID:      &module.ID,
			HostUserID:    hostID,
			QuestionCount: 10,
		}
		if err := db.Create(&mpQuiz).Error; err != nil {
			log.Printf("SeedQuiz: error creating multiplayer quiz: %v", err)
		}

		// Attach participants (host + optional second)
		_ = db.Create(&models.Participant{UserID: hostID, QuizID: mpQuiz.ID}).Error
		if len(users) == 2 {
			_ = db.Create(&models.Participant{UserID: users[1].ID, QuizID: mpQuiz.ID}).Error
		}

		// Seed sample questions if module has < 10 questions
		var qCount int64
		db.Model(&models.Question{}).Where("module_id = ?", module.ID).Count(&qCount)
		needed := 10 - int(qCount)
		for i := 0; i < needed; i++ {
			modID := module.ID
			q := models.Question{
				ModuleID:      &modID,
				Question:      "Sample question " + title + " #" + strconv.Itoa(i+1),
				AnswerTime:    5,
				ReadTime:      5,
				Options:       models.Options{OptionA: "A", OptionB: "B", OptionC: "C", OptionD: "D"},
				CorrectAnswer: "A",
				Explanation:   "Sample explanation",
			}
			if err := db.Create(&q).Error; err != nil {
				log.Printf("SeedQuiz: create question err: %v", err)
			}
		}

		log.Printf("SeedQuiz: seeded quizzes for module '%s'", title)
	}
}
