package common

import (
	"encoding/json"
)

// Structure Json Definition
type Message struct {
	Action  string          `json:"action"`            // Contoh: "send_message", "join_room"
	Payload json.RawMessage `json:"payload,omitempty"` // Isi pesan teks
	Target  string          `json:"target,omitempty"`  // Nama room yang dituju
	Sender  *Client         `json:"-"`                 // (Opsional) Nama pengirim (gunakan string untuk menghindari import cycle)
}

// Payload for User Join and User Leave events
type UserEventPayload struct {
	UserID   uint   `json:"userID"`   // Nama pengguna yang bergabung atau meninggalkan
	Username string `json:"username"` // Nama pengguna yang bergabung atau meninggalkan
	Room     string `json:"room"`     // Nama room yang dimasuki atau ditinggalkan
	Message  string `json:"message"`  // Pesan yang ditampilkan saat bergabung atau meninggalkan
}

// Payload for both of users are in same room
type UserInRoomPayload struct {
	FirstUsername  string `json:"firstUsername"`  // Nama pengguna pertama
	SecondUsername string `json:"secondUsername"` // Nama pengguna kedua
	Room           string `json:"room"`           // Nama room tempat kedua pengguna berada
	Message        string `json:"message"`        // Pesan yang ditampilkan saat kedua pengguna berada di room yang sama
}

type AnsweredQuestionEvent struct {
	QuestionID uint `json:"question_id"` // ID pertanyaan yang dijawab
	UserID     uint `json:"user_id"`     // ID pengguna yang menjawab
	IsCorrect  bool `json:"is_correct"`  // Apakah jawaban benar
	Score      int  `json:"score"`       // Skor yang diberikan untuk jawaban
}

type QuizCompletedPayload struct {
	FinalScores map[string]int `json:"final_scores"` // Skor akhir untuk setiap pemain
	Winner      string         `json:"winner"`       // Nama pemenang
	Message     string         `json:"message"`      // Pesan yang ditampilkan saat quiz selesai
}

type AnswerPayload struct {
	QuestionID     uint   `json:"question_id"`
	UserID         uint   `json:"user_id"`
	OptionSelected string `json:"option_selected"`
}

type PrequizCompletedPayload struct {
	FinalScores map[string]uint `json:"final_result"` // Skor akhir untuk setiap pemain
	Message     string          `json:"message"`      // Pesan yang ditampilkan saat prequiz selesai
}

// AnswerEvent adalah event yang dikirim ketika pemain menjawab pertanyaan
type AnswerEvent struct {
	Player  *Client
	Payload AnswerPayload
}

type AnsweredPreQuizEvent struct {
	ID        uint `json:"question_id"` // ID pertanyaan yang dijawab
	IsCorrect bool `json:"is_correct"`  // Apakah jawaban benar
}
