package prequiz

import (
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
)

// NewPrequizSession creates a new prequiz session
func NewPrequizSession(hub *common.Hub, roomName string, player *common.Client, questions []models.Prequiz) *PrequizSession {
	return &PrequizSession{
		Hub:                  hub,
		RoomName:             roomName,
		Player:               player,
		Questions:            questions,
		CurrentQuestionIndex: 0,
		State:                "waiting",
		Answers:              make(chan *common.AnswerEvent, 100),
	}
}

type PrequizSession struct {
	Hub                  *common.Hub
	RoomName             string
	Player               *common.Client // Hanya satu pemain
	Questions            []models.Prequiz
	CurrentQuestionIndex int
	State                string
	Answers              chan *common.AnswerEvent // Kita bisa gunakan ulang AnswerEvent
	QuestionStartTime    time.Time
}

func (s *PrequizSession) GetState() string {
	return s.State
}

func (s *PrequizSession) GetAnswersChannel() chan *common.AnswerEvent {
	return s.Answers
}

func (s *PrequizSession) RunPrequizLoop() {
	s.State = "running"
	log.Printf("Room '%s': Memulai sesi pre-quiz", s.RoomName)

	for i, question := range s.Questions {
		s.CurrentQuestionIndex = i

		s.GetCurrentQuestion(&question)

		s.HandleAnswer(&question)

		s.ConcludeQuestion(&question)
	}

	s.ConcludePrequiz()
}
