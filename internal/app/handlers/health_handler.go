package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

func NewHealthHandler(db *gorm.DB, startTime time.Time) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: startTime,
	}
}

// Health - Basic health check endpoint
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(h.startTime).String(),
	})
}

// Ready - Readiness check endpoint (comprehensive)
func (h *HealthHandler) Ready(c *gin.Context) {
	status := gin.H{
		"status":      "ready",
		"timestamp":   time.Now().Unix(),
		"uptime":      time.Since(h.startTime).String(),
		"environment": os.Getenv("ENV"),
	}

	// Test database connection
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			status["database"] = gin.H{
				"status":  "error",
				"message": "Failed to get database instance",
				"error":   err.Error(),
			}
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}

		// Ping database
		if err := sqlDB.Ping(); err != nil {
			status["database"] = gin.H{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			}
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}

		// Get database stats
		stats := sqlDB.Stats()
		status["database"] = gin.H{
			"status":  "connected",
			"message": "Database connection successful",
			"stats": gin.H{
				"open_connections": stats.OpenConnections,
				"in_use":           stats.InUse,
				"idle":             stats.Idle,
				"wait_count":       stats.WaitCount,
				"wait_duration":    stats.WaitDuration.String(),
			},
		}
	} else {
		status["database"] = gin.H{
			"status":  "error",
			"message": "Database instance is nil",
		}
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	// Check environment variables
	envVars := gin.H{
		"ENV":  os.Getenv("ENV"),
		"PORT": os.Getenv("PORT"),
		"DB_USER": func() string {
			if os.Getenv("DB_USER") != "" {
				return "***set***"
			}
			return "not_set"
		}(),
		"DB_NAME": func() string {
			if os.Getenv("DB_NAME") != "" {
				return "***set***"
			}
			return "not_set"
		}(),
		"INSTANCE_CONNECTION_NAME": func() string {
			if os.Getenv("INSTANCE_CONNECTION_NAME") != "" {
				return "***set***"
			}
			return "not_set"
		}(),
	}
	status["environment"] = envVars

	c.JSON(http.StatusOK, status)
}

// Healthy - Liveness check endpoint (simple)
func (h *HealthHandler) Healthy(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(h.startTime).String(),
	}

	// Simple database ping
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			status["database"] = "error"
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			status["database"] = "disconnected"
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}

		status["database"] = "connected"
	} else {
		status["database"] = "not_initialized"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
