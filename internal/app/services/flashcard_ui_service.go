package services

import (
	"fmt"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

// FlashcardUIService provides UI-friendly flashcard operations
type FlashcardUIServiceInterface interface {
	GetFlashcardsWithIntervals(userID uint, moduleID *uint) ([]FlashcardWithInterval, error)
	GetNextReviewSchedule(userID uint, flashcardID uint, grade int) (ReviewSchedule, error)
	InitializeModuleFlashcards(userID uint, moduleID uint) error
}

type flashcardUIService struct {
	fsrsService           FSRSServiceInterface
	flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface
}

func NewFlashcardUIService(
	fsrsService FSRSServiceInterface,
	flashcardProgressRepo repositories.FlashcardProgressRepositoryInterface,
) FlashcardUIServiceInterface {
	return &flashcardUIService{
		fsrsService:           fsrsService,
		flashcardProgressRepo: flashcardProgressRepo,
	}
}

// FlashcardWithInterval represents a flashcard with user-friendly interval display
type FlashcardWithInterval struct {
	ID                uint                 `json:"id"`
	FrontText         string               `json:"front_text"`
	BackText          string               `json:"back_text"`
	ModuleID          uint                 `json:"module_id"`
	ModuleTitle       string               `json:"module_title"`
	Progress          *models.UserFlashcardProgress `json:"progress,omitempty"`
	IntervalDisplay   string               `json:"interval_display"`   // e.g., "1m", "5m", "7h", "10h"
	NextReviewTime    *time.Time           `json:"next_review_time"`
	IsDue             bool                 `json:"is_due"`
	State             string               `json:"state"`               // "new", "learning", "review", "relearning"
	DifficultyLevel   string               `json:"difficulty_level"`    // "easy", "normal", "hard"
	ReviewCount       int                  `json:"review_count"`
	AverageGrade      float64              `json:"average_grade"`
}

// ReviewSchedule shows the next review intervals for different grades
type ReviewSchedule struct {
	Again ReviewOption `json:"again"` // Grade 1
	Hard  ReviewOption `json:"hard"`  // Grade 2
	Good  ReviewOption `json:"good"`  // Grade 3
	Easy  ReviewOption `json:"easy"`  // Grade 4
}

type ReviewOption struct {
	Grade           int        `json:"grade"`
	IntervalDisplay string     `json:"interval_display"`
	NextReviewTime  time.Time  `json:"next_review_time"`
	Color          string     `json:"color"`           // UI color coding
	Description    string     `json:"description"`     // User-friendly description
}

// GetFlashcardsWithIntervals returns flashcards with user-friendly interval displays
func (s *flashcardUIService) GetFlashcardsWithIntervals(userID uint, moduleID *uint) ([]FlashcardWithInterval, error) {
	// Get all user's flashcard progress
	allProgress, err := s.flashcardProgressRepo.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	var flashcards []FlashcardWithInterval

	// Create a map for quick lookup
	progressMap := make(map[uint]*models.UserFlashcardProgress)
	for i := range allProgress {
		progressMap[allProgress[i].FlashcardID] = &allProgress[i]
	}

	// Process existing progress data
	for _, progress := range allProgress {
		// Filter by moduleID if specified
		if moduleID != nil && progress.Flashcard.ModuleID != *moduleID {
			continue
		}

		flashcardItem := FlashcardWithInterval{
			ID:          progress.Flashcard.ID,
			FrontText:   progress.Flashcard.FrontText,
			BackText:    progress.Flashcard.BackText,
			ModuleID:    progress.Flashcard.ModuleID,
			Progress:    &progress,
			IntervalDisplay: formatInterval(progress.Due),
			NextReviewTime:  &progress.Due,
			IsDue:          progress.IsDue(),
			State:          getStateDescription(progress.State),
			DifficultyLevel: getDifficultyLevel(progress.Difficulty),
			ReviewCount:    progress.ReviewCount,
			AverageGrade:   progress.AverageGrade,
		}

		flashcards = append(flashcards, flashcardItem)
	}

	return flashcards, nil
}

// GetNextReviewSchedule returns the next review schedule for different grades
func (s *flashcardUIService) GetNextReviewSchedule(userID uint, flashcardID uint, grade int) (ReviewSchedule, error) {
	// This would require extending the FSRS service to preview schedules
	// For now, return a simplified version
	baseTime := time.Now()
	
	schedule := ReviewSchedule{
		Again: ReviewOption{
			Grade:           1,
			IntervalDisplay: "1m",
			NextReviewTime:  baseTime.Add(1 * time.Minute),
			Color:          "#ef4444", // red
			Description:    "Ulangi",
		},
		Hard: ReviewOption{
			Grade:           2,
			IntervalDisplay: "5m",
			NextReviewTime:  baseTime.Add(5 * time.Minute),
			Color:          "#f59e0b", // orange
			Description:    "Sulit",
		},
		Good: ReviewOption{
			Grade:           3,
			IntervalDisplay: "7h",
			NextReviewTime:  baseTime.Add(7 * time.Hour),
			Color:          "#10b981", // green
			Description:    "Lumayan",
		},
		Easy: ReviewOption{
			Grade:           4,
			IntervalDisplay: "10h",
			NextReviewTime:  baseTime.Add(10 * time.Hour),
			Color:          "#3b82f6", // blue
			Description:    "Mudah",
		},
	}

	return schedule, nil
}

// InitializeModuleFlashcards initializes all flashcards in a module for a user
func (s *flashcardUIService) InitializeModuleFlashcards(userID uint, moduleID uint) error {
	// This would require getting all flashcards in a module
	// For now, return nil as FSRS service handles initialization on first review
	return nil
}

// Helper functions

// formatInterval converts a due time to user-friendly interval display
func formatInterval(dueTime time.Time) string {
	now := time.Now()
	
	// If overdue, show as "Due"
	if dueTime.Before(now) {
		return "Due"
	}
	
	duration := dueTime.Sub(now)
	
	// Convert to different units based on duration
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes < 1 {
			return "Now"
		}
		return fmt.Sprintf("%dm", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
}

// getStateDescription converts state number to user-friendly description
func getStateDescription(state int) string {
	switch state {
	case models.StateNew:
		return "new"
	case models.StateLearning:
		return "learning"
	case models.StateReview:
		return "review"
	case models.StateRelearning:
		return "relearning"
	default:
		return "unknown"
	}
}

// getDifficultyLevel converts difficulty score to user-friendly level
func getDifficultyLevel(difficulty float64) string {
	if difficulty < 3.0 {
		return "easy"
	} else if difficulty < 7.0 {
		return "normal"
	} else {
		return "hard"
	}
}

// GetDueFlashcardsCount returns the count of due flashcards per module
func (s *flashcardUIService) GetDueFlashcardsCount(userID uint) (map[uint]int, error) {
	dueFlashcards, err := s.fsrsService.GetDueFlashcards(userID)
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int)
	for _, progress := range dueFlashcards {
		if progress.Flashcard.ModuleID > 0 {
			counts[progress.Flashcard.ModuleID]++
		}
	}

	return counts, nil
}
