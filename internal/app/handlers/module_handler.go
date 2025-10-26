package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

type ModuleHandler struct {
	moduleService services.ModuleServiceInterface
}

func NewModuleHandler(moduleService services.ModuleServiceInterface) *ModuleHandler {
	return &ModuleHandler{moduleService}
}

func (h *ModuleHandler) CreateModule(c *gin.Context) {
	// Check if request is multipart form (for file upload) or JSON
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle multipart form with optional icon upload
		h.createModuleWithIcon(c)
	} else {
		// Handle JSON request (original behavior)
		h.createModuleJSON(c)
	}
}

func (h *ModuleHandler) CreateModuleWithVideo(c *gin.Context) {
	var request request.CreateModuleWithVideoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result, err := h.moduleService.CreateModuleWithVideo(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create module with video material",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Module with video material created successfully",
		"data":    result,
	})
}

func (h *ModuleHandler) createModuleJSON(c *gin.Context) {
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

func (h *ModuleHandler) createModuleWithIcon(c *gin.Context) {
	// Parse form data
	title := c.PostForm("title")
	description := c.PostForm("description")

	if title == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and description are required"})
		return
	}

	// Parse optional fields
	var offsetX, offsetY *float64
	if offsetXStr := c.PostForm("offset_x"); offsetXStr != "" {
		if val, err := strconv.ParseFloat(offsetXStr, 64); err == nil {
			offsetX = &val
		}
	}
	if offsetYStr := c.PostForm("offset_y"); offsetYStr != "" {
		if val, err := strconv.ParseFloat(offsetYStr, 64); err == nil {
			offsetY = &val
		}
	}

	// Create module request
	moduleRequest := request.ModuleRequest{
		Title:       title,
		Description: description,
		OffsetX:     offsetX,
		OffsetY:     offsetY,
	}

	// Create module first
	result, err := h.moduleService.CreateModule(moduleRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create module",
		})
		return
	}

	// Check if icon file is provided
	_, header, err := c.Request.FormFile("icon")
	if err == nil && header != nil {
		// Icon file provided, upload it
		config := utils.DefaultImageUploadConfig()
		filePath, uploadErr := utils.UploadFile(c, "icon", config)
		if uploadErr != nil {
			// Module created but icon upload failed, just log and continue
			c.JSON(http.StatusCreated, gin.H{
				"message":    "Module created successfully, but icon upload failed",
				"data":       result,
				"icon_error": uploadErr.Error(),
			})
			return
		}

		// Update module with icon path
		offsetXFloat := float64(result.OffsetX)
		offsetYFloat := float64(result.OffsetY)
		updateRequest := request.UpdateModuleRequest{
			Title:       result.Title,
			Description: result.Description,
			OffsetX:     &offsetXFloat,
			OffsetY:     &offsetYFloat,
			Icon:        &filePath,
		}

		updatedResult, updateErr := h.moduleService.UpdateModule(uint32(result.ID), updateRequest)
		if updateErr != nil {
			// Clean up uploaded file if update fails
			utils.DeleteFile(filePath)
			c.JSON(http.StatusCreated, gin.H{
				"message":    "Module created successfully, but failed to save icon",
				"data":       result,
				"icon_error": updateErr.Error(),
			})
			return
		}

		// Success with icon
		iconURL := utils.GetFileURL(c, filePath)
		c.JSON(http.StatusCreated, gin.H{
			"message":  "Module created successfully with icon",
			"data":     updatedResult,
			"icon_url": iconURL,
		})
	} else {
		// No icon file provided, return module without icon
		c.JSON(http.StatusCreated, gin.H{
			"message": "Module created successfully",
			"data":    result,
		})
	}
}

func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	// Get user ID from context (set by JWT middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// JWT middleware sets userID as float64, so we need to convert it
	userIDFloat, ok := userIDInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID := uint(userIDFloat)

	// Use the new method with progress information
	result, err := h.moduleService.GetModuleByIDWithProgress(uint32(id), userID)
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
			"message": "Failed to retrieve modules",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved successfully",
		"data":    results,
	})
}

func (h *ModuleHandler) GetAllModulesWithProgress(c *gin.Context) {
	// Extract user ID from auth middleware
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user context"})
		return
	}

	results, err := h.moduleService.GetAllModulesWithProgress(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve modules with progress",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ENHANCED MODULES retrieved successfully", // Changed message to verify
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

	var request request.UpdateModuleRequest
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

// UploadModuleIcon uploads an icon for a module
func (h *ModuleHandler) UploadModuleIcon(c *gin.Context) {
	// Get module ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	// Check if module exists
	module, err := h.moduleService.GetModuleByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Module not found",
			"message": "Module with the given ID does not exist",
		})
		return
	}

	// Configure file upload
	config := utils.DefaultImageUploadConfig()

	// Upload the file
	filePath, err := utils.UploadFile(c, "icon", config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload icon",
		})
		return
	}

	// Delete old icon if exists
	if module.Icon != "" {
		utils.DeleteFile(module.Icon)
	}

	// Update module with new icon path
	offsetXFloat := float64(module.OffsetX)
	offsetYFloat := float64(module.OffsetY)
	updateRequest := request.UpdateModuleRequest{
		Title:       module.Title,
		Description: module.Description,
		OffsetX:     &offsetXFloat,
		OffsetY:     &offsetYFloat,
		Icon:        &filePath,
	}

	updatedModule, err := h.moduleService.UpdateModule(uint32(id), updateRequest)
	if err != nil {
		// If update fails, delete uploaded file
		utils.DeleteFile(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update module with icon",
		})
		return
	}

	// Generate full URL for the icon
	iconURL := utils.GetFileURL(c, filePath)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Icon uploaded successfully",
		"data":     updatedModule,
		"icon_url": iconURL,
	})
}

// DeleteModuleIcon deletes the icon of a module
func (h *ModuleHandler) DeleteModuleIcon(c *gin.Context) {
	// Get module ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	// Check if module exists
	module, err := h.moduleService.GetModuleByID(uint32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Module not found",
			"message": "Module with the given ID does not exist",
		})
		return
	}

	// Check if module has an icon
	if module.Icon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No icon to delete",
			"message": "Module does not have an icon",
		})
		return
	}

	// Delete the icon file
	utils.DeleteFile(module.Icon)

	// Update module to remove icon reference
	offsetXFloat := float64(module.OffsetX)
	offsetYFloat := float64(module.OffsetY)
	updateRequest := request.UpdateModuleRequest{
		Title:       module.Title,
		Description: module.Description,
		OffsetX:     &offsetXFloat,
		OffsetY:     &offsetYFloat,
		Icon:        nil,
	}

	updatedModule, err := h.moduleService.UpdateModule(uint32(id), updateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update module",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Icon deleted successfully",
		"data":    updatedModule,
	})
}

// GetAllModulesWithUnlockStatus returns modules with unlock status and progress for the authenticated user
func (h *ModuleHandler) GetAllModulesWithUnlockStatus(c *gin.Context) {
	// Extract user ID from auth middleware
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	case uint:
		userID = v
	case int:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid user context"})
		return
	}

	results, err := h.moduleService.GetAllModulesWithUnlockStatus(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve modules with unlock status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modules retrieved successfully with unlock status",
		"data":    results,
	})
}

func (h *ModuleHandler) UpdateModuleWithVideo(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var request request.UpdateModuleWithVideoRequest
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

	result, err := h.moduleService.UpdateModuleWithVideo(uint32(moduleId.ID), request)

	// Bentuk respons sesuai format yang Anda minta
	response := response.CustomModuleResponse{
		ID:            uint32(result.ID),
		Title:         result.Title,
		Description:   result.Description,
		Thumbnail:     result.Thumbnail,
		Icon:          result.Icon,
		OffsetX:       int(result.OffsetX),
		OffsetY:       int(result.OffsetY),
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
		VideoMaterial: *result.VideoMaterial,
		ARExperiment:  *result.ARExperiment,
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to update module with video",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module with video updated successfully",
		"data":    response,
	})
}

func (h *ModuleHandler) AddARExperimentToModule(c *gin.Context) {
	var request request.AddARExperimentToModuleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if user is admin
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found in context",
			"message": "Authentication required",
		})
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Invalid user data",
			"message": "Failed to process user information",
		})
		return
	}

	if !userModel.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Admin access required",
			"message": "Only administrators can add AR experiments to modules",
		})
		return
	}

	// Check if module exists
	_, err := h.moduleService.GetModuleByID(request.ModuleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error(),
			"message": "Module not found",
		})
		return
	}

	result, err := h.moduleService.AddARExperimentToModule(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to add AR experiment to module",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "AR experiment added to module successfully",
		"data":    result,
	})
}

// GetAllModulesAdmin returns all modules for admin use (no user-specific data) with pagination
func (h *ModuleHandler) GetAllModulesAdmin(c *gin.Context) {
	// Parse pagination parameters
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Set maximum limit to prevent abuse
	if limit > 100 {
		limit = 100
	}

	// Calculate offset
	offset := (page - 1) * limit

	results, totalCount, err := h.moduleService.GetAllModulesWithPagination(limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve modules",
			"error":   err.Error(),
		})
		return
	}

	// Calculate pagination metadata
	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	hasNext := page < int(totalPages)
	hasPrev := page > 1

	c.JSON(http.StatusOK, gin.H{
		"message": "All modules retrieved successfully",
		"data":    results,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_count":  totalCount,
			"total_pages":  totalPages,
			"has_next":     hasNext,
			"has_previous": hasPrev,
		},
	})
}
