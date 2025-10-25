package request

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Class    string `json:"class" binding:"required"`
	RoleName string `json:"role" binding:"omitempty,oneof=user admin teacher"` // Role name instead of role
}
