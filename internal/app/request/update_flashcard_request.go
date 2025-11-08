package request

type CreateFlashcardRequest struct {
	ModuleID  uint   `json:"module_id" binding:"required"`
	FrontText string `json:"front_text" binding:"required"`
	BackText  string `json:"back_text" binding:"required"`
	Order     int    `json:"order" binding:"omitempty"`
}

type UpdateFlashcardRequest struct {
	ModuleID  *uint   `json:"module_id" binding:"omitempty"`
	FrontText *string `json:"front_text" binding:"omitempty"`
	BackText  *string `json:"back_text" binding:"omitempty"`
	Order     *int    `json:"order" binding:"omitempty"`
}

