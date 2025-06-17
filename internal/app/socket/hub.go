package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

// Hub mengelola semua room, client, dan meneruskan pesan.
type Hub struct {
	// Kumpulan client yang terdaftar.
	Clients map[*Client]bool

	// Pesan masuk dari client yang akan di-broadcast ke room yang tepat.
	Message chan Message

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

	QuestionService    services.QuestionServiceInterface
	QuizService        services.QuizServiceInterface
	ParticipantService services.ParticipantServiceInterface
	QuizSession        map[string]*QuizSession // Menyimpan sesi quiz untuk setiap room
}

// NewHub membuat instance Hub baru.
func NewHub(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
) *Hub {
	return &Hub{
		Message:    make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		// Kita menggunakan map[*Client]bool sebagai 'set' untuk efisiensi.
		Rooms:              make(map[string]map[*Client]bool),
		RoomNumber:         1,
		QuestionService:    questionService,
		QuizService:        quizService,
		ParticipantService: participantService,
		QuizSession:        make(map[string]*QuizSession),
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

			user_join_payload := UserEventPayload{
				UserID:   client.UserID,
				Username: client.Username,
				Room:     roomName,
				Message:  fmt.Sprintf("%s telah bergabung ke room %s", client.Username, roomName),
			}

			userJoinPayloadBytes, _ := json.Marshal(user_join_payload)

			joinMessage := Message{
				Action:  "user_join",
				Payload: userJoinPayloadBytes,
				Target:  roomName,
				Sender:  client,
			}

			h.BroadcastToRoom(joinMessage)

			if len(h.Rooms[roomName]) == 2 {
				getQuestion, err := h.QuestionService.GetRandomQuestion(3)
				if err != nil || len(*getQuestion) < 3 {
					log.Printf("Gagal mendapatkan pertanyaan untuk room '%s': %v", roomName, err)
					h.Mu.Unlock()
					continue
				}

				questionIDs := make([]uint, len(*getQuestion))
				for i, q := range *getQuestion {
					questionIDs[i] = q.ID
				}

				createQuizReq := request.CreateQuizRequest{
					Status:       "in_progress", // Langsung set in_progress
					QuestionsIDs: questionIDs,
				}

				newQuiz, err := h.QuizService.CreateQuiz(createQuizReq)

				if err != nil {
					log.Printf("Gagal membuat quiz di database: %v", err)
					h.Mu.Unlock()
					continue
				}
				log.Printf("Quiz dengan ID %d berhasil dibuat di database.", newQuiz.ID)

				players := make([]*Client, 0, 2)
				for p := range h.Rooms[roomName] {
					players = append(players, p)
				}

				participants := make([]*models.Participant, 0, len(players))
				for _, player := range players {
					createParticipantReq := request.CreateParticipantRequest{
						UserID: player.UserID,
						QuizID: newQuiz.ID,
					}
					newParticipant, err := h.ParticipantService.CreateParticipant(createParticipantReq)
					if err != nil {
						log.Printf("Gagal membuat participant di database untuk user %d: %v", player.UserID, err)
						continue // Lanjutkan ke pemain berikutnya
					}
					participants = append(participants, newParticipant)
				}

				if len(participants) != 2 {
					log.Printf("Gagal membuat semua participant, membatalkan quiz.")
					// TODO: Hapus quiz yang sudah terlanjur dibuat
					h.Mu.Unlock()
					continue
				}

				session := &QuizSession{
					Hub:           h,
					RoomName:      roomName,
					Players:       players,
					Questions:     *getQuestion,
					QuizID:        newQuiz.ID,
					State:         "waiting",
					Answers:       make(chan *AnswerEvent, 10), // Buffer untuk jawaban
					PlayerAnswers: make(map[*Client]bool),
					PlayerScores:  make(map[*Client]int),
				}

				h.QuizSession[roomName] = session

				for _, player := range players {
					player.Session = session
				}

				go session.RunGameLoop()
			}

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

					user_leave_payload := UserEventPayload{
						UserID:   client.UserID,
						Username: client.Username,
						Room:     roomName,
						Message:  fmt.Sprintf("%s telah meninggalkan room %s", client.Username, roomName),
					}

					userLeavePayloadBytes, _ := json.Marshal(user_leave_payload)

					leaveMessage := Message{
						Action:  "user_leave",
						Payload: userLeavePayloadBytes,
						Target:  roomName,
						Sender:  client,
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
			if msg.Sender.Session == nil {
				log.Printf("Pesan dari client tanpa sesi: %s", msg.Sender.Username)
				continue
			}

			if msg.Action == "submit_answer" {
				var payload AnswerPayload
				if err := json.Unmarshal(msg.Payload, &payload); err != nil {
					log.Printf("Gagal unmarshal payload: %v", err)
					continue
				}

				if msg.Sender.Session.State == "running" {
					msg.Sender.Session.Answers <- &AnswerEvent{
						Player:  msg.Sender,
						Payload: payload,
					}
				}

			}
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

func (h *Hub) BroadcastToRoom(msg Message) {
	// Cari room yang dituju oleh pesan.
	roomName := msg.Target
	if clients, ok := h.Rooms[roomName]; ok {
		log.Printf("Broadcast pesan ke room '%s': %s", roomName, msg.Payload)
		// Kirim pesan ke setiap client di dalam room.
		for client := range clients {
			// Kirim pesan ke channel 'send' milik client.
			// Client akan menanganinya secara asynchronous.
			client.Send <- &msg
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

// Helper function untuk mengirim pesan ke client tertentu.
func (h *Hub) SendMessage(client *Client, message Message) {
	select {
	case client.Send <- &message:
	default:
		close(client.Send) // Jika channel penuh, tutup untuk menghindari deadlock
	}
}

func (h *Hub) RemoveSession(roomName string) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	delete(h.QuizSession, roomName)
}
