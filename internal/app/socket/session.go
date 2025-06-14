package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

// Mengelola State dari Quiz
type QuizSession struct {
	Hub                  *Hub
	RoomName             string
	Players              []*Client
	Questions            []models.Question
	QuizID               uint
	CurrentQuestionIndex int
	State                string
	Answers              chan *AnswerEvent
	PlayerAnswers        map[*Client]bool
}

type AnswerEvent struct {
	Player  *Client
	Payload request.CreateAnswerRequest
}

func (s *QuizSession) RunGameLoop() {
	log.Printf("Memulai quiz di room '%s'...", s.RoomName)
	s.State = "running"

	for i, question := range s.Questions {
		s.CurrentQuestionIndex = i
		s.PlayerAnswers = make(map[*Client]bool)
		log.Printf("Room '%s': Mengirim pertanyaan #%d", s.RoomName, i+1)
		questionPayload, _ := json.Marshal(question)

		s.Hub.BroadcastToRoom(Message{Action: "question", Payload: questionPayload, Target: s.RoomName})

		timer := time.NewTimer(time.Duration(question.AnswerTime+question.ReadTime) * time.Second)

	questionLoop:
		for {
			select {
			case answerEvent := <-s.Answers:
				log.Printf("Room '%s': Player '%s' memberikan jawaban untuk pertanyaan #%d", s.RoomName, answerEvent.Player.Username, i+1)
				s.PlayerAnswers[answerEvent.Player] = true

				if len(s.PlayerAnswers) == len(s.Players) {
					log.Printf("Room '%s': Semua pemain telah memberikan jawaban untuk pertanyaan #%d", s.RoomName, i+1)
					timer.Stop()
					break questionLoop
				}

			case <-timer.C:
				log.Printf("Room '%s': Waktu habis untuk pertanyaan #%d", s.RoomName, i+1)
				break questionLoop
			}

			resultPayload, _ := json.Marshal(map[string]interface{}{"message": "Waktu habis / Semua sudah menjawab", "correctAnswer": question.CorrectAnswer})
			s.Hub.BroadcastToRoom(Message{Action: "question_result", Payload: resultPayload, Target: s.RoomName})

			time.Sleep(5 * time.Second)
		}
	}

	s.State = "finished"
	log.Printf("Quiz di room '%s' selesai!", s.RoomName)
	quizFinishMessage := Message{Action: "quiz_finished", Payload: json.RawMessage(`{"message":"Permainan Selesai!"}`), Target: s.RoomName}
	s.Hub.BroadcastToRoom(quizFinishMessage)

	s.Hub.RemoveSession(s.RoomName)
}
