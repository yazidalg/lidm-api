package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🔐 LIDM Admin Creator")
	fmt.Println("====================")

	// Load environment
	config.LoadEnv()
	db := config.ConnectDB()

	reader := bufio.NewReader(os.Stdin)

	// Get admin details
	fmt.Print("Nama Admin: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Email Admin: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Password Admin: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Validate input
	if name == "" || email == "" || password == "" {
		fmt.Println("❌ Semua field harus diisi!")
		return
	}

	// Check if email already exists
	var existingUser models.User
	result := db.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		fmt.Printf("❌ Email %s sudah digunakan!\n", email)
		return
	}

	// Get admin role
	var adminRole models.Role
	if err := db.Where("name = ?", models.RoleAdminName).First(&adminRole).Error; err != nil {
		fmt.Printf("❌ Role admin tidak ditemukan: %v\n", err)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("❌ Error hashing password: %v\n", err)
		return
	}

	// Create admin user
	admin := models.User{
		Name:       name,
		Email:      email,
		Password:   string(hashedPassword),
		Class:      "", // Admin doesn't need class
		RoleID:     adminRole.ID,
		IsVerified: true, // Auto verify admin
		Point:      0,
		TotalXP:    0,
	}

	if err := db.Create(&admin).Error; err != nil {
		fmt.Printf("❌ Error membuat admin: %v\n", err)
		return
	}

	fmt.Println("\n✅ Admin berhasil dibuat!")
	fmt.Printf("📧 Email: %s\n", email)
	fmt.Printf("👤 Nama: %s\n", name)
	fmt.Println("\n🚀 Silakan login menggunakan credentials di atas!")
}
