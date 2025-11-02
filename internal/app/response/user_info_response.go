package response

type UserAdminResponse struct {
	ID                uint32 `json:"ID"`
	Email             string `json:"Email"`
	Name              string `json:"Name"`
	IsVerified        bool   `json:"IsVerified"`
	VerificationToken string `json:"VerificationToken"`
	RoleID            uint   `json:"RoleID"`
	RoleName          string `json:"RoleName"`
}
