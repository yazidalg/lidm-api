package prequiz

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
)

func (s *PrequizSession) GetCurrentQuestion(prequiz *models.Prequiz) {
	log.Printf("Room '%s': Mengirim pertanyaan pre-quiz #%d", s.RoomName, s.CurrentQuestionIndex+1)

	prequizPayload, _ := json.Marshal(prequiz)
	s.Hub.SendMessage(s.Player, common.Message{Action: "prequiz_question", Payload: prequizPayload, Target: s.RoomName})
	s.QuestionStartTime = time.Now()
}

func (s *PrequizSession) HandleAnswer(prequiz *models.Prequiz) {
	answerEvent := <-s.Answers
	fmt.Printf("hahahahahaha %v\n", answerEvent.Payload.OptionSelected)
	s.AnswerProcess(answerEvent, prequiz)
}

func (s *PrequizSession) AnswerProcess(answerEvent *common.AnswerEvent, currentPrequiz *models.Prequiz) {
	isCorrect := answerEvent.Payload.OptionSelected == currentPrequiz.CorrectAnswer

	if isCorrect {
		s.Player.RightAnswer++
	} else {
		s.Player.WrongAnswer++
	}

	resultPayload, _ := json.Marshal(map[string]interface{}{
		"is_correct": isCorrect,
	})

	s.Hub.SendMessage(s.Player, common.Message{
		Action:  "prequiz_answer_result",
		Payload: resultPayload,
		Target:  s.RoomName,
	})
}

func (s *PrequizSession) ConcludeQuestion(prequiz *models.Prequiz) {
	payload, _ := json.Marshal(map[string]interface{}{
		"message":        "Pertanyaan pre-quiz selesai",
		"correct_answer": prequiz.CorrectAnswer,
	})

	s.Hub.SendMessage(s.Player, common.Message{
		Action:  "prequiz_question_concluded",
		Payload: payload,
		Target:  s.RoomName,
	})
	time.Sleep(2 * time.Second)
}

func (s *PrequizSession) ConcludePrequiz() {
	s.State = "finished"
	log.Printf("Pre-quiz di room '%s' selesai! Menyimpan hasil...", s.RoomName)

	finalAnswer := map[string]uint{
		"Benar": s.Player.RightAnswer,
		"Salah": s.Player.WrongAnswer,
	}

	payload, _ := json.Marshal(common.PrequizCompletedPayload{
		FinalScores: finalAnswer,
		Message:     "Pre-quiz selesai! Skor akhir telah dihitung.",
	})

	s.Hub.BroadcastToRoom(common.Message{
		Action:  "prequiz_completed",
		Payload: payload,
		Target:  s.RoomName,
	})

	s.Hub.RemoveSession(s.RoomName)
}
