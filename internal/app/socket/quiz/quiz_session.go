package quiz

import (
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
)

// NewQuizSession creates a new quiz session
func NewQuizSession(hub *common.Hub, roomName string, players []*common.Client, questions []models.Question, participants []*models.Participant, quizID uint) *QuizSession {
	return &QuizSession{
		Hub:                  hub,
		RoomName:             roomName,
		Players:              players,
		Questions:            questions,
		Participants:         participants,
		QuizID:               quizID,
		CurrentQuestionIndex: 0,
		State:                "waiting",
		Answers:              make(chan *common.AnswerEvent, 100),
		PlayerAnswers:        make(map[*common.Client]bool),
		PlayerScores:         make(map[*common.Client]int),
	}
}

// Mengelola State dari Quiz
type QuizSession struct {
	Hub                  *common.Hub
	RoomName             string
	Players              []*common.Client
	Questions            []models.Question
	Participants         []*models.Participant
	QuizID               uint
	CurrentQuestionIndex int
	State                string
	Answers              chan *common.AnswerEvent
	PlayerAnswers        map[*common.Client]bool
	PlayerScores         map[*common.Client]int
	QuestionStartTime    time.Time
}

func (s *QuizSession) GetState() string {
	return s.State
}

func (s *QuizSession) GetAnswersChannel() chan *common.AnswerEvent {
	return s.Answers
}

func (s *QuizSession) RunQuizLoop() {
	log.Printf("Room '%s': Memulai sesi quiz", s.RoomName)
	s.State = "running"
	s.InitializeScores()

	log.Println("hahsdhasdhsa", s.QuizID, s.Questions)

	for i, question := range s.Questions {
		s.CurrentQuestionIndex = i
		s.PlayerAnswers = make(map[*common.Client]bool)

		s.SendQuestion(&question)

		s.HandleAnswer(&question)

		s.ConcludeQuestion(&question)
	}

	s.ConcludeQuiz()
}
