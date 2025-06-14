package utils

import "encoding/json"

// Structure Json Definition
type Message struct {
	Action  string          `json:"action"`            // Contoh: "send_message", "join_room"
	Message json.RawMessage `json:"payload,omitempty"` // Isi pesan teks
	Target  string          `json:"target,omitempty"`  // Nama room yang dituju
	Sender  string          `json:"sender"`            // (Opsional) Nama pengirim
}

// Payload for User Join and User Leave events
type UserEventPayload struct {
	UserID   uint   `json:"userID"`   // Nama pengguna yang bergabung atau meninggalkan
	Username string `json:"username"` // Nama pengguna yang bergabung atau meninggalkan
	Room     string `json:"room"`     // Nama room yang dimasuki atau ditinggalkan
	Message  string `json:"message"`  // Pesan yang ditampilkan saat bergabung atau meninggalkan
}
