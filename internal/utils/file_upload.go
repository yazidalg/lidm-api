package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// FileUploadConfig contains configuration for file uploads
type FileUploadConfig struct {
	UploadDir      string
	MaxFileSize    int64
	AllowedTypes   []string
	GenerateUnique bool
}

// DefaultImageUploadConfig returns default config for image uploads
func DefaultImageUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		UploadDir:      "./uploads/icons",
		MaxFileSize:    5 * 1024 * 1024, // 5MB
		AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".svg", ".gif", ".webp"},
		GenerateUnique: true,
	}
}

// LargeImageUploadConfig returns config for larger images (thumbnails, banners)
func LargeImageUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		UploadDir:      "./uploads/images",
		MaxFileSize:    10 * 1024 * 1024, // 10MB
		AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".webp"},
		GenerateUnique: true,
	}
}

// UploadFile uploads a file and returns the file path
func UploadFile(c *gin.Context, fieldName string, config FileUploadConfig) (string, error) {
	// Get the file from form
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("failed to get file: %v", err)
	}
	defer file.Close()

	// Check file size
	if header.Size > config.MaxFileSize {
		return "", fmt.Errorf("file size exceeds limit of %d bytes", config.MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedFileType(ext, config.AllowedTypes) {
		return "", fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, config.AllowedTypes)
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate filename
	filename := header.Filename
	if config.GenerateUnique {
		filename = generateUniqueFilename(header.Filename)
	}

	// Full path
	fullPath := filepath.Join(config.UploadDir, filename)

	// Create the file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the created file
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return relative path for database storage
	relativePath := filepath.Join("uploads/icons", filename)
	return relativePath, nil
}

// UploadMultipleFiles uploads multiple files
func UploadMultipleFiles(c *gin.Context, fieldName string, config FileUploadConfig) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form: %v", err)
	}

	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in field %s", fieldName)
	}

	var uploadedFiles []string

	for _, fileHeader := range files {
		// Check file size
		if fileHeader.Size > config.MaxFileSize {
			return nil, fmt.Errorf("file %s exceeds size limit", fileHeader.Filename)
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !isAllowedFileType(ext, config.AllowedTypes) {
			return nil, fmt.Errorf("file type %s not allowed for %s", ext, fileHeader.Filename)
		}

		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		// Create upload directory if not exists
		if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %v", err)
		}

		// Generate filename
		filename := fileHeader.Filename
		if config.GenerateUnique {
			filename = generateUniqueFilename(fileHeader.Filename)
		}

		// Full path
		fullPath := filepath.Join(config.UploadDir, filename)

		// Create the file
		dst, err := os.Create(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create file %s: %v", filename, err)
		}
		defer dst.Close()

		// Copy the uploaded file to the created file
		if _, err := io.Copy(dst, file); err != nil {
			return nil, fmt.Errorf("failed to save file %s: %v", filename, err)
		}

		// Add to uploaded files list
		relativePath := filepath.Join("uploads/icons", filename)
		uploadedFiles = append(uploadedFiles, relativePath)
	}

	return uploadedFiles, nil
}

// DeleteFile deletes a file from the filesystem
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil // Nothing to delete
	}

	// Convert relative path to absolute path
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(".", filePath)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	// Delete the file
	return os.Remove(filePath)
}

// isAllowedFileType checks if the file extension is allowed
func isAllowedFileType(ext string, allowedTypes []string) bool {
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// generateUniqueFilename generates a unique filename with timestamp
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// GetFileURL returns the public URL for a file
func GetFileURL(c *gin.Context, filePath string) string {
	if filePath == "" {
		return ""
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)
	return fmt.Sprintf("%s/%s", baseURL, filePath)
}
