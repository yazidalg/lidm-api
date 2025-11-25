package request

type UpdateAccountRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PhotoProfile string `json:"photo_profile"`
}

