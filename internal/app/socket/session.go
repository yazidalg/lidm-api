package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

// Mengelola State dari Quiz
type QuizSession struct {
	Hub                  *Hub
	RoomName             string
	Players              []*Client
	Questions            []models.Question
	CurrentQuestionIndex int
	State                string
}

func (s *QuizSession) RunGameLoop() {
	log.Printf("Memulai quiz di room '%s'...", s.RoomName)
	s.State = "running"

	player1 := s.Players[0]
	player2 := s.Players[1]

	gameStartPayloadP1, _ := json.Marshal(map[string]interface{}{"opponent": map[string]interface{}{"userId": player2.UserID, "username": player2.Username}})
	gameStartPayloadP2, _ := json.Marshal(map[string]interface{}{"opponent": map[string]interface{}{"userId": player1.UserID, "username": player1.Username}})

	s.Hub.SendMessage(player1, &utils.Message{Action: "game_start", Payload: gameStartPayloadP1, Target: s.RoomName})
	s.Hub.SendMessage(player2, &utils.Message{Action: "game_start", Payload: gameStartPayloadP2, Target: s.RoomName})

	time.Sleep(3 * time.Second) // Tunggu 3 detik sebelum pertanyaan pertama

	s.State = "in_progress"

	for i, question := range s.Questions {
		s.CurrentQuestionIndex = i
		log.Printf("Pertanyaan %d: %s", i+1, question.Question)

		questionPayload, _ := json.Marshal(question)
		questionMessage := &utils.Message{
			Action:  "question",
			Payload: questionPayload,
			Target:  s.RoomName,
		}

		s.Hub.BroadcastToRoom(questionMessage)

		// Tunggu sesuai durasi pertanyaan (misal: read_time + answer_time)
		totalDuration := time.Duration(question.AnswerTime+question.ReadTime) * time.Second
		time.Sleep(totalDuration)
	}

	s.State = "finished"
	log.Printf("Quiz di room '%s' selesai!", s.RoomName)
	quizFinishMessage := &utils.Message{Action: "quiz_finished", Payload: json.RawMessage(`{"message":"Permainan Selesai!"}`), Target: s.RoomName}
	s.Hub.BroadcastToRoom(quizFinishMessage)

	s.Hub.RemoveSession(s.RoomName)
}
