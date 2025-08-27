package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

type FileUploadHandler struct{}

func NewFileUploadHandler() *FileUploadHandler {
	return &FileUploadHandler{}
}

// UploadImage uploads a general image file
func (h *FileUploadHandler) UploadImage(c *gin.Context) {
	// Get upload type from query parameter (default: images)
	uploadType := c.DefaultQuery("type", "images")
	
	var config utils.FileUploadConfig
	
	switch uploadType {
	case "icons":
		config = utils.DefaultImageUploadConfig()
	case "avatars":
		config = utils.FileUploadConfig{
			UploadDir:      "./uploads/avatars",
			MaxFileSize:    3 * 1024 * 1024, // 3MB
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".webp"},
			GenerateUnique: true,
		}
	case "banners":
		config = utils.FileUploadConfig{
			UploadDir:      "./uploads/banners",
			MaxFileSize:    10 * 1024 * 1024, // 10MB
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".webp"},
			GenerateUnique: true,
		}
	default:
		config = utils.LargeImageUploadConfig()
	}

	// Upload the file
	filePath, err := utils.UploadFile(c, "image", config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload image",
		})
		return
	}

	// Return the file path for frontend use
	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"data": gin.H{
			"filename": filepath.Base(filePath),
			"path":     filePath,
			"url":      "/" + filePath, // Frontend can use this directly
			"type":     uploadType,
		},
	})
}

// UploadMultipleImages uploads multiple images at once
func (h *FileUploadHandler) UploadMultipleImages(c *gin.Context) {
	uploadType := c.DefaultQuery("type", "images")
	
	var config utils.FileUploadConfig
	switch uploadType {
	case "icons":
		config = utils.DefaultImageUploadConfig()
	default:
		config = utils.LargeImageUploadConfig()
	}

	// Upload multiple files
	filePaths, err := utils.UploadMultipleFiles(c, "images", config)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload images",
		})
		return
	}

	// Format response
	var uploadedFiles []gin.H
	for _, filePath := range filePaths {
		uploadedFiles = append(uploadedFiles, gin.H{
			"filename": filepath.Base(filePath),
			"path":     filePath,
			"url":      "/" + filePath,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Images uploaded successfully",
		"data": gin.H{
			"files": uploadedFiles,
			"count": len(uploadedFiles),
			"type":  uploadType,
		},
	})
}
