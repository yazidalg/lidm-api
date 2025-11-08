package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type FlashcardHandler struct {
	fsrsService   services.FSRSServiceInterface
	moduleRepo    repositories.ModuleRepositoryInterface
	flashcardRepo repositories.FlashcardRepositoryInterface
}

func NewFlashcardHandler(fsrsService services.FSRSServiceInterface, moduleRepo repositories.ModuleRepositoryInterface, flashcardRepo repositories.FlashcardRepositoryInterface) *FlashcardHandler {
	return &FlashcardHandler{
		fsrsService:   fsrsService,
		moduleRepo:    moduleRepo,
		flashcardRepo: flashcardRepo,
	}
}

// GetAllFlashcards - Get all flashcards
func (h *FlashcardHandler) GetAllFlashcards(c *gin.Context) {
	flashcards, err := h.flashcardRepo.GetAllFlashcards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch flashcards",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All flashcards retrieved successfully",
		"count":   len(flashcards),
		"data":    flashcards,
	})
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

	// Get user-friendly time display based on grade
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
			"progress":     updatedProgress,
			"grade":        request.Grade,
			"time_display": timeDisplay,
			"description":  description,
			"color":        color,
		},
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

// GetFlashcardIntervals - Get flashcard intervals by module or individual flashcard
func (h *FlashcardHandler) GetFlashcardIntervals(c *gin.Context) {
	flashcardIDParam := c.Query("flashcard_id")
	moduleIDParam := c.Query("module_id")

	// If no parameters provided, return default intervals
	if flashcardIDParam == "" && moduleIDParam == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Default flashcard review intervals",
			"data": gin.H{
				"options": gin.H{
					"ulang": gin.H{
						"grade":       1,
						"time":        "1m",
						"description": "Ulangi",
						"color":       "#ef4444",
					},
					"sulit": gin.H{
						"grade":       2,
						"time":        "5m",
						"description": "Sulit",
						"color":       "#f59e0b",
					},
					"lumayan": gin.H{
						"grade":       3,
						"time":        "7h",
						"description": "Lumayan",
						"color":       "#10b981",
					},
					"mudah": gin.H{
						"grade":       4,
						"time":        "10h",
						"description": "Mudah",
						"color":       "#3b82f6",
					},
				},
				"usage": "Use the grade (1-4) when calling the review endpoint",
			},
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

	// Handle module-based intervals
	if moduleIDParam != "" {
		moduleID, err := strconv.ParseUint(moduleIDParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid module ID",
			})
			return
		}

		// Get all flashcards in the module
		allFlashcards, err := h.moduleRepo.GetFlashcardsByModule(uint(moduleID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get flashcards for module",
			})
			return
		}

		// Filter to only show flashcards that are due or never reviewed
		dueFlashcardsInModule := make([]models.Flashcard, 0)

		for _, flashcard := range allFlashcards {
			// Check if this flashcard is due for this user
			isDue, err := h.fsrsService.IsFlashcardDue(uid, flashcard.ID)
			if err != nil {
				// If error checking due status, assume it's new (never reviewed)
				dueFlashcardsInModule = append(dueFlashcardsInModule, flashcard)
				continue
			}

			// Only include if it's due (includes never reviewed flashcards)
			if isDue {
				dueFlashcardsInModule = append(dueFlashcardsInModule, flashcard)
			}
		}

		// Generate schedule preview for each due flashcard (using array instead of map)
		flashcardSchedules := make([]interface{}, 0)

		// Track overall module review statistics from ALL flashcards in module (not just due ones)
		moduleReviewStats := map[string]int{
			"u": 0, // total ulang reviews
			"s": 0, // total sulit reviews
			"l": 0, // total lumayan reviews
			"m": 0, // total mudah reviews
		}

		// Calculate module stats from ALL flashcards in the module
		for _, flashcard := range allFlashcards {
			// Get review statistics for this flashcard and add to module total
			reviewStats, err := h.fsrsService.GetFlashcardReviewStats(uid, flashcard.ID)
			if err == nil {
				moduleReviewStats["u"] += reviewStats["u"]
				moduleReviewStats["s"] += reviewStats["s"]
				moduleReviewStats["l"] += reviewStats["l"]
				moduleReviewStats["m"] += reviewStats["m"]
			}
		}

		// Generate schedule preview only for due flashcards
		for _, flashcard := range dueFlashcardsInModule {
			schedule, err := h.fsrsService.GetNextReviewSchedule(uid, flashcard.ID, 0)
			if err != nil {
				// If error, use basic info
				flashcardSchedules = append(flashcardSchedules, gin.H{
					"front_text":   flashcard.FrontText,
					"back_text":    flashcard.BackText,
					"order":        flashcard.Order,
					"flashcard_id": flashcard.ID,
					"error":        "Failed to get schedule",
				})
				continue
			}

			flashcardSchedules = append(flashcardSchedules, gin.H{
				"front_text":     flashcard.FrontText,
				"back_text":      flashcard.BackText,
				"order":          flashcard.Order,
				"flashcard_id":   flashcard.ID,
				"scheduleTimers": schedule["scheduleTimers"],
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     fmt.Sprintf("Module %d flashcard intervals", moduleID),
			"module_id":   moduleID,
			"count":       len(allFlashcards), // Total flashcards in module, not just due ones
			"reviewStats": moduleReviewStats,  // Overall module review statistics from ALL flashcards
			"data":        flashcardSchedules,
		})
		return
	}

	// Handle single flashcard interval (existing logic)
	flashcardID, err := strconv.ParseUint(flashcardIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flashcard ID",
		})
		return
	}

	// Get schedule preview
	schedule, err := h.fsrsService.GetNextReviewSchedule(uid, uint(flashcardID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get schedule preview",
		})
		return
	}

	// Get review statistics
	reviewStats, err := h.fsrsService.GetFlashcardReviewStats(uid, uint(flashcardID))
	if err != nil {
		// If error getting stats, set to default
		reviewStats = map[string]int{
			"u": 0, // ulang
			"s": 0, // sulit
			"l": 0, // lumayan
			"m": 0, // mudah
		}
	}

	// Add review stats to the response
	scheduleWithStats := schedule
	if scheduleData, ok := schedule["data"].(map[string]interface{}); ok {
		scheduleData["reviewStats"] = reviewStats
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard schedule preview",
		"data":    scheduleWithStats,
	})
}

// CreateFlashcard - Create a new flashcard (Admin/Teacher only)
func (h *FlashcardHandler) CreateFlashcard(c *gin.Context) {
	var req request.CreateFlashcardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Validate that the module exists
	_, err := h.moduleRepo.GetModuleByID(uint32(req.ModuleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Module not found",
			"details": "The specified module does not exist",
		})
		return
	}

	// Create flashcard model
	flashcard := &models.Flashcard{
		ModuleID:  req.ModuleID,
		FrontText: req.FrontText,
		BackText:  req.BackText,
		Order:     req.Order,
	}

	// Create flashcard
	createdFlashcard, err := h.flashcardRepo.CreateFlashcard(flashcard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create flashcard",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Flashcard created successfully",
		"data":    createdFlashcard,
	})
}

// UpdateFlashcard - Update a flashcard (Admin/Teacher only)
func (h *FlashcardHandler) UpdateFlashcard(c *gin.Context) {
	flashcardIDParam := c.Param("flashcard_id")
	flashcardID, err := strconv.ParseUint(flashcardIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flashcard ID",
		})
		return
	}

	// Check if flashcard exists
	existingFlashcard, err := h.flashcardRepo.GetFlashcardByID(uint(flashcardID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Flashcard not found",
		})
		return
	}

	var req request.UpdateFlashcardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update only provided fields
	if req.ModuleID != nil {
		existingFlashcard.ModuleID = *req.ModuleID
	}
	if req.FrontText != nil {
		existingFlashcard.FrontText = *req.FrontText
	}
	if req.BackText != nil {
		existingFlashcard.BackText = *req.BackText
	}
	if req.Order != nil {
		existingFlashcard.Order = *req.Order
	}

	// Update flashcard
	updatedFlashcard, err := h.flashcardRepo.UpdateFlashcard(uint(flashcardID), existingFlashcard)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update flashcard",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard updated successfully",
		"data":    updatedFlashcard,
	})
}

// DeleteFlashcard - Delete a flashcard (Admin/Teacher only)
func (h *FlashcardHandler) DeleteFlashcard(c *gin.Context) {
	flashcardIDParam := c.Param("flashcard_id")
	flashcardID, err := strconv.ParseUint(flashcardIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flashcard ID",
		})
		return
	}

	// Check if flashcard exists
	_, err = h.flashcardRepo.GetFlashcardByID(uint(flashcardID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Flashcard not found",
		})
		return
	}

	// Delete flashcard
	err = h.flashcardRepo.DeleteFlashcard(uint(flashcardID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete flashcard",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flashcard deleted successfully",
	})
}

// InitializeModuleFlashcards - Copy/initialize all flashcards in a module for user
func (h *FlashcardHandler) InitializeModuleFlashcards(c *gin.Context) {
	moduleIDParam := c.Param("module_id")
	moduleID, err := strconv.ParseUint(moduleIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid module ID",
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

	// Initialize all flashcards in the module
	initializedCount, err := h.fsrsService.InitializeModuleFlashcards(uid, uint(moduleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize module flashcards",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Module flashcards initialized successfully",
		"data": gin.H{
			"module_id":         moduleID,
			"user_id":           uid,
			"initialized_count": initializedCount,
		},
	})
}
