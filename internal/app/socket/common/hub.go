package common

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

// Session interface for both quiz and prequiz sessions
type Session interface {
	GetState() string
	GetAnswersChannel() chan *AnswerEvent
}

// QuizSessionInterface for quiz sessions
type QuizSessionInterface interface {
	Session
	RunQuizLoop()
}

// PrequizSessionInterface for prequiz sessions
type PrequizSessionInterface interface {
	Session
	RunPrequizLoop()
}

// Factory functions for creating sessions
type QuizSessionFactory func(hub *Hub, roomName string, players []*Client, questions []models.Question, participants []*models.Participant, quizID uint) QuizSessionInterface
type PrequizSessionFactory func(hub *Hub, roomName string, player *Client, questions []models.Prequiz) PrequizSessionInterface

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
	PreQuizService     services.PrequizServiceInterface
	ParticipantService services.ParticipantServiceInterface
	UserService        services.UserServiceInterface // Baru: untuk update lives & xp

	QuizSession    map[string]QuizSessionInterface    // Menyimpan sesi quiz untuk setiap room pada mode 1 vs 1
	PrequizSession map[string]PrequizSessionInterface // Menyimpan sesi pre-quiz untuk setiap room pada mode pre-quiz

	// Factory functions for creating sessions
	QuizSessionFactory    QuizSessionFactory
	PrequizSessionFactory PrequizSessionFactory
}

// NewHub membuat instance Hub baru.
func NewHub(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
	prequizService services.PrequizServiceInterface,
	quizSessionFactory QuizSessionFactory,
	prequizSessionFactory PrequizSessionFactory,
	userService services.UserServiceInterface,
) *Hub {
	return &Hub{
		Message:    make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		// Kita menggunakan map[*Client]bool sebagai 'set' untuk efisiensi.
		Rooms:                 make(map[string]map[*Client]bool),
		RoomNumber:            1,
		QuestionService:       questionService,
		QuizService:           quizService,
		ParticipantService:    participantService,
		PreQuizService:        prequizService,
		UserService:           userService,
		QuizSession:           make(map[string]QuizSessionInterface),
		PrequizSession:        make(map[string]PrequizSessionInterface),
		QuizSessionFactory:    quizSessionFactory,
		PrequizSessionFactory: prequizSessionFactory,
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

			if strings.HasPrefix(roomName, "prequiz-") {
				prequizQuestion, err := h.PreQuizService.GetAllPrequizzes()
				if err != nil || len(prequizQuestion) == 0 {
					log.Printf("Gagal mendapatkan pertanyaan pre-quiz: %v", err)
					h.Mu.Unlock()
					continue
				}

				session := h.PrequizSessionFactory(h, roomName, client, prequizQuestion)
				h.PrequizSession[roomName] = session
				client.Session = session

				go session.RunPrequizLoop()

			}

			if strings.HasPrefix(roomName, "quiz-") && len(h.Rooms[roomName]) == 2 {
				log.Printf("Room kuis '%s' telah penuh. Memulai sesi kuis...", roomName)

				var quizId uint
				fmt.Sscanf(roomName, "quiz-%d", &quizId)

				quiz, err := h.QuizService.GetQuizByID(quizId)

				if err != nil || quiz.Status != "in_progress" {
					log.Printf("Gagal memulai kuis untuk room '%s': Kuis tidak ditemukan atau belum siap. Error: %v", roomName, err)
					h.Mu.Unlock()
					continue
				}

				// Asumsi QuestionService punya fungsi ini
				questionCount := 5 // atau ambil dari quiz.QuestionCount
				questions, err := h.QuestionService.GetRandomQuestionsByModule(*quiz.ModuleID, questionCount)
				if err != nil || len(*questions) == 0 {
					log.Printf("Gagal mendapatkan pertanyaan untuk modul %d: %v", *quiz.ModuleID, err)
					h.Mu.Unlock()
					continue
				}

				players := make([]*Client, 0, 2)
				for p := range h.Rooms[roomName] {
					players = append(players, p)
				}

				participants := make([]*models.Participant, len(quiz.Participants))
				for i, p := range quiz.Participants {
					tempP := p
					participants[i] = &tempP
				}

				session := h.QuizSessionFactory(h, roomName, players, *questions, participants, quiz.ID)

				for _, player := range players {
					player.Session = session
				}

				go session.RunQuizLoop()
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

				if msg.Sender.Session.GetState() == "running" {
					msg.Sender.Session.GetAnswersChannel() <- &AnswerEvent{
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
