package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type AuthMiddleware struct {
	authService services.AuthServiceInterface
}

func NewAuthMiddleware(authService services.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{authService}
}

// RequireAuth - middleware untuk memastikan user sudah login
func (m *AuthMiddleware) RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("authorization")
	if err != nil {
		// Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization token required",
			})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenString = authHeader
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token expired",
			})
			c.Abort()
			return
		}

		userIdFloat, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}

		user, err := m.authService.LoginUser(int(userIdFloat))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		c.Abort()
		return
	}
}

// RequireAdmin - middleware untuk memastikan user adalah admin
func (m *AuthMiddleware) RequireAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not found in context",
		})
		c.Abort()
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user data",
		})
		c.Abort()
		return
	}

	if !userModel.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Admin access required",
		})
		c.Abort()
		return
	}

	c.Next()
}

// RequireUserOrAdmin - middleware untuk memastikan user adalah owner resource atau admin
func (m *AuthMiddleware) RequireUserOrAdmin(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User not found in context",
		})
		c.Abort()
		return
	}

	userModel, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Invalid user data",
		})
		c.Abort()
		return
	}

	// Jika admin, langsung izinkan
	if userModel.IsAdmin() {
		c.Next()
		return
	}

	// Jika bukan admin, cek apakah user adalah pemilik resource
	requestedUserID := c.Param("user_id")
	if requestedUserID == "" {
		requestedUserID = c.Param("id")
	}

	if requestedUserID != "" && requestedUserID != fmt.Sprint(userModel.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Access denied. You can only access your own data",
		})
		c.Abort()
		return
	}

	c.Next()
}
