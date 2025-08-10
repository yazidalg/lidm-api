package services

import (
	"time"

	"github.com/open-spaced-repetition/go-fsrs/v3"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type FSRSServiceInterface interface {
	InitializeFlashcard(userID, flashcardID uint) (*models.UserFlashcardProgress, error)
	ReviewFlashcard(userID, flashcardID uint, grade int) (*models.UserFlashcardProgress, error)
	GetDueFlashcards(userID uint) ([]models.UserFlashcardProgress, error)
	GetFlashcardProgress(userID, flashcardID uint) (*models.UserFlashcardProgress, error)
	GetUserRetentionStats(userID uint) (map[string]interface{}, error)
}

type fsrsService struct {
	flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface
	fsrsAlgorithm         *fsrs.FSRS
}

func NewFSRSService(flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface) FSRSServiceInterface {
	// Initialize FSRS with default parameters
	params := fsrs.DefaultParam()
	algorithm := fsrs.NewFSRS(params)

	return &fsrsService{
		flashcardProgressRepo: flashcardProgressRepo,
		fsrsAlgorithm:         algorithm,
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

	// Update average grade
	if progress.ReviewCount == 1 {
		progress.AverageGrade = float64(grade)
	} else {
		progress.AverageGrade = (progress.AverageGrade*float64(progress.ReviewCount-1) + float64(grade)) / float64(progress.ReviewCount)
	}

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
