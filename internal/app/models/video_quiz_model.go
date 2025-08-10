package models

import "gorm.io/gorm"

type VideoQuiz struct {
	gorm.Model
	VideoMaterialID uint    `gorm:"not null;index" json:"video_material_id"` // Foreign key ke VideoMaterial
	Question        string  `gorm:"not null" json:"question"`                // Pertanyaan quiz
	TimestampStart  int     `gorm:"not null" json:"timestamp_start"`         // Detik ke berapa quiz muncul (dalam detik)
	TimestampEnd    int     `gorm:"not null" json:"timestamp_end"`           // Detik ke berapa quiz berakhir
	Options         Options `gorm:"embedded" json:"options"`                 // Pilihan jawaban
	CorrectAnswer   string  `gorm:"not null" json:"correct_answer"`          // Jawaban yang benar
	Explanation     string  `gorm:"type:text" json:"explanation"`            // Penjelasan jawaban
	Order           int     `gorm:"not null;default:0" json:"order"`         // Urutan quiz dalam video

	// Relationships
	VideoMaterial *VideoMaterial `gorm:"foreignKey:VideoMaterialID" json:"video_material,omitempty"`
	UserAnswers   []VideoQuizUserAnswer `gorm:"foreignKey:VideoQuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_answers,omitempty"`
}

// Model untuk menyimpan jawaban user pada video quiz
type VideoQuizUserAnswer struct {
	gorm.Model
	VideoQuizID   uint   `gorm:"not null;index" json:"video_quiz_id"`   // Foreign key ke VideoQuiz
	UserID        uint   `gorm:"not null;index" json:"user_id"`         // Foreign key ke User
	SelectedAnswer string `gorm:"not null" json:"selected_answer"`       // Jawaban yang dipilih user (A/B/C/D)
	IsCorrect     bool   `gorm:"not null" json:"is_correct"`            // Apakah jawaban benar
	AnsweredAt    int64  `gorm:"not null" json:"answered_at"`           // Timestamp kapan dijawab (Unix timestamp)
	ResponseTime  int    `gorm:"not null" json:"response_time"`         // Waktu respons dalam detik

	// Relationships
	VideoQuiz *VideoQuiz `gorm:"foreignKey:VideoQuizID" json:"video_quiz,omitempty"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
