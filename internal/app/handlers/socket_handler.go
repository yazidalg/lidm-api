package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (you may want to restrict this in production)
	},
}

// MatchManager handles the matchmaking and quiz sessions
type MatchManager struct {
	mutex              sync.Mutex
	waitingPlayers     map[uint]*Player
	activeSessions     map[string]*QuizSession
	questionService    services.QuestionServiceInterface
	quizService        services.QuizServiceInterface
	participantService services.ParticipantServiceInterface
	answerService      services.AnswerServiceInterface
}

// Player represents a connected user
type Player struct {
	UserID   uint
	Username string
	Conn     *websocket.Conn
	Session  *QuizSession
}

// QuizSession represents an active quiz between two players
type QuizSession struct {
	SessionID       string
	Player1         *Player
	Player2         *Player
	QuizID          uint
	Questions       []*models.Question
	CurrentQuestion int
	Scores          map[uint]int
	StartTime       time.Time
	IsActive        bool
}

// SocketMessage represents the structure of messages sent over the socket
type SocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// NewMatchManager creates a new match manager
func NewMatchManager(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
	answerService services.AnswerServiceInterface,
) *MatchManager {
	return &MatchManager{
		waitingPlayers:     make(map[uint]*Player),
		activeSessions:     make(map[string]*QuizSession),
		questionService:    questionService,
		quizService:        quizService,
		participantService: participantService,
		answerService:      answerService,
	}
}

// SocketHandler handles websocket connections
type SocketHandler struct {
	matchManager *MatchManager
}

// NewSocketHandler creates a new socket handler
func NewSocketHandler(
	questionService services.QuestionServiceInterface,
	quizService services.QuizServiceInterface,
	participantService services.ParticipantServiceInterface,
	answerService services.AnswerServiceInterface,
) *SocketHandler {
	return &SocketHandler{
		matchManager: NewMatchManager(questionService, quizService, participantService, answerService),
	}
}

// HandleConnection handles a new websocket connection
func (h *SocketHandler) HandleConnection(c *gin.Context) {
	userID := c.Query("user_id")
	username := c.Query("username")

	if userID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and username are required",
		})
		return
	}

	// Parse user ID
	var uid uint
	_, err := fmt.Sscanf(userID, "%d", &uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Upgrade the HTTP connection to a websocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create a new player
	player := &Player{
		UserID:   uid,
		Username: username,
		Conn:     conn,
	}

	// Handle the connection
	go h.handlePlayer(player)
}

// handlePlayer processes messages from a connected player
func (h *SocketHandler) handlePlayer(player *Player) {
	defer func() {
		player.Conn.Close()
		h.matchManager.removePlayer(player)
	}()

	// Send a welcome message
	welcomeMsg := SocketMessage{
		Type: "connected",
		Payload: map[string]interface{}{
			"message": "Connected to quiz matchmaking",
			"user_id": player.UserID,
		},
	}

	if err := player.Conn.WriteJSON(welcomeMsg); err != nil {
		log.Printf("Error sending welcome message: %v", err)
		return
	}

	// Process incoming messages
	for {
		_, message, err := player.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse the message
		var socketMsg SocketMessage
		if err := json.Unmarshal(message, &socketMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch socketMsg.Type {
		case "find_match":
			h.matchManager.findMatch(player)
		case "answer_question":
			if player.Session != nil {
				// Parse the answer payload
				payloadBytes, _ := json.Marshal(socketMsg.Payload)
				var answerPayload struct {
					QuestionID uint   `json:"question_id"`
					Answer     string `json:"answer"`
				}
				json.Unmarshal(payloadBytes, &answerPayload)

				h.matchManager.handleAnswer(player, answerPayload.QuestionID, answerPayload.Answer)
			}
		case "cancel_matchmaking":
			h.matchManager.cancelMatchmaking(player)
		}
	}
}

// findMatch attempts to find a match for the player or adds them to the waiting queue
func (mm *MatchManager) findMatch(player *Player) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Check if there are any waiting players
	if len(mm.waitingPlayers) == 0 {
		// Add the player to the waiting queue
		mm.waitingPlayers[player.UserID] = player

		// Notify the player they're in the queue
		queueMsg := SocketMessage{
			Type: "matchmaking_status",
			Payload: map[string]interface{}{
				"status":  "queued",
				"message": "Waiting for an opponent",
			},
		}
		player.Conn.WriteJSON(queueMsg)
		return
	}

	// Find an opponent (first waiting player)
	var opponent *Player
	for _, p := range mm.waitingPlayers {
		if p.UserID != player.UserID {
			opponent = p
			delete(mm.waitingPlayers, opponent.UserID)
			break
		}
	}

	if opponent == nil {
		// No suitable opponent found, add player to waiting queue
		mm.waitingPlayers[player.UserID] = player

		// Notify the player they're in the queue
		queueMsg := SocketMessage{
			Type: "matchmaking_status",
			Payload: map[string]interface{}{
				"status":  "queued",
				"message": "Waiting for an opponent",
			},
		}
		player.Conn.WriteJSON(queueMsg)
		return
	}

	// Create a new quiz session
	sessionID := generateSessionID()
	session := &QuizSession{
		SessionID: sessionID,
		Player1:   opponent,
		Player2:   player,
		Scores:    make(map[uint]int),
		IsActive:  true,
		StartTime: time.Now(),
	}

	// Load questions for the quiz
	questions, err := mm.loadQuestions(5) // Load 5 questions for the quiz
	if err != nil {
		log.Printf("Failed to load questions: %v", err)

		// Notify both players of the error
		errorMsg := SocketMessage{
			Type: "error",
			Payload: map[string]interface{}{
				"message": "Failed to create quiz session",
			},
		}
		player.Conn.WriteJSON(errorMsg)
		opponent.Conn.WriteJSON(errorMsg)
		return
	}

	session.Questions = questions

	// Update player session references
	player.Session = session
	opponent.Session = session

	// Add session to active sessions
	mm.activeSessions[sessionID] = session

	// Create a quiz record in the database
	quiz, err := mm.createQuizRecord(session)
	if err != nil {
		log.Printf("Failed to create quiz record: %v", err)
	} else {
		session.QuizID = quiz.ID
	}

	// Create participant records
	mm.createParticipants(session)

	// Notify both players that the match has started
	matchStartMsg := SocketMessage{
		Type: "match_started",
		Payload: map[string]interface{}{
			"session_id": sessionID,
			"opponent": map[string]interface{}{
				"user_id":  opponent.UserID,
				"username": opponent.Username,
			},
			"quiz_id": session.QuizID,
		},
	}
	player.Conn.WriteJSON(matchStartMsg)

	matchStartMsg.Payload.(map[string]interface{})["opponent"] = map[string]interface{}{
		"user_id":  player.UserID,
		"username": player.Username,
	}
	opponent.Conn.WriteJSON(matchStartMsg)

	// Start the quiz
	go mm.runQuizSession(session)
}

// runQuizSession sends questions to both players and manages the quiz flow
func (mm *MatchManager) runQuizSession(session *QuizSession) {
	// Wait a moment before starting
	time.Sleep(3 * time.Second)

	// Send questions one by one
	for i, question := range session.Questions {
		// Check if session is still active
		if !session.IsActive {
			return
		}

		session.CurrentQuestion = i

		// Prepare question data without the correct answer
		questionData := map[string]interface{}{
			"question_id": question.ID,
			"text":        question.Question,
			"options": map[string]string{
				"A": question.Options.OptionA,
				"B": question.Options.OptionB,
				"C": question.Options.OptionC,
				"D": question.Options.OptionD,
			},
			"read_time":       question.ReadTime,
			"answer_time":     question.AnswerTime,
			"question_number": i + 1,
			"total_questions": len(session.Questions),
		}

		// Send the question to both players
		questionMsg := SocketMessage{
			Type:    "new_question",
			Payload: questionData,
		}

		session.Player1.Conn.WriteJSON(questionMsg)
		session.Player2.Conn.WriteJSON(questionMsg)

		// Wait for the read time
		time.Sleep(time.Duration(question.ReadTime) * time.Second)

		// Notify that answering period has started
		startAnsweringMsg := SocketMessage{
			Type: "start_answering",
			Payload: map[string]interface{}{
				"question_id":    question.ID,
				"time_remaining": question.AnswerTime,
			},
		}

		session.Player1.Conn.WriteJSON(startAnsweringMsg)
		session.Player2.Conn.WriteJSON(startAnsweringMsg)

		// Wait for the answer time
		time.Sleep(time.Duration(question.AnswerTime) * time.Second)

		// Send correct answer and current scores
		correctAnswerMsg := SocketMessage{
			Type: "question_ended",
			Payload: map[string]interface{}{
				"question_id":    question.ID,
				"correct_answer": question.CorrectAnswer,
				"explanation":    question.Explanation,
				"scores": map[string]interface{}{
					fmt.Sprintf("%d", session.Player1.UserID): session.Scores[session.Player1.UserID],
					fmt.Sprintf("%d", session.Player2.UserID): session.Scores[session.Player2.UserID],
				},
			},
		}

		session.Player1.Conn.WriteJSON(correctAnswerMsg)
		session.Player2.Conn.WriteJSON(correctAnswerMsg)

		// Wait a moment before the next question
		time.Sleep(5 * time.Second)
	}

	// Quiz completed - send final results
	winner := determineWinner(session)

	finalResultMsg := SocketMessage{
		Type: "quiz_completed",
		Payload: map[string]interface{}{
			"final_scores": map[string]interface{}{
				fmt.Sprintf("%d", session.Player1.UserID): session.Scores[session.Player1.UserID],
				fmt.Sprintf("%d", session.Player2.UserID): session.Scores[session.Player2.UserID],
			},
			"winner": winner,
		},
	}

	session.Player1.Conn.WriteJSON(finalResultMsg)
	session.Player2.Conn.WriteJSON(finalResultMsg)

	// Update participant records with final scores
	mm.updateParticipantScores(session)

	// Clean up the session
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	session.IsActive = false
	delete(mm.activeSessions, session.SessionID)
}

// handleAnswer processes a player's answer to a question
func (mm *MatchManager) handleAnswer(player *Player, questionID uint, answer string) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	session := player.Session
	if session == nil || !session.IsActive {
		return
	}

	// Find the current question
	var currentQuestion *models.Question
	if session.CurrentQuestion < len(session.Questions) {
		currentQuestion = session.Questions[session.CurrentQuestion]
	} else {
		return
	}

	// Check if the question ID matches the current question
	if currentQuestion.ID != questionID {
		return
	}

	// Check if the answer is correct
	isCorrect := currentQuestion.CorrectAnswer == answer

	// Update score (add 10 points for correct answer)
	if isCorrect {
		session.Scores[player.UserID] += 10
	}

	// Notify the player about their answer result
	answerResultMsg := SocketMessage{
		Type: "answer_result",
		Payload: map[string]interface{}{
			"question_id": questionID,
			"is_correct":  isCorrect,
			"score":       session.Scores[player.UserID],
		},
	}

	player.Conn.WriteJSON(answerResultMsg)

	// Create an answer record in the database
	mm.createAnswerRecord(player, questionID, answer, isCorrect)
}

// removePlayer removes a player from matchmaking or active sessions
func (mm *MatchManager) removePlayer(player *Player) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Remove from waiting players if present
	delete(mm.waitingPlayers, player.UserID)

	// Handle active session
	if player.Session != nil {
		session := player.Session

		// Notify the other player that the opponent has disconnected
		var otherPlayer *Player
		if session.Player1.UserID == player.UserID {
			otherPlayer = session.Player2
		} else {
			otherPlayer = session.Player1
		}

		if otherPlayer != nil && otherPlayer.Conn != nil {
			disconnectMsg := SocketMessage{
				Type: "opponent_disconnected",
				Payload: map[string]interface{}{
					"message": "Your opponent has disconnected",
				},
			}
			otherPlayer.Conn.WriteJSON(disconnectMsg)
		}

		// Update the session
		session.IsActive = false
		delete(mm.activeSessions, session.SessionID)
	}
}

// cancelMatchmaking removes a player from the matchmaking queue
func (mm *MatchManager) cancelMatchmaking(player *Player) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()

	// Check if the player is in the waiting queue
	if _, exists := mm.waitingPlayers[player.UserID]; exists {
		delete(mm.waitingPlayers, player.UserID)

		// Notify the player
		cancelMsg := SocketMessage{
			Type: "matchmaking_cancelled",
			Payload: map[string]interface{}{
				"message": "Matchmaking cancelled",
			},
		}
		player.Conn.WriteJSON(cancelMsg)
	}
}

// Helper functions

// generateSessionID creates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// loadQuestions loads a set of random questions
func (mm *MatchManager) loadQuestions(count int) ([]*models.Question, error) {
	// This is a simplified version - you should implement proper question selection
	questions, err := mm.questionService.GetAllQuestions()
	if err != nil {
		return nil, err
	}

	// Shuffle and take 'count' questions
	result := make([]*models.Question, 0)
	for i := 0; i < len(questions) && i < count; i++ {
		q := questions[i]
		result = append(result, &q)
	}

	return result, nil
}

// createQuizRecord creates a new quiz record in the database
func (mm *MatchManager) createQuizRecord(session *QuizSession) (*models.Quiz, error) {
	// Extract question IDs
	questionIDs := make([]uint, len(session.Questions))
	for i, q := range session.Questions {
		questionIDs[i] = q.ID
	}

	// Create quiz request
	quizRequest := request.CreateQuizRequest{
		ParticipantID: session.Player1.UserID,
		QuestionsIDs:  questionIDs,
		Status:        "active",
	}

	return mm.quizService.CreateQuiz(quizRequest)
}

// createParticipants creates participant records for both players
func (mm *MatchManager) createParticipants(session *QuizSession) {
	// Create participant for player 1
	player1Request := request.CreateParticipantRequest{
		UserID:     session.Player1.UserID,
		QuizID:     session.QuizID,
		TotalScore: 0,
		TotalXP:    0,
		Result:     "pending",
	}

	// Create participant for player 2
	player2Request := request.CreateParticipantRequest{
		UserID:     session.Player2.UserID,
		QuizID:     session.QuizID,
		TotalScore: 0,
		TotalXP:    0,
		Result:     "pending",
	}

	// Create the participants in the database
	mm.participantService.CreateParticipant(player1Request)
	mm.participantService.CreateParticipant(player2Request)
}

// updateParticipantScores updates the participant records with final scores
func (mm *MatchManager) updateParticipantScores(session *QuizSession) {
	// Update participant for player 1
	player1Result := "lost"
	if session.Scores[session.Player1.UserID] > session.Scores[session.Player2.UserID] {
		player1Result = "won"
	} else if session.Scores[session.Player1.UserID] == session.Scores[session.Player2.UserID] {
		player1Result = "draw"
	}

	player1Request := request.UpdateParticipantRequest{
		UserID:     session.Player1.UserID,
		QuizID:     session.QuizID,
		TotalScore: session.Scores[session.Player1.UserID],
		TotalXP:    calculateXP(session.Scores[session.Player1.UserID]),
		Result:     player1Result,
	}

	// Update participant for player 2
	player2Result := "lost"
	if session.Scores[session.Player2.UserID] > session.Scores[session.Player1.UserID] {
		player2Result = "won"
	} else if session.Scores[session.Player2.UserID] == session.Scores[session.Player1.UserID] {
		player2Result = "draw"
	}

	player2Request := request.UpdateParticipantRequest{
		UserID:     session.Player2.UserID,
		QuizID:     session.QuizID,
		TotalScore: session.Scores[session.Player2.UserID],
		TotalXP:    calculateXP(session.Scores[session.Player2.UserID]),
		Result:     player2Result,
	}

	// Get participant IDs from the database
	player1Participants, _ := mm.participantService.GetParticipantsByUserID(session.Player1.UserID)
	player2Participants, _ := mm.participantService.GetParticipantsByUserID(session.Player2.UserID)

	// Find matching participant records
	for _, p := range player1Participants {
		if p.QuizID == session.QuizID {
			mm.participantService.UpdateParticipant(int32(p.ID), player1Request)
			break
		}
	}

	for _, p := range player2Participants {
		if p.QuizID == session.QuizID {
			mm.participantService.UpdateParticipant(int32(p.ID), player2Request)
			break
		}
	}
}

// createAnswerRecord creates a record of a player's answer
func (mm *MatchManager) createAnswerRecord(player *Player, questionID uint, optionSelected string, isCorrect bool) {
	score := 0
	if isCorrect {
		score = 10
	}

	answerRequest := request.CreateAnswerRequest{
		QuestionID:     questionID,
		UserID:         player.UserID,
		QuizID:         []uint{player.Session.QuizID},
		OptionSelected: optionSelected,
		IsCorrect:      isCorrect,
		Score:          score,
	}

	// Create the answer in the database
	mm.answerService.CreateAnswer(answerRequest)
}

// determineWinner determines the winner of the quiz
func determineWinner(session *QuizSession) map[string]interface{} {
	player1Score := session.Scores[session.Player1.UserID]
	player2Score := session.Scores[session.Player2.UserID]

	if player1Score > player2Score {
		return map[string]interface{}{
			"user_id":  session.Player1.UserID,
			"username": session.Player1.Username,
			"score":    player1Score,
		}
	} else if player2Score > player1Score {
		return map[string]interface{}{
			"user_id":  session.Player2.UserID,
			"username": session.Player2.Username,
			"score":    player2Score,
		}
	} else {
		return map[string]interface{}{
			"draw":  true,
			"score": player1Score,
		}
	}
}

// calculateXP calculates XP based on score
func calculateXP(score int) int {
	// Simple XP calculation - adjust as needed
	return score * 2
}
