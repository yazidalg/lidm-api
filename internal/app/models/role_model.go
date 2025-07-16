package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Users []User `gorm:"foreignKey:RoleID" json:"users,omitempty"`
}

// Role constants
const (
	RoleUserName  = "user"
	RoleAdminName = "admin"
)

// Helper methods
func (r *Role) IsAdmin() bool {
	return r.Name == RoleAdminName
}

func (r *Role) IsUser() bool {
	return r.Name == RoleUserName
}

// Seed data function
func SeedRoles(db *gorm.DB) error {
	roles := []Role{
		{
			Name:        RoleUserName,
			Description: "Regular user with basic permissions",
		},
		{
			Name:        RoleAdminName,
			Description: "Administrator with full permissions",
		},
	}

	for _, role := range roles {
		var existingRole Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.Error != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
