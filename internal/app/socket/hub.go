package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/yazidalg/lidm_backend/internal/utils"
)

// Hub mengelola semua room, client, dan meneruskan pesan.
type Hub struct {
	// Kumpulan client yang terdaftar.
	Clients map[*Client]bool

	// Pesan masuk dari client yang akan di-broadcast ke room yang tepat.
	Message chan *utils.Message

	// Channel untuk mendaftarkan client baru.
	Register chan *Client

	// Channel untuk membatalkan pendaftaran client.
	Unregister chan *Client

	// Rooms mengelompokkan client. Key adalah nama room,
	// value adalah sebuah set dari client di dalam room tersebut.
	Rooms map[string]map[*Client]bool

	// Mu adalah mutex untuk mengamankan akses ke Rooms dan Clients.
	Mu sync.RWMutex

	// Counter Room Number
	RoomNumber int

	UserID uint

	Username string
}

// NewHub membuat instance Hub baru.
func NewHub() *Hub {
	return &Hub{
		Message:    make(chan *utils.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		// Kita menggunakan map[*Client]bool sebagai 'set' untuk efisiensi.
		Rooms:      make(map[string]map[*Client]bool),
		RoomNumber: 1,
	}
}

// Run menjalankan Hub. Ini harus dijalankan sebagai goroutine.
func (h *Hub) Run() {
	for {
		// Menunggu salah satu dari empat channel menerima data.
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			roomName := client.Room
			if _, ok := h.Rooms[roomName]; !ok {
				h.Rooms[roomName] = make(map[*Client]bool)
			}
			h.Rooms[roomName][client] = true
			log.Printf("Client terdaftar di room '%s'. Total di room: %d", roomName, len(h.Rooms[roomName]))

			user_join_payload := &utils.UserEventPayload{
				UserID:   client.UserID,
				Username: client.Username,
				Room:     roomName,
				Message:  fmt.Sprintf("%s telah bergabung ke room %s", client.Username, roomName),
			}

			userJoinPayloadBytes, _ := json.Marshal(user_join_payload)

			joinMessage := &utils.Message{
				Action:  "user_join",
				Message: userJoinPayloadBytes,
				Target:  roomName,
				Sender:  client.Username,
			}

			h.BroadcastToRoom(joinMessage)
			h.Mu.Unlock()

		case client := <-h.Unregister:
			h.Mu.Lock()
			roomName := client.Room
			if clients, ok := h.Rooms[roomName]; ok {
				if _, ok := clients[client]; ok {
					// Hapus client dari room
					delete(clients, client)
					// Tutup channel send untuk menghentikan writePump goroutine client tersebut
					close(client.Send)

					user_leave_payload := &utils.UserEventPayload{
						UserID:   client.UserID,
						Username: client.Username,
						Room:     roomName,
						Message:  fmt.Sprintf("%s telah meninggalkan room %s", client.Username, roomName),
					}

					userLeavePayloadBytes, _ := json.Marshal(user_leave_payload)

					leaveMessage := &utils.Message{
						Action:  "user_leave",
						Message: userLeavePayloadBytes,
						Target:  roomName,
						Sender:  client.Username,
					}

					h.BroadcastToRoom(leaveMessage)

					// Jika room menjadi kosong, hapus room dari memori
					if len(clients) == 0 {
						delete(h.Rooms, roomName)
						log.Printf("Room '%s' kosong dan dihapus.", roomName)
					}
				}
			}
			h.Mu.Unlock()

		case msg := <-h.Message:
			h.Mu.RLock() // Gunakan Read Lock untuk broadcast
			roomName := msg.Target
			if clients, ok := h.Rooms[roomName]; ok {
				// Broadcast pesan ke semua client di room
				for client := range clients {
					select {
					case client.Send <- msg:
					default:
						// Jika channel send client penuh (mungkin client lambat/hang),
						// kita tutup koneksinya untuk mencegah Hub terblokir.
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.Mu.RUnlock()
		}
	}
}

func (h *Hub) RegisterClient(client *Client) {
	// Daftarkan client ke room yang dia tuju.
	roomName := client.Room
	if _, ok := h.Rooms[roomName]; !ok {
		// Jika room belum ada, buat room baru.
		h.Rooms[roomName] = make(map[*Client]bool)
	}
	// Tambahkan client ke room.
	h.Rooms[roomName][client] = true
	log.Printf("Client terdaftar di room '%s'", roomName)
}

func (h *Hub) UnregisterClient(client *Client) {
	// Hapus client dari room.
	roomName := client.Room
	if clients, ok := h.Rooms[roomName]; ok {
		// Hapus client dari map.
		delete(clients, client)
		// Jika room menjadi kosong setelah client keluar, hapus room tersebut.
		if len(clients) == 0 {
			delete(h.Rooms, roomName)
			log.Printf("Room '%s' kosong dan dihapus", roomName)
		}
	}
}

func (h *Hub) BroadcastToRoom(msg *utils.Message) {
	// Cari room yang dituju oleh pesan.
	roomName := msg.Target
	if clients, ok := h.Rooms[roomName]; ok {
		log.Printf("Broadcast pesan ke room '%s': %s", roomName, msg.Message)
		// Kirim pesan ke setiap client di dalam room.
		for client := range clients {
			// Kirim pesan ke channel 'send' milik client.
			// Client akan menanganinya secara asynchronous.
			client.Send <- msg
		}
	}
}

func (h *Hub) GetClientsInRoom(roomName string) int {
	h.Mu.RLock() // Gunakan Read Lock karena kita hanya membaca data

	defer h.Mu.RUnlock()

	// Kode di dalam RLock dan RUnlock dijamin aman dari race conditions.
	return len(h.Rooms[roomName])
}

func (h *Hub) FindAndAssignRoom() string {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	fmt.Printf("Jumlah rooms: %d\n", len(h.Rooms))

	for name, clients := range h.Rooms {
		if len(h.Rooms) == 0 {
			log.Println("Tidak ada room yang tersedia, membuat room baru")
			roomName := fmt.Sprintf("Room-%d", h.RoomNumber)
			h.Rooms[roomName] = make(map[*Client]bool)
			h.RoomNumber++
			log.Printf("Room baru '%s' dibuat", roomName)
		} else if len(clients) == 1 {
			log.Printf("Room '%s' ditemukan dengan %d client, akan digunakan. menambahkan 1 pemain", name, len(clients))
			return name
		} else {
			log.Printf("Room '%s' penuh dengan %d client, mencari room lain", name, len(clients))
		}
	}

	log.Println("Semua room penuh, membuat room baru")
	roomName := fmt.Sprintf("Room-%d", h.RoomNumber)
	h.Rooms[roomName] = make(map[*Client]bool)
	h.RoomNumber++

	log.Printf("Room baru '%s' dibuat", roomName)
	return roomName
}
