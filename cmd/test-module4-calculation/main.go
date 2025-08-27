package main

import (
	"fmt"

	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/config"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to database
	db := config.ConnectDB()

	// Initialize repositories
	moduleProgressRepo := repositories.NewModuleProgressRepository(db)
	userRepo := repositories.NewUserRepository(db)
	moduleRepo := repositories.NewModuleRepository(db)
	prequizRepo := repositories.NewPrequizRepository(db)
	videoQuizRepo := repositories.NewVideoQuizRepository(db)

	// Initialize service
	moduleProgressService := services.NewModuleProgressService(
		moduleProgressRepo,
		userRepo,
		moduleRepo,
		prequizRepo,
		videoQuizRepo,
	)

	fmt.Printf("Testing Module 4 progress calculation...\n")
	
	// Test the calculation method directly
	calculatedProgress, err := moduleProgressService.UpdateUserProgress(2, 4)
	if err != nil {
		fmt.Printf("Error calculating progress: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Calculated Progress: %.2f%%\n", calculatedProgress.Progress)
	fmt.Printf("✅ Is Complete: %t\n", calculatedProgress.IsComplete)
	fmt.Printf("✅ Is Unlocked: %t\n", calculatedProgress.IsUnlocked)
	
	// Manual verification
	fmt.Printf("\nManual verification:\n")
	
	// Check prequizzes
	allPrequizzes, _ := prequizRepo.GetPrequizzesByModule(4, 0)
	prequizAnswers, _ := prequizRepo.GetUserPrequizAnswers(2)
	
	answeredPrequizMap := make(map[uint]bool)
	for _, answer := range prequizAnswers {
		answeredPrequizMap[answer.PrequizID] = true
	}
	
	answeredPrequizzes := 0
	for _, prequiz := range allPrequizzes {
		if answeredPrequizMap[prequiz.ID] {
			answeredPrequizzes++
		}
	}
	
	fmt.Printf("- Module 4 Prequizzes: %d/%d answered\n", answeredPrequizzes, len(allPrequizzes))
	
	// Check video quizzes
	module, _ := moduleRepo.GetModuleByID(4)
	if module.VideoMaterial != nil && len(module.VideoMaterial.VideoQuizzes) > 0 {
		videoAnswers, _ := videoQuizRepo.GetAllUserVideoQuizAnswers(2)
		
		answeredVideoMap := make(map[uint]bool)
		for _, answer := range videoAnswers {
			answeredVideoMap[answer.VideoQuizID] = true
		}
		
		answeredVideoQuizzes := 0
		for _, videoQuiz := range module.VideoMaterial.VideoQuizzes {
			if answeredVideoMap[videoQuiz.ID] {
				answeredVideoQuizzes++
			}
		}
		
		fmt.Printf("- Module 4 Video Quizzes: %d/%d answered\n", answeredVideoQuizzes, len(module.VideoMaterial.VideoQuizzes))
		
		totalQuizzes := len(allPrequizzes) + len(module.VideoMaterial.VideoQuizzes)
		totalAnswered := answeredPrequizzes + answeredVideoQuizzes
		expectedProgress := float32(totalAnswered) / float32(totalQuizzes) * 100.0
		
		fmt.Printf("- Expected Progress: %.2f%% (%d/%d total)\n", expectedProgress, totalAnswered, totalQuizzes)
	}
}
