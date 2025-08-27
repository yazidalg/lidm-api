package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	userID := 2
	moduleID := 1

	fmt.Printf("Checking current quiz answers for User ID: %d, Module ID: %d\n\n", userID, moduleID)

	// Check prequiz answers
	var prequizCount int64
	db.Model(&struct{}{}).Table("prequizzes").Where("module_id = ?", moduleID).Count(&prequizCount)
	
	var answeredPrequizCount int64
	db.Model(&struct{}{}).Table("prequiz_user_answers").
		Joins("JOIN prequizzes ON prequiz_user_answers.prequiz_id = prequizzes.id").
		Where("prequizzes.module_id = ? AND prequiz_user_answers.user_id = ?", moduleID, userID).
		Count(&answeredPrequizCount)

	fmt.Printf("Prequizzes: %d/%d answered\n", answeredPrequizCount, prequizCount)

	// Check video quiz answers
	var videoQuizCount int64
	db.Model(&struct{}{}).Table("video_quizzes").
		Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
		Where("video_materials.module_id = ?", moduleID).
		Count(&videoQuizCount)

	var answeredVideoQuizCount int64
	if videoQuizCount > 0 {
		db.Model(&struct{}{}).Table("video_quiz_user_answers").
			Joins("JOIN video_quizzes ON video_quiz_user_answers.video_quiz_id = video_quizzes.id").
			Joins("JOIN video_materials ON video_quizzes.video_material_id = video_materials.id").
			Where("video_materials.module_id = ? AND video_quiz_user_answers.user_id = ?", moduleID, userID).
			Count(&answeredVideoQuizCount)
	}

	fmt.Printf("Video Quizzes: %d/%d answered\n", answeredVideoQuizCount, videoQuizCount)

	// Calculate expected progress using our new logic
	allPrequizzesAnswered := answeredPrequizCount == prequizCount
	var expectedProgress float32

	if videoQuizCount > 0 {
		// Module has video quizzes
		allVideoQuizzesAnswered := answeredVideoQuizCount == videoQuizCount
		
		if allPrequizzesAnswered && allVideoQuizzesAnswered {
			expectedProgress = 100.0
		} else {
			totalQuizzes := float32(prequizCount + videoQuizCount)
			answeredQuizzes := float32(answeredPrequizCount + answeredVideoQuizCount)
			expectedProgress = (answeredQuizzes / totalQuizzes) * 100.0
		}
	} else {
		// Module has no video quizzes
		if allPrequizzesAnswered {
			expectedProgress = 100.0
		} else {
			if prequizCount > 0 {
				expectedProgress = (float32(answeredPrequizCount) / float32(prequizCount)) * 100.0
			} else {
				expectedProgress = 0.0
			}
		}
	}

	fmt.Printf("\nExpected Progress: %.2f%%\n", expectedProgress)

	// Check if module progress record exists
	var progressExists bool
	result := db.Model(&struct{}{}).Table("module_progresses").
		Where("user_id = ? AND module_id = ? AND deleted_at IS NULL", userID, moduleID).
		Select("1").
		Limit(1)
	
	if result.Error == nil {
		progressExists = true
	}

	fmt.Printf("Progress record exists: %t\n", progressExists)
}
