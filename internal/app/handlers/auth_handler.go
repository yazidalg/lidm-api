package handlers

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authServiceInterface services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authServiceInterface}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var body request.UserRegisterRequest

	if err := c.Bind(&body); err != nil {
		log.Printf("Bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Request body: %+v", body)

	// Validasi role
	if body.RoleName != "" && body.RoleName != "user" && body.RoleName != "admin" && body.RoleName != "teacher" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid role. Must be 'user', 'admin', or 'teacher'",
		})
		return
	}

	// Validasi class - hanya wajib untuk user biasa, admin dan teacher boleh kosong
	if body.RoleName != "admin" && body.RoleName != "teacher" && body.Class == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Class is required for user role",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to hash password",
		})
		return
	}

	// Set default role jika tidak ada
	if body.RoleName == "" {
		body.RoleName = "user"
	}

	user := request.UserRegisterRequest{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		Class:    body.Class,
		RoleName: body.RoleName,
	}

	result, err := h.authService.RegisterUser(user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": err.Error(),
			"message": "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"data":    result,
	})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	findByEmail, err := h.authService.GetByEmail(body.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not found",
			"details": err.Error(),
		})
		return
	}

	if findByEmail.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email or Password is wrong",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(findByEmail.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email or Password is wrong",
		})
		return
	}

	// Generate JWT token with consistent structure
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SECRET") // fallback to SECRET if JWT_SECRET not set
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      findByEmail.ID,
		"email":   findByEmail.Email,
		"role_id": findByEmail.RoleID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenStr, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
			"details": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenStr, 3600*24*7, "", "", false, true) // 7 days

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"token":   tokenStr,
		"user":    findByEmail,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("verificationToken")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token is required",
		})
		return
	}

	user, err := h.authService.GetByVerificationToken(token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Token is invalid",
			"details": err.Error(),
		})
		return
	}

	user, err = h.authService.VerifyUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to verify user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User verified successfully",
		"data":    user,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var body struct {
		IdToken string `json:"id_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Initialize Google Auth (you should get CLIENT_ID from environment)
	googleAuth := utils.NewGoogleAuth(os.Getenv("GOOGLE_CLIENT_ID"))

	// Verify Google token
	userInfo, err := googleAuth.VerifyGoogleToken(body.IdToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Google token",
			"details": err.Error(),
		})
		return
	}

	// Check if user exists in database
	user, err := h.authService.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist, create new user
		registerReq := request.UserRegisterRequest{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: "", // No password for Google auth users
			Class:    "", // Will need to be updated later
		}

		user, err = h.authService.RegisterUserWithGoogle(registerReq, userInfo.Picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create user",
				"details": err.Error(),
			})
			return
		}
	}

	// Generate JWT token with consistent structure
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SECRET") // fallback to SECRET if JWT_SECRET not set
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user":    user,
	})
}

func (h *AuthHandler) BelajarLogin(c *gin.Context) {
	var body struct {
		IdToken string `json:"id_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Initialize Google Auth (you should get CLIENT_ID from environment)
	googleAuth := utils.NewGoogleAuth(os.Getenv("GOOGLE_CLIENT_ID"))

	// Verify Google token
	userInfo, err := googleAuth.VerifyGoogleToken(body.IdToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid Google token",
			"details": err.Error(),
		})
		return
	}

	// Check if email domain is belajar.id
	// Validasi email harus jenjang.belajar.id (sd, smp, sma)
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@(sd|smp|sma)\.belajar\.id$`)
	if !re.MatchString(userInfo.Email) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya akun belajar.id dengan jenjang sd/smp/sma yang diperbolehkan",
		})
		return
	}

	// Check if user exists in database
	user, err := h.authService.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist, create new user
		registerReq := request.UserRegisterRequest{
			Name:     userInfo.Name,
			Email:    userInfo.Email,
			Password: "", // No password for Google auth users
			Class:    "", // Will need to be updated later
		}

		user, err = h.authService.RegisterUserWithGoogle(registerReq, userInfo.Picture)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create user",
				"details": err.Error(),
			})
			return
		}
	}

	// Generate JWT token with consistent structure
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("SECRET") // fallback to SECRET if JWT_SECRET not set
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"user":    user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear authorization cookie
	c.SetCookie("authorization", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
