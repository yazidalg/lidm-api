package utils

// Structure Json Definition
type Message struct {
	Action  string `json:"action"`  // Contoh: "send_message", "join_room"
	Message string `json:"message"` // Isi pesan teks
	Target  string `json:"target"`  // Nama room yang dituju
	Sender  string `json:"sender"`  // (Opsional) Nama pengirim
}
