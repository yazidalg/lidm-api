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

type UpdateAccountResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
