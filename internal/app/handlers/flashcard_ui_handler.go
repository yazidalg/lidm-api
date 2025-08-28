package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type FlashcardUIHandler struct {
	flashcardUIService services.FlashcardUIServiceInterface
	fsrsService       services.FSRSServiceInterface
}

func NewFlashcardUIHandler(
	flashcardUIService services.FlashcardUIServiceInterface, 
	fsrsService services.FSRSServiceInterface,
) *FlashcardUIHandler {
	return &FlashcardUIHandler{
		flashcardUIService: flashcardUIService,
		fsrsService:       fsrsService,
	}
}

// GetFlashcardWithOptions - Get a single flashcard with review options (1m, 5m, 7h, 10h)
func (h *FlashcardUIHandler) GetFlashcardWithOptions(c *gin.Context) {
	flashcardIDParam := c.Param("flashcard_id")
	flashcardID, err := strconv.ParseUint(flashcardIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flashcard ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Convert userID to uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get preview of next review schedule
	schedule, err := h.flashcardUIService.GetNextReviewSchedule(uid, uint(flashcardID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get review schedule",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard review options",
		"data": gin.H{
			"flashcard_id": flashcardID,
			"options": gin.H{
				"ulang": gin.H{
					"grade":           1,
					"time_display":    "1m",
					"color":          "#ef4444",
					"description":    "Ulangi",
				},
				"sulit": gin.H{
					"grade":           2,
					"time_display":    "5m", 
					"color":          "#f59e0b",
					"description":    "Sulit",
				},
				"lumayan": gin.H{
					"grade":           3,
					"time_display":    "7h",
					"color":          "#10b981", 
					"description":    "Lumayan",
				},
				"mudah": gin.H{
					"grade":           4,
					"time_display":    "10h",
					"color":          "#3b82f6",
					"description":    "Mudah",
				},
			},
			"schedule": schedule,
		},
	})
}

// ReviewFlashcardWithGrade - Review flashcard and save to progress
func (h *FlashcardUIHandler) ReviewFlashcardWithGrade(c *gin.Context) {
	flashcardIDParam := c.Param("flashcard_id")
	flashcardID, err := strconv.ParseUint(flashcardIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flashcard ID",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var request struct {
		Grade int `json:"grade" binding:"required,min=1,max=4"` // 1=Ulang, 2=Sulit, 3=Lumayan, 4=Mudah
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert userID to uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Check if flashcard progress exists, if not initialize it
	_, err = h.fsrsService.GetFlashcardProgress(uid, uint(flashcardID))
	if err != nil {
		// Initialize new flashcard
		_, err = h.fsrsService.InitializeFlashcard(uid, uint(flashcardID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to initialize flashcard",
			})
			return
		}
	}

	// Review the flashcard and save to progress
	updatedProgress, err := h.fsrsService.ReviewFlashcard(uid, uint(flashcardID), request.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to review flashcard",
		})
		return
	}

	// Format response with time display
	var timeDisplay, description, color string
	switch request.Grade {
	case 1:
		timeDisplay = "1m"
		description = "Ulang"
		color = "#ef4444"
	case 2:
		timeDisplay = "5m"
		description = "Sulit"
		color = "#f59e0b"
	case 3:
		timeDisplay = "7h"
		description = "Lumayan"
		color = "#10b981"
	case 4:
		timeDisplay = "10h"
		description = "Mudah"
		color = "#3b82f6"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard reviewed successfully",
		"data": gin.H{
			"flashcard_id":   flashcardID,
			"grade":          request.Grade,
			"time_display":   timeDisplay,
			"description":    description,
			"color":          color,
			"next_review":    updatedProgress.Due,
			"progress":       updatedProgress,
		},
	})
}

// GetUserFlashcardProgress - Get user's flashcard progress (from saved data)
func (h *FlashcardUIHandler) GetUserFlashcardProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Get module ID from query parameter (optional)
	moduleIDParam := c.Query("module_id")
	var moduleID *uint
	if moduleIDParam != "" {
		if mid, err := strconv.ParseUint(moduleIDParam, 10, 32); err == nil {
			modID := uint(mid)
			moduleID = &modID
		}
	}

	// Convert userID to uint
	var uid uint
	switch v := userID.(type) {
	case uint:
		uid = v
	case float64:
		uid = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid = uint(id)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user ID",
			})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get flashcards with intervals
	flashcards, err := h.flashcardUIService.GetFlashcardsWithIntervals(uid, moduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get flashcard progress",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User flashcard progress retrieved",
		"data":    flashcards,
	})
}
