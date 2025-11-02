package request

type UpdateFlashcardRequest struct {
	ModuleID  *uint   `json:"module_id" binding:"omitempty"`
	FrontText *string `json:"front_text" binding:"omitempty"`
	BackText  *string `json:"back_text" binding:"omitempty"`
	Order     *int    `json:"order" binding:"omitempty"`
}

