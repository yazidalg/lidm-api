package handlers

import (
	"fmt"
	"net/http"
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

	verifiedUser, err := h.authService.GetVerifiedUser(findByEmail.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User is not verified",
			"details": err.Error(),
		})
		return
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(verifiedUser.Password), []byte(body.Password))

	if passwordError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password doesn't match",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": verifiedUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("SECRET"))

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenString, 3600, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"data":    tokenString,
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
