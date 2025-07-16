package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authServiceInterface services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authServiceInterface}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
		Class    string
		Role     string `json:"role,omitempty"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
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

	user := request.UserRegisterRequest{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hash),
		Class:    body.Class,
		Role:     body.Role,
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": findByEmail.ID,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate token",
			"details": err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenStr, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    tokenStr,
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
