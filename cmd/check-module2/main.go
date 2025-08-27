package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Query the module progress for Module 2
	var progress models.ModuleProgress
	result := db.Where("user_id = ? AND module_id = ?", 2, 2).First(&progress)
	
	if result.Error != nil {
		fmt.Printf("No progress record found for User ID: 2, Module ID: 2: %v\n", result.Error)
	} else {
		fmt.Printf("Current database values for User ID: 2, Module ID: 2:\n")
		fmt.Printf("Progress: %.2f%%\n", progress.Progress)
		fmt.Printf("IsComplete: %t\n", progress.IsComplete)
		fmt.Printf("IsUnlocked: %t\n", progress.IsUnlocked)
		fmt.Printf("StartedAt: %v\n", progress.StartedAt)
		fmt.Printf("CompletedAt: %v\n", progress.CompletedAt)
	}

	// Check how many prequizzes Module 2 has
	var prequizCount int64
	db.Model(&models.Prequiz{}).Where("module_id = ?", 2).Count(&prequizCount)
	fmt.Printf("\nModule 2 has %d prequizzes\n", prequizCount)

	// Check how many prequizzes user has answered for Module 2
	var answeredCount int64
	db.Table("prequiz_user_answers").
		Joins("JOIN prequizzes ON prequiz_user_answers.prequiz_id = prequizzes.id").
		Where("prequizzes.module_id = ? AND prequiz_user_answers.user_id = ?", 2, 2).
		Count(&answeredCount)
	fmt.Printf("User has answered %d prequizzes for Module 2\n", answeredCount)

	// Check video quizzes for Module 2
	var videoQuizCount int64
	db.Table("video_quizzes").
		Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Where("video_materials.module_id = ?", 2).
		Count(&videoQuizCount)
	fmt.Printf("Module 2 has %d video quizzes\n", videoQuizCount)

	// Check answered video quizzes
	var answeredVideoCount int64
	db.Table("video_quiz_user_answers").
		Joins("JOIN video_quizzes ON video_quiz_user_answers.video_quiz_id = video_quizzes.id").
		Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Where("video_materials.module_id = ? AND video_quiz_user_answers.user_id = ?", 2, 2).
		Count(&answeredVideoCount)
	fmt.Printf("User has answered %d video quizzes for Module 2\n", answeredVideoCount)
}
