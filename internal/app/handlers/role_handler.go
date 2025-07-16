package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type RoleHandler struct {
	roleService services.RoleServiceInterface
}

func NewRoleHandler(roleService services.RoleServiceInterface) *RoleHandler {
	return &RoleHandler{roleService}
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve roles",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Roles retrieved successfully",
		"data":    roles,
	})
}

func (h *RoleHandler) GetRoleById(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Role ID is required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(roleID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role ID",
		})
		return
	}

	role, err := h.roleService.GetRoleById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Role not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role retrieved successfully",
		"data":    role,
	})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var body struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	role := models.Role{
		Name:        body.Name,
		Description: body.Description,
	}

	createdRole, err := h.roleService.CreateRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"data":    createdRole,
	})
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Role ID is required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(roleID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role ID",
		})
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	// Get existing role first
	existingRole, err := h.roleService.GetRoleById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Role not found",
			"error":   err.Error(),
		})
		return
	}

	// Update fields
	if body.Name != "" {
		existingRole.Name = body.Name
	}
	if body.Description != "" {
		existingRole.Description = body.Description
	}

	updatedRole, err := h.roleService.UpdateRole(existingRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
		"data":    updatedRole,
	})
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("id")
	if roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Role ID is required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(roleID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role ID",
		})
		return
	}

	err := h.roleService.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}
