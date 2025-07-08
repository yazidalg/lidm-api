package socket

import (
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
)

// Mengelola State dari Quiz
type QuizSession struct {
	Hub                  *Hub
	RoomName             string
	Players              []*Client
	Questions            []models.Question
	Participants         []*models.Participant
	PrequizQuestion      []models.Prequiz
	QuizID               uint
	CurrentQuestionIndex int
	State                string
	Answers              chan *AnswerEvent
	PlayerAnswers        map[*Client]bool
	PlayerScores         map[*Client]int
	QuestionStartTime    time.Time
}

type AnswerEvent struct {
	Player  *Client
	Payload AnswerPayload
}

func (s *QuizSession) RunGameLoop() {
	log.Printf("Room '%s': Memulai sesi quiz", s.RoomName)
	s.State = "running"
	s.InitializeScores()

	log.Println("hahsdhasdhsa", s.QuizID, s.Questions)

	for i, question := range s.Questions {
		s.CurrentQuestionIndex = i
		s.PlayerAnswers = make(map[*Client]bool)

		s.SendQuestion(&question)

		s.HandleAnswer(&question)

		s.ConcludeQuestion(&question)
	}

	s.ConcludeQuiz()
}
