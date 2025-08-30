package services

import (
	"fmt"
	"time"

	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type FSRSServiceInterface interface {
	InitializeFlashcard(userID, flashcardID uint) (*models.UserFlashcardProgress, error)
	InitializeModuleFlashcards(userID, moduleID uint) (int, error) // Returns count of initialized flashcards
	ReviewFlashcard(userID, flashcardID uint, grade int) (*models.UserFlashcardProgress, error)
	GetDueFlashcards(userID uint) ([]models.UserFlashcardProgress, error)
	GetFlashcardProgress(userID, flashcardID uint) (*models.UserFlashcardProgress, error)
	GetUserRetentionStats(userID uint) (map[string]interface{}, error)
	GetNextReviewSchedule(userID, flashcardID uint, currentState int) (map[string]interface{}, error) // Preview next review times
	IsFlashcardDue(userID, flashcardID uint) (bool, error) // Check if flashcard is due for review
	GetFlashcardReviewStats(userID, flashcardID uint) (map[string]int, error) // Get review statistics (u, s, l, m counts)
}

type fsrsService struct {
	flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface
	moduleRepo           repositories.ModuleRepositoryInterface
	fsrsAlgorithm        *fsrs.FSRS
}

func NewFSRSService(flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface, moduleRepo repositories.ModuleRepositoryInterface) FSRSServiceInterface {
	// Initialize FSRS with custom parameters
	params := fsrs.DefaultParam()
	params.RequestRetention = 0.98 // Set retention rate to 98%
	
	algorithm := fsrs.NewFSRS(params)

	return &fsrsService{
		flashcardProgressRepo: flashcardProgressRepo,
		moduleRepo:           moduleRepo,
		fsrsAlgorithm:        algorithm,
	}
}

func (s *fsrsService) InitializeFlashcard(userID, flashcardID uint) (*models.UserFlashcardProgress, error) {
	// Check if progress already exists
	existing, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new card using FSRS
	card := fsrs.NewCard()
	now := time.Now()

	// Set Due to current time for new flashcards (immediately available for review)
	if card.Due.IsZero() {
		card.Due = now
	}

	progress := &models.UserFlashcardProgress{
		UserID:      userID,
		FlashcardID: flashcardID,
		Stability:   card.Stability,
		Difficulty:  card.Difficulty,
		Elapsed:     int(card.ElapsedDays),
		Scheduled:   int(card.ScheduledDays),
		Reps:        int(card.Reps),
		Lapses:      int(card.Lapses),
		State:       int(card.State),
		Due:         card.Due,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.flashcardProgressRepo.Create(progress)
}

// InitializeModuleFlashcards initializes all flashcards in a module for a user
func (s *fsrsService) InitializeModuleFlashcards(userID, moduleID uint) (int, error) {
	// Get module with flashcards
	module, err := s.moduleRepo.GetModuleByID(uint32(moduleID))
	if err != nil {
		return 0, err
	}

	if module == nil {
		return 0, fmt.Errorf("module not found")
	}

	initializedCount := 0
	
	for _, flashcard := range module.Flashcards {
		// Check if this flashcard is already initialized for this user
		existing, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcard.ID)
		if err == nil && existing != nil {
			// Already exists, skip
			continue
		}

		// Initialize this flashcard for the user
		_, err = s.InitializeFlashcard(userID, flashcard.ID)
		if err != nil {
			// Log error but continue with other flashcards
			continue
		}
		
		initializedCount++
	}

	return initializedCount, nil
}

func (s *fsrsService) ReviewFlashcard(userID, flashcardID uint, grade int) (*models.UserFlashcardProgress, error) {
	// Get current progress
	progress, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
	if err != nil {
		return nil, err
	}

	// Convert to FSRS card
	card := s.progressToCard(progress)

	// Create review log
	reviewTime := time.Now()
	rating := fsrs.Rating(grade)

	// Process review using FSRS algorithm
	schedulingCards := s.fsrsAlgorithm.Repeat(card, reviewTime)

	// Get the updated card based on rating
	var updatedCard fsrs.Card
	switch rating {
	case fsrs.Again:
		updatedCard = schedulingCards[rating].Card
	case fsrs.Hard:
		updatedCard = schedulingCards[rating].Card
	case fsrs.Good:
		updatedCard = schedulingCards[rating].Card
	case fsrs.Easy:
		updatedCard = schedulingCards[rating].Card
	default:
		updatedCard = schedulingCards[fsrs.Good].Card // Default to Good
	}

	// Update progress with new values
	progress.Stability = updatedCard.Stability
	progress.Difficulty = updatedCard.Difficulty
	progress.Elapsed = int(updatedCard.ElapsedDays)
	progress.Scheduled = int(updatedCard.ScheduledDays)
	progress.Reps = int(updatedCard.Reps)
	progress.Lapses = int(updatedCard.Lapses)
	progress.State = int(updatedCard.State)
	progress.Due = updatedCard.Due
	progress.LastReview = &reviewTime
	progress.ReviewCount++

	// Update average grade using UNIQUE system - only store the last grade
	// This supports the unique l/m/s/u statistics where each flashcard has only one active status
	fmt.Printf("DEBUG: Setting grade %d for flashcard %d (was %f)\n", grade, flashcardID, progress.AverageGrade)
	progress.AverageGrade = float64(grade) // Store last grade directly

	return s.flashcardProgressRepo.Update(progress)
}

func (s *fsrsService) GetDueFlashcards(userID uint) ([]models.UserFlashcardProgress, error) {
	return s.flashcardProgressRepo.GetDueByUser(userID, time.Now())
}

func (s *fsrsService) GetFlashcardProgress(userID, flashcardID uint) (*models.UserFlashcardProgress, error) {
	return s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
}

func (s *fsrsService) GetUserRetentionStats(userID uint) (map[string]interface{}, error) {
	allProgress, err := s.flashcardProgressRepo.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_cards":        len(allProgress),
		"new_cards":          0,
		"learning_cards":     0,
		"review_cards":       0,
		"average_stability":  0.0,
		"average_difficulty": 0.0,
		"total_reviews":      0,
		"average_grade":      0.0,
	}

	if len(allProgress) == 0 {
		return stats, nil
	}

	var totalStability, totalDifficulty, totalGrade float64
	var totalReviews int

	for _, progress := range allProgress {
		switch progress.State {
		case models.StateNew:
			stats["new_cards"] = stats["new_cards"].(int) + 1
		case models.StateLearning, models.StateRelearning:
			stats["learning_cards"] = stats["learning_cards"].(int) + 1
		case models.StateReview:
			stats["review_cards"] = stats["review_cards"].(int) + 1
		}

		totalStability += progress.Stability
		totalDifficulty += progress.Difficulty
		totalReviews += progress.ReviewCount
		totalGrade += progress.AverageGrade * float64(progress.ReviewCount)
	}

	stats["average_stability"] = totalStability / float64(len(allProgress))
	stats["average_difficulty"] = totalDifficulty / float64(len(allProgress))
	stats["total_reviews"] = totalReviews

	if totalReviews > 0 {
		stats["average_grade"] = totalGrade / float64(totalReviews)
	}

	return stats, nil
}

// Helper function to convert UserFlashcardProgress to FSRS Card
func (s *fsrsService) progressToCard(progress *models.UserFlashcardProgress) fsrs.Card {
	card := fsrs.Card{
		Due:           progress.Due,
		Stability:     progress.Stability,
		Difficulty:    progress.Difficulty,
		ElapsedDays:   uint64(progress.Elapsed),
		ScheduledDays: uint64(progress.Scheduled),
		Reps:          uint64(progress.Reps),
		Lapses:        uint64(progress.Lapses),
		State:         fsrs.State(progress.State),
	}

	if progress.LastReview != nil {
		card.LastReview = *progress.LastReview
	}

	return card
}

// GetNextReviewSchedule gets preview of next review schedule for all grades
func (s *fsrsService) GetNextReviewSchedule(userID, flashcardID uint, currentState int) (map[string]interface{}, error) {
	// Get current progress or create default
	progress, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
	if err != nil {
		// Use default new card if no progress exists
		card := fsrs.NewCard()
		return s.generateSchedulePreview(card), nil
	}

	// Convert progress to FSRS card
	card := s.progressToCard(progress)
	return s.generateSchedulePreview(card), nil
}

// Helper function to generate schedule preview for all grades
func (s *fsrsService) generateSchedulePreview(card fsrs.Card) map[string]interface{} {
	now := time.Now()
	
	// Calculate schedule for each grade
	againResult := s.fsrsAlgorithm.Next(card, now, fsrs.Again)
	hardResult := s.fsrsAlgorithm.Next(card, now, fsrs.Hard)
	goodResult := s.fsrsAlgorithm.Next(card, now, fsrs.Good)
	easyResult := s.fsrsAlgorithm.Next(card, now, fsrs.Easy)

	return map[string]interface{}{
		"scheduleTimers": map[string]interface{}{
			"ulang": map[string]interface{}{
				"grade":           1,
				"due":             againResult.Card.Due,
				"due_display":     s.formatDueDisplay(againResult.Card.Due, now),
				"stability":       againResult.Card.Stability,
				"difficulty":      againResult.Card.Difficulty,
				"elapsed_days":    againResult.Card.ElapsedDays,
				"scheduled_days":  againResult.Card.ScheduledDays,
				"reps":            againResult.Card.Reps,
				"lapses":          againResult.Card.Lapses,
				"state":           againResult.Card.State,
				"description":     "Ulangi",
				"color":           "#ef4444",
			},
			"sulit": map[string]interface{}{
				"grade":           2,
				"due":             hardResult.Card.Due,
				"due_display":     s.formatDueDisplay(hardResult.Card.Due, now),
				"stability":       hardResult.Card.Stability,
				"difficulty":      hardResult.Card.Difficulty,
				"elapsed_days":    hardResult.Card.ElapsedDays,
				"scheduled_days":  hardResult.Card.ScheduledDays,
				"reps":            hardResult.Card.Reps,
				"lapses":          hardResult.Card.Lapses,
				"state":           hardResult.Card.State,
				"description":     "Sulit",
				"color":           "#f59e0b",
			},
			"lumayan": map[string]interface{}{
				"grade":           3,
				"due":             goodResult.Card.Due,
				"due_display":     s.formatDueDisplay(goodResult.Card.Due, now),
				"stability":       goodResult.Card.Stability,
				"difficulty":      goodResult.Card.Difficulty,
				"elapsed_days":    goodResult.Card.ElapsedDays,
				"scheduled_days":  goodResult.Card.ScheduledDays,
				"reps":            goodResult.Card.Reps,
				"lapses":          goodResult.Card.Lapses,
				"state":           goodResult.Card.State,
				"description":     "Lumayan",
				"color":           "#10b981",
			},
			"mudah": map[string]interface{}{
				"grade":           4,
				"due":             easyResult.Card.Due,
				"due_display":     s.formatDueDisplay(easyResult.Card.Due, now),
				"stability":       easyResult.Card.Stability,
				"difficulty":      easyResult.Card.Difficulty,
				"elapsed_days":    easyResult.Card.ElapsedDays,
				"scheduled_days":  easyResult.Card.ScheduledDays,
				"reps":            easyResult.Card.Reps,
				"lapses":          easyResult.Card.Lapses,
				"state":           easyResult.Card.State,
				"description":     "Mudah",
				"color":           "#3b82f6",
			},
		},
	}
}

// Helper function to format due time display
func (s *fsrsService) formatDueDisplay(due time.Time, now time.Time) string {
	duration := due.Sub(now)
	
	if duration < 0 {
		return "now"
	}
	
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes < 1 {
			return "now"
		}
		return fmt.Sprintf("%dm", minutes)
	}
	
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	}
	
	days := int(duration.Hours() / 24)
	return fmt.Sprintf("%dd", days)
}

// IsFlashcardDue checks if a flashcard is due for review
func (s *fsrsService) IsFlashcardDue(userID, flashcardID uint) (bool, error) {
	// Get flashcard progress
	progress, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
	if err != nil {
		// If no progress found, it means flashcard has never been reviewed, so it's "due"
		return true, nil
	}

	// Check if current time is past the due time
	now := time.Now()
	return now.After(progress.Due), nil
}

// GetFlashcardReviewStats returns review statistics for a flashcard
// This uses UNIQUE logic - each flashcard can only have one active review status
func (s *fsrsService) GetFlashcardReviewStats(userID, flashcardID uint) (map[string]int, error) {
	// Get flashcard progress
	progress, err := s.flashcardProgressRepo.GetByUserAndFlashcard(userID, flashcardID)
	if err != nil {
		// If no progress found, return all zeros (never reviewed)
		return map[string]int{
			"u": 0, // ulang (grade 1)
			"s": 0, // sulit (grade 2)
			"l": 0, // lumayan (grade 3)
			"m": 0, // mudah (grade 4)
		}, nil
	}

	// Initialize stats with all zeros
	stats := map[string]int{
		"u": 0, // ulang (grade 1)
		"s": 0, // sulit (grade 2)  
		"l": 0, // lumayan (grade 3)
		"m": 0, // mudah (grade 4)
	}

	// UNIQUE LOGIC: Each flashcard has only ONE active status based on last review
	if progress.ReviewCount > 0 {
		// Use the last grade to determine current status
		// Since we store last grade directly in AverageGrade field
		
		lastGrade := progress.AverageGrade
		
		// Debug logging for flashcard 12
		if flashcardID == 12 {
			fmt.Printf("DEBUG flashcard 12: reviewCount=%d, lastGrade=%f\n", progress.ReviewCount, lastGrade)
		}
		
		// Determine category based on last review grade
		if lastGrade <= 1.5 {
			stats["u"] = 1 // This flashcard is currently in "ulang" status
		} else if lastGrade <= 2.5 {
			stats["s"] = 1 // This flashcard is currently in "sulit" status
		} else if lastGrade <= 3.5 {
			stats["l"] = 1 // This flashcard is currently in "lumayan" status
		} else {
			stats["m"] = 1 // This flashcard is currently in "mudah" status
		}
	}

	return stats, nil
}
