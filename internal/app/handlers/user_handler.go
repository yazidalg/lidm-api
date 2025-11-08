package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type UserHandler struct {
	userService        services.UserServiceInterface
	leaderboardService *services.LeaderboardService
}

func NewUserHandler(userService services.UserServiceInterface, leaderboardService *services.LeaderboardService) *UserHandler {
	return &UserHandler{
		userService:        userService,
		leaderboardService: leaderboardService,
	}
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	userVal, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not unauthorized"})
		return
	}

	user := userVal.(models.User)

	// Update user position before returning profile
	if h.leaderboardService != nil {
		_ = h.leaderboardService.UpdateUserPosition(user.ID)

		// Refresh user data to get updated position information
		if updatedUser, err := h.userService.GetUserByIDUint(user.ID); err == nil && updatedUser != nil {
			user = *updatedUser
		}
	}

	// Add the two fields you requested to the user response
	responseData := user

	// Add position information in a simple format
	response := gin.H{
		"message":         "User found",
		"data":            responseData,
		"type":            user.PositionType,   // "increasing" or "decreasing"
		"position_change": user.PositionChange, // 1, 2, 3, etc.
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserAdmin(c *gin.Context) {
	userVal, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not unauthorized"})
		return
	}
	user := userVal.(models.User)

	responseData := response.UserAdminResponse{
		ID:                uint32(user.ID),
		Email:             user.Email,
		Name:              user.Name,
		IsVerified:        user.IsVerified,
		VerificationToken: user.VerificationToken,
		RoleID:            user.RoleID,
		RoleName:          user.Role.Name,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User admin found",
		"data":    responseData,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	var body struct {
		Role string `json:"role"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Validate role
	if body.Role != "user" && body.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role. Must be 'user' or 'admin'",
		})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	// Convert string to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	updatedUser, err := h.userService.UpdateUserRole(id, body.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user role",
			"error":   err.Error(),
		})
		return
}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"data":    updatedUser,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	// Convert string to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) UpdateAccount(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	// Convert string to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}

	var req request.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	updatedUser, err := h.userService.UpdateAccount(id, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update account",
			"error":   err.Error(),
		})
		return
	}

	responseData := response.UpdateAccountResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"data":    responseData,
	})
}
