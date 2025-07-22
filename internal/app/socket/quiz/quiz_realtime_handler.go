package quiz

import (
	"encoding/json"
	"log"

	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/app/socket/common"
)

type QuizEventHandler struct {
	quizSessionService services.QuizSessionServiceInterface
}

func NewQuizEventHandler(quizSessionService services.QuizSessionServiceInterface) *QuizEventHandler {
	return &QuizEventHandler{
		quizSessionService: quizSessionService,
	}
}

// HandleMessage processes incoming WebSocket messages for quiz
func (h *QuizEventHandler) HandleMessage(client *common.Client, message *common.Message) {
	switch message.Action {
	case "join_quiz":
		h.handleJoinQuiz(client, message)
	case "answer_question":
		h.handleAnswerQuestion(client, message)
	case "get_quiz_state":
		h.handleGetQuizState(client, message)
	case "finish_quiz":
		h.handleFinishQuiz(client, message)
	default:
		h.sendError(client, "unknown_message_type", "Unknown message type")
	}
}

func (h *QuizEventHandler) handleJoinQuiz(client *common.Client, message *common.Message) {
	var req request.JoinQuizRequest
	if err := json.Unmarshal(message.Payload, &req); err != nil {
		h.sendError(client, "invalid_data", "Invalid join quiz data")
		return
	}

	result, err := h.quizSessionService.JoinQuiz(client.UserID, req)
	if err != nil {
		h.sendError(client, "join_failed", err.Error())
		return
	}

	// Send success response to client
	responseData, _ := json.Marshal(result)
	response := &common.Message{
		Action:  "quiz_joined",
		Payload: responseData,
	}

	select {
	case client.Send <- response:
	default:
		close(client.Send)
	}

	// Notify other participants in the room
	h.broadcastQuizUpdate(client, result.QuizID)
}

func (h *QuizEventHandler) handleAnswerQuestion(client *common.Client, message *common.Message) {
	var req request.AnswerQuestionRequest
	if err := json.Unmarshal(message.Payload, &req); err != nil {
		h.sendError(client, "invalid_data", "Invalid answer data")
		return
	}

	result, err := h.quizSessionService.AnswerQuestion(client.UserID, req)
	if err != nil {
		h.sendError(client, "answer_failed", err.Error())
		return
	}

	// Send answer result to client
	responseData, _ := json.Marshal(result)
	response := &common.Message{
		Action:  "answer_submitted",
		Payload: responseData,
	}

	select {
	case client.Send <- response:
	default:
		close(client.Send)
	}

	// Broadcast quiz state update to all participants
	h.broadcastQuizUpdate(client, result.QuizID)
}

func (h *QuizEventHandler) handleGetQuizState(client *common.Client, message *common.Message) {
	var data struct {
		QuizID uint `json:"quiz_id"`
	}

	if err := json.Unmarshal(message.Payload, &data); err != nil {
		h.sendError(client, "invalid_data", "Invalid quiz ID")
		return
	}

	result, err := h.quizSessionService.GetQuizSession(data.QuizID)
	if err != nil {
		h.sendError(client, "get_state_failed", err.Error())
		return
	}

	responseData, _ := json.Marshal(result)
	response := &common.Message{
		Action:  "quiz_state",
		Payload: responseData,
	}

	select {
	case client.Send <- response:
	default:
		close(client.Send)
	}
}

func (h *QuizEventHandler) handleFinishQuiz(client *common.Client, message *common.Message) {
	var data struct {
		QuizID uint `json:"quiz_id"`
	}

	if err := json.Unmarshal(message.Payload, &data); err != nil {
		h.sendError(client, "invalid_data", "Invalid quiz ID")
		return
	}

	err := h.quizSessionService.FinishQuiz(data.QuizID, client.UserID)
	if err != nil {
		h.sendError(client, "finish_failed", err.Error())
		return
	}

	// Send finish confirmation
	responsePayload := map[string]interface{}{
		"quiz_id": data.QuizID,
		"message": "Quiz finished successfully",
	}
	responseData, _ := json.Marshal(responsePayload)
	response := &common.Message{
		Action:  "quiz_finished",
		Payload: responseData,
	}

	select {
	case client.Send <- response:
	default:
		close(client.Send)
	}

	// Get and broadcast final results
	result, err := h.quizSessionService.GetQuizResult(data.QuizID)
	if err == nil {
		resultData, _ := json.Marshal(result)
		resultResponse := &common.Message{
			Action:  "quiz_results",
			Payload: resultData,
		}
		h.broadcastToRoom(client, resultResponse)
	}
}

func (h *QuizEventHandler) broadcastQuizUpdate(client *common.Client, quizID uint) {
	result, err := h.quizSessionService.GetQuizSession(quizID)
	if err != nil {
		log.Printf("Failed to get quiz session for broadcast: %v", err)
		return
	}

	responseData, _ := json.Marshal(result)
	response := &common.Message{
		Action:  "quiz_update",
		Payload: responseData,
	}

	h.broadcastToRoom(client, response)
}

func (h *QuizEventHandler) broadcastToRoom(client *common.Client, message *common.Message) {
	// Get all clients in the same room
	if client.Hub != nil && client.Hub.Rooms != nil {
		if room, exists := client.Hub.Rooms[client.Room]; exists {
			for roomClient := range room {
				if roomClient != client {
					select {
					case roomClient.Send <- message:
					default:
						close(roomClient.Send)
					}
				}
			}
		}
	}
}

func (h *QuizEventHandler) sendError(client *common.Client, errorType, errorMessage string) {
	errorPayload := map[string]interface{}{
		"error_type": errorType,
		"message":    errorMessage,
	}
	errorData, _ := json.Marshal(errorPayload)

	errorResponse := &common.Message{
		Action:  "error",
		Payload: errorData,
	}

	select {
	case client.Send <- errorResponse:
	default:
		close(client.Send)
	}
}
