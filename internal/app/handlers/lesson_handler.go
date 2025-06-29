package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type LessonHandler struct {
	lessonService services.LessonServiceInterface
	moduleService services.ModuleServiceInterface
}

func NewLessonHandler(
	lessonService services.LessonServiceInterface,
	moduleService services.ModuleServiceInterface,
) *LessonHandler {
	return &LessonHandler{
		lessonService: lessonService,
		moduleService: moduleService,
	}
}

func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var lessonRequest request.LessonRequest

	module, err := h.moduleService.GetModuleByID(uint32(lessonRequest.ModuleID))

	if err != nil && module == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Module not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&lessonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := h.lessonService.CreateLesson(lessonRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create lesson",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created successfully",
		"data":    result,
	})
}

func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	result, err := h.lessonService.GetLessonByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Lesson not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson retrieved successfully",
		"data":    result,
	})
}

func (h *LessonHandler) GetAllLessons(c *gin.Context) {
	results, err := h.lessonService.GetAllLessons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve lessons",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons retrieved successfully",
		"data":    results,
	})
}

func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	lesson, err := h.lessonService.GetLessonByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Lesson not found",
		})
		return
	}

	var lessonRequest request.LessonRequest
	if err := c.ShouldBindJSON(&lessonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.lessonService.UpdateLesson(uint32(lesson.ID), lessonRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update lesson",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson updated successfully",
		"data":    result,
	})
}

func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}

	lesson, err := h.lessonService.GetLessonByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Lesson not found",
		})
		return
	}

	err = h.lessonService.DeleteLesson(uint32(lesson.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete lesson",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson deleted successfully",
	})
}
