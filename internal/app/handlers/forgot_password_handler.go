package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

type ForgotPasswordHandler struct {
	forgotPasswordService services.ForgotPasswordServiceInterface
}

func NewForgotPasswordHandler(forgotPasswordService services.ForgotPasswordServiceInterface) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{forgotPasswordService}
}

// RequestPasswordReset handles the initial password reset request that sends OTP
func (h *ForgotPasswordHandler) RequestPasswordReset(c *gin.Context) {
	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	err := h.forgotPasswordService.RequestPasswordReset(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to process password reset request",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset OTP has been sent to your email",
	})
}

// VerifyOTP validates the OTP code from the user
func (h *ForgotPasswordHandler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	isValid, err := h.forgotPasswordService.VerifyOTP(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to verify OTP",
			"error":   err.Error(),
		})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid or expired OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OTP verified successfully",
	})
}

// ResetPassword resets the user's password with the new password
func (h *ForgotPasswordHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	err := h.forgotPasswordService.ResetPassword(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to reset password",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password has been reset successfully",
	})
}
