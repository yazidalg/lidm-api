package request

type UserRegisterRequest struct {
	Name     string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Class    string `json:"class" binding:"required"`
}
