package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
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
	log.Printf("Quiz di room '%s' (ID: %d) selesai! Menyimpan hasil...", s.RoomName, s.QuizID)

	finalScores := make(map[string]int)
	winnerUsername := "Seri"
	var winnerID *uint
	maxScore := -1

	// Hitung skor akhir dan tentukan pemenang
	for player, score := range s.PlayerScores {
		finalScores[player.Username] = score
		if score > maxScore {
			maxScore = score
			winnerUsername = player.Username
			// Ambil ID pemenang
			id := player.UserID
			winnerID = &id
		} else if score == maxScore && maxScore != -1 {
			winnerUsername = "Seri"
			winnerID = nil
		}
	}

	// 1. Update setiap participant dengan skor akhirnya
	for _, p := range s.Participants {
		// Cari client yang sesuai dengan participant
		for client, score := range s.PlayerScores {
			if client.UserID == p.UserID {
				updateReq := request.UpdateParticipantRequest{
					TotalScore: score,
				}
				_, err := s.Hub.ParticipantService.UpdateParticipant(int32(p.ID), updateReq)
				if err != nil {
					log.Printf("Gagal update skor untuk participant %d: %v", p.ID, err)
				} else {
					log.Printf("Skor untuk participant %d (User: %d) berhasil diupdate menjadi %d.", p.ID, p.UserID, score)
				}
				break // Lanjut ke participant berikutnya
			}
		}
	}

	// 2. Update status quiz dan pemenangnya
	updateQuizReq := request.UpdateQuizRequest{
		Status:   "completed",
		WinnerID: winnerID,
	}
	// GORM akan handle jika winnerID adalah nil (tidak akan update field)
	// Kita harus memodifikasi QuizService untuk menerima WinnerID
	// Namun, untuk sementara kita bisa modifikasi modelnya langsung di service
	// Mari kita asumsikan service bisa meng-handle ini
	// TODO: Pastikan UpdateQuiz di service bisa update WinnerID
	_, err := s.Hub.QuizService.UpdateQuiz(s.QuizID, updateQuizReq)
	if err != nil {
		log.Printf("Gagal update status quiz %d: %v", s.QuizID, err)
	} else {
		log.Printf("Quiz %d berhasil diupdate dengan status 'completed'. Pemenang: %s", s.QuizID, winnerUsername)
	}

	// Kirim pesan ke client
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
