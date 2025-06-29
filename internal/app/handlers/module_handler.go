package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type ModuleHandler struct {
	moduleService services.ModuleServiceInterface
}

func NewModuleHandler(moduleService services.ModuleServiceInterface) *ModuleHandler {
	return &ModuleHandler{moduleService}
}

func (h *ModuleHandler) CreateModule(c *gin.Context) {
	var moduleRequest request.ModuleRequest

	if err := c.ShouldBindJSON(&moduleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := h.moduleService.CreateModule(moduleRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create module",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Module created successfully",
		"data":    result,
	})
}

func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	result, err := h.moduleService.GetModuleByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Module not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module retrieved successfully",
		"data":    result,
	})
}

func (h *ModuleHandler) GetAllModules(c *gin.Context) {
	results, err := h.moduleService.GetAllModules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve modules",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved successfully",
		"data":    results,
	})
}

func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var request request.ModuleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	moduleId, err := h.moduleService.GetModuleByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Module not found",
		})
		return
	}

	result, err := h.moduleService.UpdateModule(uint32(moduleId.ID), request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update module",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module updated successfully",
		"data":    result,
	})
}

func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	moduleId, err := h.moduleService.GetModuleByID(uint32(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Module not found",
		})
		return
	}

	err = h.moduleService.DeleteModule(uint32(moduleId.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete module",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module deleted successfully",
	})
}
