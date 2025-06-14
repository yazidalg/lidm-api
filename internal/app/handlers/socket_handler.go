package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yazidalg/lidm_backend/internal/app/socket"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

// Upgrader untuk mengubah koneksi HTTP menjadi WebSocket.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Izinkan semua koneksi untuk development.
		// Di produksi, Anda harus memvalidasi origin di sini.
		return true
	},
}

// SocketHandler menangani logika terkait WebSocket.
type SocketHandler struct {
	hub *socket.Hub
}

func NewSocketHandler(hub *socket.Hub) *SocketHandler {
	return &SocketHandler{hub: hub}
}

// ServeWs adalah handler Gin untuk endpoint WebSocket.
// Ia akan meng-upgrade koneksi dan mendaftarkan client ke Hub.
func (h *SocketHandler) ServeWs(c *gin.Context) {
	// Ekstrak nama room dari URL.
	roomName := c.Param("room")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room name is required"})
		return
	}

	// Mendapatkan jumlah client yang sudah ada di room.
	clientCount := h.hub.GetClientsInRoom(roomName)

	// Jika jumlah client sudah mencapai batas (2), tolak permintaan.
	if clientCount >= 2 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "room is full"})
		return
	}

	// Upgrade koneksi HTTP ke WebSocket.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for room %s: %v", roomName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room is Full. cannot join"})
		return
	}

	// Buat client baru.
	client := &socket.Client{
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan *utils.Message, 256), // Buffered channel
		Room: roomName,
	}

	// Daftarkan client ke hub.
	client.Hub.Register <- client

	// Jalankan readPump dan writePump sebagai goroutine.
	// Ini memungkinkan client untuk membaca dan menulis secara bersamaan.
	go client.WritePump()
	go client.ReadPump()
}

func (h *SocketHandler) MatchMaking(c *gin.Context) {
	roomName := h.hub.FindAndAssignRoom()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection for matchmaking: %v", err)
		return
	}

	client := &socket.Client{
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan *utils.Message, 256), // Buffered channel
		Room: roomName,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
