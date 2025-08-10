package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type FlashcardHandler struct {
	fsrsService services.FSRSServiceInterface
}

func NewFlashcardHandler(fsrsService services.FSRSServiceInterface) *FlashcardHandler {
	return &FlashcardHandler{
		fsrsService: fsrsService,
	}
}

// ReviewFlashcard - Review a flashcard with FSRS algorithm
func (h *FlashcardHandler) ReviewFlashcard(c *gin.Context) {
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
		Grade int `json:"grade" binding:"required,min=1,max=4"` // 1=Again, 2=Hard, 3=Good, 4=Easy
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

	// Initialize flashcard if not exists, then review
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

	// Review the flashcard
	updatedProgress, err := h.fsrsService.ReviewFlashcard(uid, uint(flashcardID), request.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to review flashcard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard reviewed successfully",
		"data":    updatedProgress,
	})
}

// GetDueFlashcards - Get flashcards due for review
func (h *FlashcardHandler) GetDueFlashcards(c *gin.Context) {
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

	dueFlashcards, err := h.fsrsService.GetDueFlashcards(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch due flashcards",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Due flashcards retrieved successfully",
		"count":   len(dueFlashcards),
		"data":    dueFlashcards,
	})
}

// GetRetentionStats - Get user's retention statistics
func (h *FlashcardHandler) GetRetentionStats(c *gin.Context) {
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

	stats, err := h.fsrsService.GetUserRetentionStats(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch retention statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Retention statistics retrieved successfully",
		"data":    stats,
	})
}

// InitializeFlashcard - Initialize a new flashcard for user
func (h *FlashcardHandler) InitializeFlashcard(c *gin.Context) {
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

	progress, err := h.fsrsService.InitializeFlashcard(uid, uint(flashcardID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize flashcard",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Flashcard initialized successfully",
		"data":    progress,
	})
}
