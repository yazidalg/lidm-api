package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

const (
	// Waktu untuk menunggu write ke client sebelum timeout.
	writeWait = 10 * time.Second
	// Waktu maksimum untuk membaca pesan berikutnya dari client.
	pongWait = 60 * time.Second
	// Interval untuk mengirim ping ke client. Harus lebih kecil dari pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client adalah perantara antara koneksi WebSocket dan Hub.
type Client struct {
	Hub *Hub

	// Koneksi WebSocket itu sendiri.
	Conn *websocket.Conn

	// Channel untuk mengirim pesan dari Hub ke client. Ini adalah buffered channel.
	Send chan *utils.Message

	// Room tempat client ini berada.
	Room string

	UserID   uint   // ID pengguna yang terhubung
	Username string // Nama pengguna yang terhubung
}

// readPump bertugas membaca pesan dari koneksi WebSocket dan mengirimkannya ke Hub.
func (c *Client) ReadPump() {
	defer func() {
		// Saat readPump berhenti (karena client disConnect),
		// kita unregister client dari Hub dan tutup koneksi.
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	// Set deadline untuk membaca pesan.
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// Set handler untuk pesan Pong dari client.
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Loop untuk terus membaca dari koneksi.
	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break // Keluar dari loop jika ada error (misal, client disConnect).
		}

		var msg utils.Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Gagal unmarshal JSON: %v", err)
			continue
		}

		// Set target dan sender dari pesan.
		msg.Target = c.Room
		msg.Sender = c.Conn.RemoteAddr().String() // Gunakan alamat IP:Port sebagai identifier

		// Kirim pesan yang sudah di-parse ke Hub.
		c.Hub.Message <- &msg
	}
}

// writePump bertugas menulis pesan dari Hub ke koneksi WebSocket.
func (c *Client) WritePump() {
	// Buat ticker untuk mengirim ping secara periodik.
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// Set deadline untuk menulis pesan.
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub menutup channel ini.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Encode pesan ke JSON dan kirim.
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Gagal menulis JSON: %v", err)
				return
			}

		case <-ticker.C:
			// Kirim pesan ping secara periodik untuk menjaga koneksi tetap hidup.
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
