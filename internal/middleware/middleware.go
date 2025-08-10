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
	userService services.UserServiceInterface
}

func NewAuthMiddleware(authService services.AuthServiceInterface, userService services.UserServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{authService: authService, userService: userService}
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

	// Get JWT secret with fallback
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SECRET")
	}

	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "JWT secret not configured",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
			"details": err.Error(),
		})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check token expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token expired",
			})
			c.Abort()
			return
		}

		// Extract the user ID from token claims
		userIdFloat, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims - user ID not found",
			})
			c.Abort()
			return
		}

		// Extract email from token claims
		email, emailExists := claims["email"].(string)
		if !emailExists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims - email not found",
			})
			c.Abort()
			return
		}

		// Extract role_id from token claims
		roleID, roleExists := claims["role_id"].(float64)
		if !roleExists {
			roleID = 1 // Default to user role if not specified
		}

		// Fetch complete user data from database using email from token
		user, err := m.authService.GetByEmail(email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not found",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Reset lives harian jika perlu (max 5) dengan persistence
		_ = m.userService.ResetLivesIfNeeded(user.ID, 5)
		// Refresh user setelah kemungkinan reset
		if refreshed, err := m.userService.GetUserByIDUint(user.ID); err == nil && refreshed != nil {
			user = *refreshed
		}

		// Set user data in context
		c.Set("user", user)
		c.Set("userID", userIdFloat)
		c.Set("user_id", userIdFloat) // For activity tracking compatibility
		c.Set("email", email)
		c.Set("roleID", roleID)
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token claims",
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
