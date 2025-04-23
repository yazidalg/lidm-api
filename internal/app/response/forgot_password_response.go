package response

type ForgotPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
