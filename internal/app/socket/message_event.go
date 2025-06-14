package socket

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
