package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userService services.UserServiceInterface) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	userVal, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not unauthorized"})
		return
	}

	user := userVal.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"message": "User found",
		"data":    user,
	})
}
