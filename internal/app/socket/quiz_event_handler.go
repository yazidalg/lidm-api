package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
)

// Inisialisasi Score untuk setiap pemain
func (s *QuizSession) InitializeScores() {
	for _, player := range s.Players {
		s.PlayerScores[player] = 0
	}
}

// Mengirim pertanyaan ke semua pemain di room
func (s *QuizSession) SendQuestion(question *models.Question) {
	log.Printf("Room '%s': Mengirim pertanyaan #%d", s.RoomName, s.CurrentQuestionIndex+1)
	questionPayload, _ := json.Marshal(question)
	s.Hub.BroadcastToRoom(Message{Action: "question", Payload: questionPayload, Target: s.RoomName})
}

func (s *QuizSession) HandleAnswer(currentQuestion *models.Question) {
	timer := time.NewTimer(time.Duration(currentQuestion.AnswerTime+currentQuestion.ReadTime) * time.Second)
	defer timer.Stop()

questionLoop:
	for {
		select {
		case answerEvent := <-s.Answers:
			s.AnswerProcess(answerEvent, currentQuestion)

			if len(s.PlayerAnswers) == len(s.Players) {
				log.Printf("Room '%s': Semua pemain telah memberikan jawaban untuk pertanyaan #%d", s.RoomName, s.CurrentQuestionIndex+1)
				break questionLoop
			}

		case <-timer.C:
			log.Printf("Room '%s': Waktu habis.", s.RoomName)
			break questionLoop
		}
	}
}

func (s *QuizSession) AnswerProcess(answerEvent *AnswerEvent, currentQuestion *models.Question) {
	if _, answered := s.PlayerAnswers[answerEvent.Player]; answered {
		log.Printf("Player '%s' sudah menjawab pertanyaan #%d", answerEvent.Player.Username, s.CurrentQuestionIndex+1)
		return
	}

	s.PlayerAnswers[answerEvent.Player] = true

	isCorrect := answerEvent.Payload.OptionSelected == currentQuestion.CorrectAnswer
	if isCorrect {
		s.PlayerScores[answerEvent.Player] += 10
	}

	resultPayload, _ := json.Marshal(map[string]interface{}{"is_correct": isCorrect, "your_score": s.PlayerScores[answerEvent.Player]})
	s.Hub.BroadcastToRoom(Message{Action: "answer_result", Payload: resultPayload, Target: s.RoomName})

	answeredEvent := &AnsweredQuestionEvent{
		QuestionID: currentQuestion.ID,
		UserID:     answerEvent.Payload.UserID,
		Score:      s.PlayerScores[answerEvent.Player],
		IsCorrect:  isCorrect,
	}

	answeredEventBytes, _ := json.Marshal(answeredEvent)
	s.Hub.BroadcastToRoom(Message{Action: "answered_question", Payload: answeredEventBytes, Target: s.RoomName})
}

// concludeQuestion mengirimkan rangkuman akhir dari sebuah pertanyaan.
func (s *QuizSession) ConcludeQuestion(question *models.Question) {
	payload, _ := json.Marshal(map[string]interface{}{
		"message":        "Waktu habis atau semua pemain telah menjawab.",
		"correct_answer": question.CorrectAnswer,
	})
	s.Hub.BroadcastToRoom(Message{Action: "question_ended", Payload: payload, Target: s.RoomName})
}

func (s *QuizSession) ConcludeQuiz() {
	s.State = "finished"
	log.Printf("Quiz di room '%s' selesai!", s.RoomName)

	finalScores := make(map[string]int)
	winnerUsername := "Seri"
	maxScore := -1

	for player, score := range s.PlayerScores {
		finalScores[player.Username] = score
		if score > maxScore {
			maxScore = score
			winnerUsername = player.Username
		} else if score == maxScore && maxScore != -1 {
			winnerUsername = "Seri"
		}
	}

	// TODO: Simpan `finalScores` ke database untuk setiap partisipan.

	completedPayload := QuizCompletedPayload{
		FinalScores: finalScores,
		Winner:      winnerUsername,
		Message:     "Permainan Selesai!",
	}
	completedPayloadBytes, _ := json.Marshal(completedPayload)

	quizFinishMessage := Message{Action: "quiz_completed", Payload: completedPayloadBytes, Target: s.RoomName}
	s.Hub.BroadcastToRoom(quizFinishMessage)

	s.Hub.RemoveSession(s.RoomName)
}
