package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type CourseHandler struct {
	courseService services.CourseServiceInterface
}

func NewCourseHandler(courseService services.CourseServiceInterface) *CourseHandler {
	return &CourseHandler{courseService}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var request request.CreateCourseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.courseService.CreateCourse(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create course",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"data":    result,
	})
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.courseService.GetCourseByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Course not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course retrieved successfully",
		"data":    result,
	})
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	results, err := h.courseService.GetAllCourses()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve courses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Courses retrieved successfully",
		"data":    results,
	})
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	// Implementation for updating a course
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	// Implementation for deleting a course
}
