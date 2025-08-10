package database

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) error {
	// Check if admin already exists
	var adminUser models.User
	result := db.Where("email = ?", "admin@lidm.com").First(&adminUser)

	if result.Error == nil {
		// Admin user already exists
		return nil
	}

	// Get admin role
	var adminRole models.Role
	if err := db.Where("name = ?", models.RoleAdminName).First(&adminRole).Error; err != nil {
		return err
	}

	// Hash default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create admin user
	admin := models.User{
		Name:       "Administrator",
		Email:      "admin@lidm.com",
		Password:   string(hashedPassword),
		Class:      "", // Admin doesn't need class
		RoleID:     adminRole.ID,
		IsVerified: true, // Auto verify admin
		Point:      0,
		TotalXP:    0,
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}
