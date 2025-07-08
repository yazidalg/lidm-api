package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
)

// Upgrader untuk mengubah koneksi HTTP menjadi Webcommon.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Izinkan semua koneksi untuk development.
		// Di produksi, Anda harus memvalidasi origin di sini.
		return true
	},
}

// SocketHandler menangani logika terkait Webcommon.
type SocketHandler struct {
	hub *common.Hub
}

func NewSocketHandler(hub *common.Hub) *SocketHandler {
	return &SocketHandler{hub: hub}
}

// ServeWs adalah handler Gin untuk endpoint Webcommon.
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

	// Upgrade koneksi HTTP ke Webcommon.
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for room %s: %v", roomName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Room is Full. cannot join"})
		return
	}

	// Buat client baru.
	client := &common.Client{
		Hub:  h.hub,
		Conn: conn,
		Send: make(chan *common.Message, 256), // Buffered channel
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
	userVal, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	roomName := h.hub.FindAndAssignRoom()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection for matchmaking: %v", err)
		return
	}

	client := &common.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan *common.Message, 256), // Buffered channel
		Room:     roomName,
		UserID:   user.ID,
		Username: user.Name,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *SocketHandler) PreQuiz(c *gin.Context) {
	userVal, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	user, ok := userVal.(models.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection for prequiz: %v", err)
		return
	}

	roomName := fmt.Sprintf("prequiz-%d", user.ID)

	client := &common.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan *common.Message, 256), // Buffered channel
		Room:     roomName,
		UserID:   user.ID,
		Username: user.Name,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
