package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type PrequizHandler struct {
	prequizService services.PrequizServiceInterface
	lessonService  services.LessonServiceInterface
	userService    services.UserServiceInterface
}

func NewPrequizHandler(
	prequizService services.PrequizServiceInterface,
	lessonService services.LessonServiceInterface,
	userService services.UserServiceInterface,
) *PrequizHandler {
	return &PrequizHandler{
		prequizService: prequizService,
		lessonService:  lessonService,
		userService:    userService,
	}
}

func (h *PrequizHandler) CreatePrequiz(c *gin.Context) {
	var request request.PrequizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	lesson, lessonErr := h.lessonService.GetLessonByID(uint32(request.LessonID))
	user, userErr := h.userService.GetUserById(int(request.UserID))
	if lessonErr != nil || userErr != nil || lesson == nil || user.ID == 0 || lesson.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Lesson or user not found",
			"error":   "Failed to create prequiz",
		})
		return
	}

	result, err := h.prequizService.CreatePrequiz(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create prequiz",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Prequiz created successfully",
		"data":    result,
	})
}

func (h *PrequizHandler) GetPrequizByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prequiz ID"})
		return
	}

	prequiz, err := h.prequizService.GetPrequizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Prequiz not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequiz retrieved successfully",
		"data":    prequiz,
	})
}

func (h *PrequizHandler) GetAllPrequizzes(c *gin.Context) {
	prequizzes, err := h.prequizService.GetAllPrequizzes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve prequizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequizzes retrieved successfully",
		"data":    prequizzes,
	})
}

func (h *PrequizHandler) UpdatePrequiz(c *gin.Context) {
	var request request.PrequizRequest
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prequiz ID"})
		return
	}

	lesson, lessonErr := h.lessonService.GetLessonByID(uint32(request.LessonID))
	user, userErr := h.userService.GetUserById(int(request.UserID))
	if lessonErr != nil || userErr != nil || lesson == nil || user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Lesson or user not found",
			"error":   "Failed to update prequiz",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	prequiz, err := h.prequizService.UpdatePrequiz(uint(id), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update prequiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prequiz updated successfully",
		"data":    prequiz,
	})
}
