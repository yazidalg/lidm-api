package socketio

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/utils"
	"github.com/zishang520/socket.io/v2/socket"
)

// Event constants
const (
	EventQuestion       = "question"
	EventAnswerSubmit   = "submit_answer"
	EventAnswerResult   = "answer_result"
	EventQuizCompleted  = "quiz_completed"
	EventReviewJawaban  = "review_jawaban"
	EventJoinQuiz       = "join_quiz"
	// New: join single-player quiz by module id (no quiz lobby needed)
	EventJoinQuizByModule = "join_quiz_by_module"
	EventStartQuiz      = "start_quiz"
	EventHostStartQuiz  = "host_start_quiz" // Event untuk host memulai quiz multiplayer
	EventLivesExhausted = "lives_exhausted"
	EventUserJoin       = "user_join"
	EventUserLeave      = "user_leave"
	EventCreateLobby    = "create_lobby"
	EventLobbyCreated   = "lobby_created"
	EventJoinLobby      = "join_lobby"
	EventLobbyJoined    = "lobby_joined"
	EventLobbyFull      = "lobby_full"
	EventLobbyNotFound  = "lobby_not_found"
	EventLobbyError     = "lobby_error"
	EventNextQuestion   = "next_question"
	EventOpponentTaunt  = "opponent_taunt"  // Event for sending taunts/reactions to opponent
	EventTauntReceived  = "taunt_received"  // Event when receiving a taunt from opponent
)

// AnswerReviewItem stores details of each answered question for review
type AnswerReviewItem struct {
	Pertanyaan   string `json:"pertanyaan"`
	JawabanUser  string `json:"jawaban_user"`
	Penjelasan   string `json:"penjelasan"`
	IsCorrect    bool   `json:"is_correct"`
	QuestionType string `json:"question_type"` // "hots" or "regular"
}

// Single player session
type SocketQuizSession struct {
	Socket      *socket.Socket
	Quiz        *models.Quiz
	Questions   []models.Question
	CurrentIdx  int
	Answered    map[uint]bool
	Mu          sync.Mutex
	QuestionSvc services.QuestionServiceInterface
	QuizSvc     services.QuizServiceInterface
	UserSvc     services.UserServiceInterface

	// scoring & summary
	CorrectCount    int
	WrongCount      int
	BoostUsed       int
	PointsEarned    int
	InitialUserPoint int // Poin user sebelum mulai quiz
	UserID          uint // ID user untuk tracking
	AnswerReview    []AnswerReviewItem // Track all answers for review
	// boosters per question id (so FE and BE agree)
	QuestionBoosters map[uint]int
}

func NewSocketQuizSession(sock *socket.Socket, quiz *models.Quiz, questions []models.Question, qSvc services.QuestionServiceInterface, quizSvc services.QuizServiceInterface, userSvc services.UserServiceInterface) *SocketQuizSession {
	return &SocketQuizSession{Socket: sock, Quiz: quiz, Questions: questions, Answered: make(map[uint]bool), QuestionSvc: qSvc, QuizSvc: quizSvc, UserSvc: userSvc}
}

func (s *SocketQuizSession) emitQuestion() {
	if s.CurrentIdx >= len(s.Questions) {
		s.emitSummary("Quiz finished")
		return
	}
	q := s.Questions[s.CurrentIdx]
	// Do not expose full question (avoids leaking correct_answer). Frontend should use questions from start_quiz.
	s.Socket.Emit(EventQuestion, gin.H{"question_id": q.ID, "index": s.CurrentIdx})
}

func (s *SocketQuizSession) handleAnswer(userID uint, option string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	// Set UserID jika belum diset
	if s.UserID == 0 {
		s.UserID = userID
		// Ambil initial point user
		if u, err := s.UserSvc.GetUserByIDUint(userID); err == nil && u != nil {
			s.InitialUserPoint = int(u.Point)
		}
	}
	
	if s.CurrentIdx >= len(s.Questions) {
		return
	}
	q := s.Questions[s.CurrentIdx]
	if s.Answered[q.ID] {
		return
	}
	// Compare by index to be case-insensitive and robust
	correct := optionIndex(option) == optionIndex(q.CorrectAnswer)
	var earnedXP int32 = 0
	
	// Determine mode
	mode := "single_player"
	if s.Quiz != nil && s.Quiz.Mode != "" {
		mode = s.Quiz.Mode
	} else if s.QuizSvc != nil && s.Quiz != nil && s.Quiz.ID != 0 {
		if quiz, err := s.QuizSvc.GetQuizByID(s.Quiz.ID); err == nil && quiz != nil && quiz.Mode != "" {
			mode = quiz.Mode
		}
	}
	
	if correct {
		// compute earned points with booster
		booster := 0
		if s.QuestionBoosters != nil {
			booster = s.QuestionBoosters[q.ID]
		}
		if booster > 0 {
			s.BoostUsed++
		}
		earned := 10
		if booster > 0 {
			earned = earned * booster
		}
		s.CorrectCount++
		s.PointsEarned += earned
		earnedXP = int32(earned)
		_ = s.UserSvc.AddXP(userID, earnedXP)
	} else {
		s.WrongCount++
		// Kurangi lives jika mode single_player
		if mode == "single_player" {
			_ = s.UserSvc.DecrementLife(userID)
		}
	}
	
	// Build result payload
	resultPayload := gin.H{
		"question_id": q.ID,
		"is_correct":  correct,
		"options": gin.H{
			"selected_index": optionIndex(option),
			"correct_index":  optionIndex(q.CorrectAnswer),
		},
		"your_score": s.PointsEarned,
		"gained_xp":  earnedXP,
	}
	
	// Selalu ambil remaining_lives untuk mode single_player
	if mode == "single_player" {
		if u, err := s.UserSvc.GetUserByIDUint(userID); err == nil && u != nil {
			resultPayload["remaining_lives"] = int(u.Lives)
			
			// Emit answer result dulu
			s.Socket.Emit(EventAnswerResult, resultPayload)
			
			// Jika lives habis, langsung finish quiz
			if u.Lives <= 0 {
				s.Socket.Emit(EventLivesExhausted, gin.H{"message": "Lives exhausted", "lives": u.Lives})
				s.emitSummary("Lives exhausted")
				s.CurrentIdx = len(s.Questions)
				return
			}
		} else {
			// Fallback jika gagal ambil user
			s.Socket.Emit(EventAnswerResult, resultPayload)
		}
	} else {
		s.Socket.Emit(EventAnswerResult, resultPayload)
	}
	s.Answered[q.ID] = true
	
	// Track answer for review - get the actual text value of the selected option
	jawabanUserText := getOptionText(q, option)
	s.AnswerReview = append(s.AnswerReview, AnswerReviewItem{
		Pertanyaan:   q.Question,
		JawabanUser:  jawabanUserText,
		Penjelasan:   q.Explanation,
		IsCorrect:    correct,
		QuestionType: q.QuestionType, // Add question type (hots/regular)
	})
	
	// no question_ended event
	s.CurrentIdx++
	// If this was the last question, emit completion immediately to avoid race/missed event
	if s.CurrentIdx >= len(s.Questions) {
		s.emitQuestion()
	} else {
		go func() { time.Sleep(time.Second); s.emitQuestion() }()
	}
}

// emitSummary sends the final quiz_completed payload for single-player
func (s *SocketQuizSession) emitSummary(message string) {
	total := len(s.Questions)
	accuracy := 0.0
	if total > 0 {
		accuracy = float64(s.CorrectCount) / float64(total)
	}
	
	// Calculate percentage (0-100%)
	percentage := int(accuracy * 100)
	if percentage > 100 {
		percentage = 100
	}
	
	// Hitung total boost yang berhasil didapat
	totalBoostValue := 0
	if s.QuestionBoosters != nil {
		for qID := range s.Answered {
			if booster, exists := s.QuestionBoosters[qID]; exists && booster > 0 {
				// Cek apakah pertanyaan ini dijawab benar
				// (asumsi: jika ada di Answered dan contribute ke PointsEarned)
				totalBoostValue += booster
			}
		}
	}

	// Determine points type for consistent UI
	pointsType := "plus"
	if s.PointsEarned < 0 {
		pointsType = "minus"
	}

	// Determine if user "won" based on performance and create positive messages
	userWon := false
	var positiveMessage string
	
	if percentage >= 80 {
		userWon = true
		positiveMessage = "Luar biasa! Skor sempurna! Kamu sangat hebat! 🌟"
	} else if percentage >= 60 {
		userWon = true
		positiveMessage = "Bagus sekali! Kamu berhasil dengan baik! Terus pertahankan! 🎯"
	} else if percentage >= 40 {
		userWon = false
		positiveMessage = "Tidak apa-apa! Kamu sudah berusaha keras. Belajar lagi ya! 📚"
	} else {
		userWon = false
		positiveMessage = "Jangan patah semangat! Setiap kesalahan adalah pelajaran. Ayo coba lagi! 💪"
	}

	s.Socket.Emit(EventQuizCompleted, gin.H{
		"message":         positiveMessage, // Updated to positive message
		"total_questions": total,
		"benar":           s.CorrectCount,
		"salah":           s.WrongCount,
		"percentage":      percentage,
		"boost":           totalBoostValue,
		"winner":          userWon, // true or false based on performance
		"points": gin.H{
			"earned":       s.PointsEarned,
			"delta_amount": s.PointsEarned, // Total poin yang didapat dari quiz ini
			"delta_type":   func() string { if s.PointsEarned >= 0 { return "plus" } else { return "minus" } }(),
		},
		"type": pointsType, // "plus" or "minus" - consistent with multiplayer format
	})
	
	// Emit review jawaban
	s.Socket.Emit(EventReviewJawaban, gin.H{
		"data": s.AnswerReview,
	})
}

// Multiplayer session (simplified)
type MultiplayerSession struct {
	QuizID           uint
	RoomName         string
	Quiz             *models.Quiz
	Questions        []models.Question
	CurrentIdx       int
	State            string
	Players          map[uint]*socket.Socket
	Scores           map[uint]int
	Answers          map[uint]bool // Tracks if user answered
	AnswerResults    map[uint]bool // Tracks if answer was correct (true) or wrong (false)
	
	// Enhanced tracking for detailed quiz_completed
	UserCorrectCount map[uint]int   // Per-user correct answer count
	UserWrongCount   map[uint]int   // Per-user wrong answer count
	UserPointsPlus   map[uint]int   // Per-user positive points earned
	UserPointsMinus  map[uint]int   // Per-user negative points lost
	UserBoostUsed    map[uint]int   // Per-user boost multipliers used
	QuestionBoosters map[uint]int   // Boost multiplier per question
	UserAnswerReview map[uint][]AnswerReviewItem // Per-user answer review
	
	Mutex          sync.Mutex
	QuestionSvc    services.QuestionServiceInterface
	QuizSvc        services.QuizServiceInterface
	ParticipantSvc services.ParticipantServiceInterface
	UserSvc        services.UserServiceInterface
	AnswerCh       chan AnswerSubmission
	TimerCancel    chan struct{}
}

type AnswerSubmission struct {
	UserID, QuestionID uint
	Option             string
}

func NewMultiplayerSession(quiz *models.Quiz, room string, qSvc services.QuestionServiceInterface, quizSvc services.QuizServiceInterface, pSvc services.ParticipantServiceInterface, uSvc services.UserServiceInterface, questions []models.Question) *MultiplayerSession {
	// Initialize boost pattern for questions
	questionBoosters := make(map[uint]int)
	boostPattern := []int{1, 1, 2, 1, 3, 1, 1, 2, 1, 4} // Predefined boost pattern like single-player
	for i, question := range questions {
		if i < len(boostPattern) {
			questionBoosters[question.ID] = boostPattern[i]
		} else {
			questionBoosters[question.ID] = 1 // Default multiplier
		}
	}
	
	return &MultiplayerSession{
		QuizID: quiz.ID, 
		RoomName: room, 
		Quiz: quiz, 
		Questions: questions, 
		Players: make(map[uint]*socket.Socket), 
		Scores: make(map[uint]int), 
		Answers: make(map[uint]bool), 
		AnswerResults: make(map[uint]bool),
		UserCorrectCount: make(map[uint]int),
		UserWrongCount: make(map[uint]int),
		UserPointsPlus: make(map[uint]int),
		UserPointsMinus: make(map[uint]int),
		UserBoostUsed: make(map[uint]int),
		QuestionBoosters: questionBoosters,
		UserAnswerReview: make(map[uint][]AnswerReviewItem),
		QuestionSvc: qSvc, 
		QuizSvc: quizSvc, 
		ParticipantSvc: pSvc, 
		UserSvc: uSvc, 
		AnswerCh: make(chan AnswerSubmission, 10), 
		TimerCancel: make(chan struct{}, 1),
	}
}

func (s *MultiplayerSession) broadcast(io *socket.Server, event string, data any) {
	io.To(socket.Room(s.RoomName)).Emit(event, data)
}
func (s *MultiplayerSession) run(io *socket.Server) {
	if s == nil {
		log.Printf("[socket.io] ERROR: Attempted to run nil MultiplayerSession")
		return
	}

	if len(s.Questions) == 0 {
		log.Printf("[socket.io] ERROR: MultiplayerSession has no questions, quiz_id=%d", s.QuizID)
		return
	}

	if s.AnswerCh == nil {
		log.Printf("[socket.io] Initializing AnswerCh for quiz_id=%d", s.QuizID)
		s.AnswerCh = make(chan AnswerSubmission, 10)
	}

	log.Printf("[socket.io] Starting multiplayer session for quiz_id=%d with %d questions", s.QuizID, len(s.Questions))
	s.State = "running"

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[socket.io] PANIC in multiplayer session: %v", r)
		}
	}()

	for s.CurrentIdx < len(s.Questions) {
		s.runQuestion(io)
		s.CurrentIdx++
	}

	s.finish(io)
}
func (s *MultiplayerSession) runQuestion(io *socket.Server) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[socket.io] PANIC in runQuestion: %v", r)
		}
	}()

	if s == nil {
		log.Printf("[socket.io] ERROR: Attempted to run nil MultiplayerSession in runQuestion")
		return
	}

	if s.CurrentIdx < 0 || s.CurrentIdx >= len(s.Questions) {
		log.Printf("[socket.io] ERROR: Invalid question index %d (total: %d) for quiz_id=%d",
			s.CurrentIdx, len(s.Questions), s.QuizID)
		return
	}

	if s.Answers == nil {
		s.Answers = make(map[uint]bool)
	}

	q := s.Questions[s.CurrentIdx]
	log.Printf("[socket.io] Running question %d/%d (ID: %d) for quiz_id=%d",
		s.CurrentIdx+1, len(s.Questions), q.ID, s.QuizID)

	// Reset per-question answers
	s.Answers = make(map[uint]bool)
	// Do not broadcast full question to avoid leaking correct_answer; clients already have questions from start_quiz
	s.broadcast(io, EventQuestion, gin.H{"question_id": q.ID, "index": s.CurrentIdx})

	// Ensure AnswerCh channel exists
	if s.AnswerCh == nil {
		s.AnswerCh = make(chan AnswerSubmission, 10)
	}

	// Wait until all players have answered (no timer)
	for {
		ans := <-s.AnswerCh
		if ans.QuestionID != q.ID {
			log.Printf("[socket.io] Ignoring answer for different question: got %d, expected %d", ans.QuestionID, q.ID)
			continue
		}
		s.processAnswer(io, q, ans)

		playerCount := 0
		if s.Players != nil {
			playerCount = len(s.Players)
		}
		if playerCount > 0 && len(s.Answers) >= playerCount {
			log.Printf("[socket.io] All %d players answered, moving to next question", playerCount)
			break
		}
	}

	// no question_ended broadcast
}
func (s *MultiplayerSession) processAnswer(io *socket.Server, q models.Question, ans AnswerSubmission) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[socket.io] PANIC in processAnswer: %v", r)
		}
	}()

	if s == nil {
		log.Printf("[socket.io] ERROR: Attempted to process answer with nil MultiplayerSession")
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Initialize maps if they don't exist
	if s.Answers == nil {
		s.Answers = make(map[uint]bool)
	}

	if s.Scores == nil {
		s.Scores = make(map[uint]int)
	}

	// Check if user already answered
	if s.Answers[ans.UserID] {
		log.Printf("[socket.io] User %d already answered question %d", ans.UserID, q.ID)
		return
	}

	log.Printf("[socket.io] Processing answer from user %d for question %d: option=%s",
		ans.UserID, q.ID, ans.Option)

	// Initialize tracking maps if needed
	if s.UserCorrectCount[ans.UserID] == 0 && s.UserWrongCount[ans.UserID] == 0 {
		s.UserCorrectCount[ans.UserID] = 0
		s.UserWrongCount[ans.UserID] = 0
		s.UserPointsPlus[ans.UserID] = 0
		s.UserPointsMinus[ans.UserID] = 0
		s.UserBoostUsed[ans.UserID] = 0
		s.UserAnswerReview[ans.UserID] = []AnswerReviewItem{}
	}

	// Compare by index to avoid case/key mismatches
	correct := optionIndex(ans.Option) == optionIndex(q.CorrectAnswer)
	
	// Get boost multiplier for this question
	booster := 1
	if s.QuestionBoosters != nil {
		if boost, exists := s.QuestionBoosters[q.ID]; exists {
			booster = boost
			if boost > 1 {
				s.UserBoostUsed[ans.UserID]++
			}
		}
	}
	
	// Calculate points with boost
	basePoints := 10
	actualPoints := basePoints
	if correct {
		actualPoints = basePoints * booster
		s.Scores[ans.UserID] += actualPoints
		s.UserCorrectCount[ans.UserID]++
		s.UserPointsPlus[ans.UserID] += actualPoints
		
		if s.UserSvc != nil {
			_ = s.UserSvc.AddXP(ans.UserID, int32(actualPoints))
		} else {
			log.Printf("[socket.io] WARNING: UserService is nil, can't add XP")
		}
		log.Printf("[socket.io] Correct answer from user %d, boost: %dx, points: +%d, new score: %d", 
			ans.UserID, booster, actualPoints, s.Scores[ans.UserID])
	} else {
		// Salah: kurangi 10 poin (bisa jadi minus) - no boost on wrong answers
		actualPoints = -basePoints
		s.Scores[ans.UserID] += actualPoints
		s.UserWrongCount[ans.UserID]++
		s.UserPointsMinus[ans.UserID] += basePoints // Store as positive value for display
		
		log.Printf("[socket.io] Incorrect answer from user %d, answer was: %s, correct: %s, points: %d, new score: %d",
			ans.UserID, ans.Option, q.CorrectAnswer, actualPoints, s.Scores[ans.UserID])
	}

	// Add to answer review - get the actual text value of the selected option
	jawabanUserText := getOptionText(q, ans.Option)
	reviewItem := AnswerReviewItem{
		Pertanyaan:   q.Question,
		JawabanUser:  jawabanUserText, // Use actual option text instead of letter
		Penjelasan:   q.Explanation,
		IsCorrect:    correct,
		QuestionType: q.QuestionType, // Add question type (hots/regular)
	}
	s.UserAnswerReview[ans.UserID] = append(s.UserAnswerReview[ans.UserID], reviewItem)

	s.Answers[ans.UserID] = true
	s.AnswerResults[ans.UserID] = correct // Store whether answer was correct
	// Emit answer_result only to the answering user's socket (no leakage)
	if sock, ok := s.Players[ans.UserID]; ok && sock != nil {
		sock.Emit(EventAnswerResult, gin.H{
			"user_id":     ans.UserID,
			"question_id": q.ID,
			"is_correct":  correct,
			"score":       s.Scores[ans.UserID],
			"options": gin.H{
				"selected_index": optionIndex(ans.Option),
				"correct_index":  optionIndex(q.CorrectAnswer),
			},
		})
	}

	// If all players have answered, emit next_question with opponent-only score (no correct_answer exposure)
	playerCount := len(s.Players)
	if playerCount > 0 && len(s.Answers) >= playerCount {
		for uid, sock := range s.Players {
			if sock == nil {
				continue
			}
			// Get current user's score
			myScore := s.Scores[uid]
			myScoreType := "plus"
			if myScore < 0 {
				myScoreType = "minus"
			} else if myScore == 0 {
				myScoreType = "zero"
			}
			
			// Get opponent's data
			oppScore := 0
			var oppID uint
			var oppAnswerCorrect *bool // nil if opponent hasn't answered yet
			for id, sc := range s.Scores {
				if id != uid {
					oppScore = sc
					oppID = id
					// Get opponent's answer result from previous question
					if result, exists := s.AnswerResults[id]; exists {
						oppAnswerCorrect = &result
					}
				}
			}
			
			// Determine opponent score type
			oppScoreType := "plus"
			if oppScore < 0 {
				oppScoreType = "minus"
			} else if oppScore == 0 {
				oppScoreType = "zero"
			}
			
			opponentData := gin.H{
				"user_id":    oppID,
				"score":      oppScore,
				"score_type": oppScoreType,
			}
			// Add answer status if opponent answered
			if oppAnswerCorrect != nil {
				opponentData["answer_correct"] = *oppAnswerCorrect
			}
			
			sock.Emit(EventNextQuestion, gin.H{
				"current":        gin.H{"question_id": q.ID},
				"your_score":     myScore,
				"your_score_type": myScoreType,
				"opponent":       opponentData,
			})
		}
	}
}
func (s *MultiplayerSession) finish(io *socket.Server) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[socket.io] PANIC in finish: %v", r)
		}
	}()

	if s == nil {
		log.Printf("[socket.io] ERROR: Attempted to finish nil MultiplayerSession")
		return
	}

	log.Printf("[socket.io] Finishing quiz session for quiz_id=%d", s.QuizID)
	s.State = "finished"

	winnerName := "Seri"
	var winnerID *uint
	maxScore := -1

	// Initialize scores map if nil
	if s.Scores == nil {
		s.Scores = make(map[uint]int)
	}

	// Check if services are available
	if s.QuizSvc == nil {
		log.Printf("[socket.io] ERROR: QuizService is nil in finish()")
		s.broadcastMultiplayerCompletion(io, winnerName)
		return
	}

	quiz, err := s.QuizSvc.GetQuizByID(s.QuizID)
	if err != nil || quiz == nil {
		log.Printf("[socket.io] ERROR: Failed to get quiz %d: %v", s.QuizID, err)
		s.broadcastMultiplayerCompletion(io, winnerName)
		return
	}

	// Check if ParticipantService is available
	if s.ParticipantSvc == nil {
		log.Printf("[socket.io] ERROR: ParticipantService is nil in finish()")
		s.broadcastMultiplayerCompletion(io, winnerName)
		return
	}

	// Calculate scores and determine winner
	for _, part := range quiz.Participants {
		score := s.Scores[part.UserID]
		log.Printf("[socket.io] Participant %s (ID: %d) scored: %d",
			part.User.Name, part.UserID, score)

		updateReq := request.UpdateParticipantRequest{TotalScore: score}
		_, err := s.ParticipantSvc.UpdateParticipant(int32(part.ID), updateReq)
		if err != nil {
			log.Printf("[socket.io] Failed to update participant score: %v", err)
		}

		if score > maxScore {
			maxScore = score
			winnerName = part.User.Name
			tempID := part.UserID
			winnerID = &tempID
			log.Printf("[socket.io] New leader: %s with score %d", winnerName, score)
		} else if score == maxScore && maxScore != -1 {
			winnerName = "Seri"
			winnerID = nil
			log.Printf("[socket.io] Tie detected with score %d", score)
		}
	}

	// Update quiz status
	updateQuizReq := request.UpdateQuizRequest{Status: "completed", WinnerID: winnerID}
	_, err = s.QuizSvc.UpdateQuiz(s.QuizID, updateQuizReq)
	if err != nil {
		log.Printf("[socket.io] Failed to update quiz status: %v", err)
	}

	log.Printf("[socket.io] Quiz %d completed, winner: %s", s.QuizID, winnerName)
	
	// Create comprehensive quiz_completed payload for multiplayer
	s.broadcastMultiplayerCompletion(io, winnerName)
}

// broadcastMultiplayerCompletion sends detailed quiz_completed event to all multiplayer participants
func (s *MultiplayerSession) broadcastMultiplayerCompletion(io *socket.Server, winnerName string) {
	totalQuestions := len(s.Questions)
	
	// Prepare per-user data for each participant
	for userID, sock := range s.Players {
		if sock == nil {
			continue
		}
		
		// Get opponent data
		var opponentID uint
		var opponentScore int
		var opponentCorrect int
		var opponentWrong int
		
		for uid := range s.Players {
			if uid != userID {
				opponentID = uid
				opponentScore = s.Scores[uid]
				opponentCorrect = s.UserCorrectCount[uid]
				opponentWrong = s.UserWrongCount[uid]
				break
			}
		}
		
		// Get user names
		var userName, opponentName string
		if s.UserSvc != nil {
			if user, err := s.UserSvc.GetUserById(int(userID)); err == nil {
				userName = getFirstName(user.Name)
			}
			if opponentID > 0 {
				if opponent, err := s.UserSvc.GetUserById(int(opponentID)); err == nil {
					opponentName = getFirstName(opponent.Name)
				}
			}
		}
		

		
		// Calculate total boost multipliers used (sum of all boost values > 1)
		userTotalBoost := 0
		opponentTotalBoost := 0
		if s.QuestionBoosters != nil {
			for qID, booster := range s.QuestionBoosters {
				if booster > 1 {
					// Check if user answered this question correctly to get boost
					for _, review := range s.UserAnswerReview[userID] {
						if review.Pertanyaan == s.getQuestionByID(qID).Question && review.IsCorrect {
							userTotalBoost += booster
							break
						}
					}
					// Check opponent
					for _, review := range s.UserAnswerReview[opponentID] {
						if review.Pertanyaan == s.getQuestionByID(qID).Question && review.IsCorrect {
							opponentTotalBoost += booster
							break
						}
					}
				}
			}
		}
		
		// Emit simplified quiz_completed to this user
		userPoints := s.Scores[userID]
		opponentPoints := opponentScore
		
		// Determine if this user won and create positive messages
		userWon := false
		var message string
		
		if winnerName == "Seri" {
			message = "Pertandingan seri! Kalian berdua hebat! 🤝"
		} else {
			// Check if current user won by comparing scores
			if userPoints > opponentPoints {
				userWon = true
				message = "Selamat! Kamu menang! Terus tingkatkan kemampuanmu! 🎉"
			} else {
				userWon = false
				message = "Jangan menyerah! Kamu sudah berusaha dengan baik. Coba lagi ya! 💪"
			}
		}
		
		// Calculate user's percentage  
		userCorrect := s.UserCorrectCount[userID]
		userWrong := s.UserWrongCount[userID]
		userPercentage := 0
		if totalQuestions > 0 {
			userPercentage = int((float64(userCorrect) / float64(totalQuestions)) * 100)
		}
		
		// Calculate opponent's percentage
		opponentPercentage := 0
		if totalQuestions > 0 {
			opponentPercentage = int((float64(opponentCorrect) / float64(totalQuestions)) * 100)
		}
		
		// Calculate boost for user
		userBoost := 0
		if s.QuestionBoosters != nil {
			for qID, booster := range s.QuestionBoosters {
				if booster > 1 {
					// Check if user answered this question correctly to get boost
					for _, review := range s.UserAnswerReview[userID] {
						if review.Pertanyaan == s.getQuestionByID(qID).Question && review.IsCorrect {
							userBoost += booster
							break
						}
					}
				}
			}
		}
		
		// Calculate boost for opponent
		opponentBoost := 0
		if s.QuestionBoosters != nil {
			for qID, booster := range s.QuestionBoosters {
				if booster > 1 {
					// Check if opponent answered this question correctly to get boost
					for _, review := range s.UserAnswerReview[opponentID] {
						if review.Pertanyaan == s.getQuestionByID(qID).Question && review.IsCorrect {
							opponentBoost += booster
							break
						}
					}
				}
			}
		}

		// Calculate earned points (points gained/lost in this quiz)
		userEarned := s.UserPointsPlus[userID] - s.UserPointsMinus[userID]
		userEarnedType := "plus"
		if userEarned < 0 {
			userEarnedType = "minus"
		}
		
		opponentEarned := s.UserPointsPlus[opponentID] - s.UserPointsMinus[opponentID]
		opponentEarnedType := "plus"
		if opponentEarned < 0 {
			opponentEarnedType = "minus"
		}
		
		// Get total scores (accumulated score from database/previous games)
		var userTotalScore int32 = 0
		var opponentTotalScore int32 = 0
		
		if s.UserSvc != nil {
			if user, err := s.UserSvc.GetUserByIDUint(userID); err == nil && user != nil {
				userTotalScore = user.TotalXP
			}
			if opponent, err := s.UserSvc.GetUserByIDUint(opponentID); err == nil && opponent != nil {
				opponentTotalScore = opponent.TotalXP
			}
		}

		sock.Emit(EventQuizCompleted, gin.H{
			"scores": s.Scores,
			"winner": userWon, // true or false
			"message": message, // Positive message
			"user": gin.H{
				"name":           userName,            // User's first name
				"total_score":    userTotalScore,      // Total accumulated score (like 1250)
				"earned_points":  userEarned,          // Points earned/lost this quiz (like +12, -4)
				"earned_type":    userEarnedType,      // "plus" or "minus"
				"benar":          userCorrect,
				"salah":          userWrong,
				"percentage":     userPercentage,
				"boost":          userBoost,
			},
			"opponent": gin.H{
				"id":             opponentID,
				"name":           opponentName,        // Opponent's first name
				"total_score":    opponentTotalScore,  // Total accumulated score (like 976)
				"earned_points":  opponentEarned,      // Points earned/lost this quiz (like +12, -4)
				"earned_type":    opponentEarnedType,  // "plus" or "minus"
				"benar":          opponentCorrect,
				"salah":          opponentWrong,
				"percentage":     opponentPercentage,
				"boost":          opponentBoost,
			},
		})
		
		// Emit answer review for this user
		sock.Emit(EventReviewJawaban, gin.H{
			"data": s.UserAnswerReview[userID],
		})
	}
}

// Helper method to get question by ID
func (s *MultiplayerSession) getQuestionByID(questionID uint) models.Question {
	for _, q := range s.Questions {
		if q.ID == questionID {
			return q
		}
	}
	return models.Question{} // Return empty question if not found
}

// getFirstName extracts the first name from a full name
func getFirstName(fullName string) string {
	if fullName == "" {
		return ""
	}
	parts := strings.Fields(strings.TrimSpace(fullName))
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

// StartSocketIOServer bootstraps the v3 socket.io server
func StartSocketIOServer(router *gin.Engine, questionSvc services.QuestionServiceInterface, quizSvc services.QuizServiceInterface, userSvc services.UserServiceInterface, participantSvc services.ParticipantServiceInterface) {
	opts := socket.DefaultServerOptions()
	opts.SetServeClient(false)
	// IMPORTANT: NewServer expects (httpServer, options). We passed (options, nil) before which caused
	// the library to treat *ServerOptions as an http handler and panic: "attach socket.io to an express request handler".
	// Pass nil for http server (we'll mount the handler manually) and opts as the second argument.
	io := socket.NewServer(nil, opts)
	sessions := struct {
		sync.Mutex
		M map[uint]*MultiplayerSession
	}{M: make(map[uint]*MultiplayerSession)}

	io.On("connection", func(cs ...any) {
		client := cs[0].(*socket.Socket)
		log.Printf("[socket.io] connected id=%s", client.Id())
		client.On("message", func(args ...any) { log.Printf("[socket.io] message args=%#v", args) })

		var singleSess *SocketQuizSession
		var userIDUint uint

		// helper: ensure multiplayer session exists and return it
		getOrCreateSession := func(q *models.Quiz) *MultiplayerSession {
			sessions.Lock()
			defer sessions.Unlock()
			if s, ok := sessions.M[q.ID]; ok && s != nil {
				return s
			}
			s := NewMultiplayerSession(q, roomName(q.ID), questionSvc, quizSvc, participantSvc, userSvc, nil)
			sessions.M[q.ID] = s
			return s
		}

		// create_lobby
		client.On(EventCreateLobby, func(args ...any) {
			log.Printf("[socket.io] create_lobby raw=%#v", args)
			var moduleIDF, hostIDF float64
			questionCount := 10
			var okM, okH bool
			if len(args) == 1 { // object or json string
				switch v := args[0].(type) {
				case map[string]any:
					if vv, ok := v["module_id"].(float64); ok {
						moduleIDF, okM = vv, true
					}
					if vv, ok := v["host_user_id"].(float64); ok {
						hostIDF, okH = vv, true
					}
					if vv, ok := v["question_count"].(float64); ok && int(vv) > 0 {
						questionCount = int(vv)
					}
				case string:
					s := strings.TrimSpace(v)
					if strings.HasPrefix(s, "{") {
						var obj map[string]any
						if err := json.Unmarshal([]byte(s), &obj); err == nil {
							if vv, ok := obj["module_id"].(float64); ok {
								moduleIDF, okM = vv, true
							}
							if vv, ok := obj["host_user_id"].(float64); ok {
								hostIDF, okH = vv, true
							}
							if vv, ok := obj["question_count"].(float64); ok && int(vv) > 0 {
								questionCount = int(vv)
							}
						}
					}
				}
			} else if len(args) >= 2 {
				moduleIDF, okM = args[0].(float64)
				hostIDF, okH = args[1].(float64)
				if len(args) >= 3 {
					if qc, ok := args[2].(float64); ok && int(qc) > 0 {
						questionCount = int(qc)
					}
				}
			}
			if !okM || !okH {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_args", "args": args})
				return
			}
			req := request.CreateQuizRequest{ModuleID: uint(moduleIDF), HostUserID: uint(hostIDF), Mode: "multiplayer", Status: "pending", InviteCode: utils.GenerateInviteCode(6), QuestionCount: questionCount}
			quiz, err := quizSvc.CreateQuiz(req)
			if err != nil || quiz == nil {
				client.Emit(EventLobbyError, gin.H{"error": "create_failed"})
				return
			}
			if _, errP := participantSvc.CreateParticipant(request.CreateParticipantRequest{UserID: req.HostUserID, QuizID: quiz.ID}); errP != nil {
				log.Printf("[socket.io] participant host create err=%v", errP)
			}
			// Join host socket to room and register in session as player
			room := roomName(quiz.ID)
			client.Join(socket.Room(room))
			sess := getOrCreateSession(quiz)
			if sess.Players == nil {
				sess.Players = make(map[uint]*socket.Socket)
			}
			sess.Players[req.HostUserID] = client
			io.To(socket.Room(room)).Emit(EventUserJoin, gin.H{"user_id": req.HostUserID})

			// Host has created lobby, but don't emit lobby_joined until opponent joins
			log.Printf("[socket.io] Host %d created lobby for quiz %d, waiting for opponent", req.HostUserID, quiz.ID)
			client.Emit(EventLobbyCreated, gin.H{"quiz_id": quiz.ID, "invite_code": quiz.InviteCode, "module_id": quiz.ModuleID, "question_count": quiz.QuestionCount})
		})

		// join_lobby
		client.On(EventJoinLobby, func(args ...any) {
			log.Printf("[socket.io] join_lobby raw=%#v", args)
			var code string
			var okC bool
			var userID uint

			if len(args) == 1 {
				switch v := args[0].(type) {
				case map[string]any:
					log.Printf("[socket.io] join_lobby got map: %#v", v)
					if vv, ok := v["invite_code"].(string); ok {
						code, okC = vv, true
					}
					if vv, ok := v["user_id"].(float64); ok {
						userID = uint(vv)
					} else if vs, ok := v["user_id"].(string); ok {
						if n, err := strconv.ParseUint(vs, 10, 64); err == nil {
							userID = uint(n)
						}
					}
				case string:
					log.Printf("[socket.io] join_lobby got string: %s", v)
					s := strings.TrimSpace(v)
					if strings.HasPrefix(s, "{") {
						var obj map[string]any
						if err := json.Unmarshal([]byte(s), &obj); err == nil {
							log.Printf("[socket.io] join_lobby parsed JSON: %#v", obj)
							if vv, ok := obj["invite_code"].(string); ok {
								code, okC = vv, true
							}
							if vv, ok := obj["user_id"].(float64); ok {
								userID = uint(vv)
							} else if vs, ok := obj["user_id"].(string); ok {
								if n, err := strconv.ParseUint(vs, 10, 64); err == nil {
									userID = uint(n)
								}
							}
						} else {
							log.Printf("[socket.io] join_lobby JSON parse error: %v", err)
						}
					}
				default:
					log.Printf("[socket.io] join_lobby unexpected type: %T", v)
				}
			} else if len(args) >= 2 {
				log.Printf("[socket.io] join_lobby got multiple args")
				code, okC = args[0].(string)
				if uidf, ok := args[1].(float64); ok {
					userID = uint(uidf)
				} else if uids, ok := args[1].(string); ok {
					if n, err := strconv.ParseUint(uids, 10, 64); err == nil {
						userID = uint(n)
					}
				}
			}

			if !okC || userID == 0 {
				log.Printf("[socket.io] join_lobby ERROR: invalid arguments")
				client.Emit(EventLobbyError, gin.H{"error": "invalid_args", "args": args})
				return
			}

			quiz, err := quizSvc.GetQuizByInviteCode(code)
			if err != nil || quiz == nil {
				log.Printf("[socket.io] join_lobby invite_code=%s lookup_err=%v", code, err)
				client.Emit(EventLobbyNotFound, gin.H{"invite_code": code})
				return
			}
			if quiz.Status != "pending" {
				log.Printf("[socket.io] join_lobby invite_code=%s status=%s (not pending)", code, quiz.Status)
				client.Emit(EventLobbyNotFound, gin.H{"invite_code": code, "status": quiz.Status})
				return
			}

			// Note: we'll join room and register socket only after eligibility confirmed
			// Check if already joined
			for _, p := range quiz.Participants {
				if p.UserID == userID {
					// Find opponent info
					var opponent *models.User
					for _, other := range quiz.Participants {
						if other.UserID != userID {
							opponent = &other.User
						}
					}
					// Eligible: ensure socket joins room and is tracked as player
					room := roomName(quiz.ID)
					client.Join(socket.Room(room))
					sess := getOrCreateSession(quiz)
					if sess.Players == nil {
						sess.Players = make(map[uint]*socket.Socket)
					}
					sess.Players[userID] = client
					io.To(socket.Room(room)).Emit(EventUserJoin, gin.H{"user_id": userID})
					// Build payload with user/opponent details
					var me, opp map[string]any
					if u, err := userSvc.GetUserByIDUint(userID); err == nil && u != nil {
						me = map[string]any{"id": u.ID, "name": u.Name, "total_xp": u.TotalXP, "profile_picture": u.ProfilePicture}
					}
					if opponent != nil {
						opp = map[string]any{"id": opponent.ID, "name": opponent.Name, "total_xp": opponent.TotalXP, "profile_picture": opponent.ProfilePicture}
					}
					// Check if current user is host (first participant)
					isHost := len(quiz.Participants) > 0 && quiz.Participants[0].UserID == userID
					client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user": me, "opponent": opp, "is_host": isHost})
					return
				}
			}
			if len(quiz.Participants) >= 2 {
				client.Emit(EventLobbyFull, gin.H{"invite_code": code})
				return
			}
			if _, err = participantSvc.CreateParticipant(request.CreateParticipantRequest{UserID: userID, QuizID: quiz.ID}); err != nil {
				client.Emit(EventLobbyError, gin.H{"error": "participant_failed"})
				return
			}
			// Eligible new participant: join room and register socket
			room := roomName(quiz.ID)
			client.Join(socket.Room(room))
			sess := getOrCreateSession(quiz)
			if sess.Players == nil {
				sess.Players = make(map[uint]*socket.Socket)
			}
			sess.Players[userID] = client
			io.To(socket.Room(room)).Emit(EventUserJoin, gin.H{"user_id": userID})

			// Refresh quiz participants
			quiz, _ = quizSvc.GetQuizByInviteCode(code)
			if len(quiz.Participants) == 2 {
				_, _ = quizSvc.UpdateQuiz(quiz.ID, request.UpdateQuizRequest{Status: "in_progress"})
				// Emit lobby_joined tailored to each player's socket with full details
				// Fetch user details first to ensure correct fields
				var u1, u2 *models.User
				if len(quiz.Participants) == 2 {
					if uu, err := userSvc.GetUserByIDUint(quiz.Participants[0].UserID); err == nil {
						u1 = uu
					}
					if uu, err := userSvc.GetUserByIDUint(quiz.Participants[1].UserID); err == nil {
						u2 = uu
					}
				}
				// Emit to participant[0] (host)
				if u1 != nil {
					if s1, ok := sess.Players[u1.ID]; ok && s1 != nil {
						s1.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID,
							"user": map[string]any{"id": u1.ID, "name": u1.Name, "total_xp": u1.TotalXP, "profile_picture": u1.ProfilePicture},
							"opponent": func() any {
								if u2 == nil {
									return nil
								}
								return map[string]any{"id": u2.ID, "name": u2.Name, "total_xp": u2.TotalXP, "profile_picture": u2.ProfilePicture}
							}(),
							"is_host": true, // First participant is always host
						})
					}
				}
				// Emit to participant[1] (joiner)
				if u2 != nil {
					if s2, ok := sess.Players[u2.ID]; ok && s2 != nil {
						s2.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID,
							"user": map[string]any{"id": u2.ID, "name": u2.Name, "total_xp": u2.TotalXP, "profile_picture": u2.ProfilePicture},
							"opponent": func() any {
								if u1 == nil {
									return nil
								}
								return map[string]any{"id": u1.ID, "name": u1.Name, "total_xp": u1.TotalXP, "profile_picture": u1.ProfilePicture}
							}(),
							"is_host": false, // Second participant is not host
						})
					}
				}
				// Emit random XP bonus to both
				bonus := []int{2, 3, 5}[time.Now().UnixNano()%3]
				io.To(socket.Room(roomName(quiz.ID))).Emit("bonus_multiplier", gin.H{"multiplier": bonus})
				
				// Emit lobby ready event (waiting for host to start)
				io.To(socket.Room(roomName(quiz.ID))).Emit("lobby_ready", gin.H{
					"quiz_id": quiz.ID,
					"message": "Lobby penuh, menunggu host memulai quiz",
					"players": len(quiz.Participants),
				})
				
				log.Printf("[socket.io] Lobby ready for quiz_id=%d, waiting for host to start", quiz.ID)
				return
			}
			// If only one participant, don't emit lobby_joined (wait for opponent)
			log.Printf("[socket.io] User %d joined lobby for quiz %d, still waiting for opponent", userID, quiz.ID)
		})

		// host_start_quiz - Host manually starts the multiplayer quiz
		client.On(EventHostStartQuiz, func(args ...any) {
			log.Printf("[socket.io] host_start_quiz received with args: %#v", args)
			
			var quizID uint
			var hostUserID uint
			
			// Parse arguments
			if len(args) == 1 {
				// Support object format: {quiz_id, host_user_id} or {quiz_id, user_id}
				if obj, ok := args[0].(map[string]any); ok {
					if qid, ok := obj["quiz_id"].(float64); ok {
						quizID = uint(qid)
					}
					// Support both host_user_id and user_id
					if uid, ok := obj["host_user_id"].(float64); ok {
						hostUserID = uint(uid)
					} else if uid, ok := obj["user_id"].(float64); ok {
						hostUserID = uint(uid)
					}
				}
			} else if len(args) >= 2 {
				// Support positional args: quiz_id, user_id
				if qid, ok := args[0].(float64); ok {
					quizID = uint(qid)
				}
				if uid, ok := args[1].(float64); ok {
					hostUserID = uint(uid)
				}
			}
			
			if quizID == 0 {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_quiz_id", "message": "Quiz ID is required"})
				return
			}
			
			if hostUserID == 0 {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_user_id", "message": "Host user ID is required"})
				return
			}
			
			// Get quiz with retry logic to handle database consistency
			var quiz *models.Quiz
			var err error
			maxRetries := 3
			for i := 0; i < maxRetries; i++ {
				quiz, err = quizSvc.GetQuizByID(quizID)
				if err != nil || quiz == nil {
					if i < maxRetries-1 {
						log.Printf("[socket.io] Quiz not found, retry %d/%d for quiz_id=%d", i+1, maxRetries, quizID)
						time.Sleep(time.Millisecond * 100) // Small delay for database consistency
						continue
					}
					client.Emit(EventLobbyError, gin.H{"error": "quiz_not_found", "message": "Quiz not found"})
					return
				}
				
				// Check if participants are loaded - if not, wait a bit and retry
				if len(quiz.Participants) == 0 && i < maxRetries-1 {
					log.Printf("[socket.io] No participants found yet, retry %d/%d for quiz_id=%d", i+1, maxRetries, quizID)
					time.Sleep(time.Millisecond * 200)
					continue
				}
				
				break // Success or final attempt
			}
			
			// Verify user is the host (first participant)
			log.Printf("[socket.io] Host validation - quiz.Participants count: %d", len(quiz.Participants))
			if len(quiz.Participants) > 0 {
				log.Printf("[socket.io] Host validation - first participant UserID: %d, provided hostUserID: %d", quiz.Participants[0].UserID, hostUserID)
				// Log all participants for debugging
				for i, p := range quiz.Participants {
					log.Printf("[socket.io] Participant[%d]: UserID=%d", i, p.UserID)
				}
			}
			
			if len(quiz.Participants) == 0 {
				client.Emit(EventLobbyError, gin.H{"error": "no_participants", "message": "No participants found in quiz"})
				return
			}
			
			if quiz.Participants[0].UserID != hostUserID {
				// Check if the user is in the participants list at all
				isParticipant := false
				for _, p := range quiz.Participants {
					if p.UserID == hostUserID {
						isParticipant = true
						break
					}
				}
				
				if isParticipant {
					client.Emit(EventLobbyError, gin.H{"error": "not_host", "message": "User is participant but not the host"})
				} else {
					client.Emit(EventLobbyError, gin.H{"error": "not_participant", "message": "User is not a participant in this quiz"})
				}
				return
			}
			
			// Start the quiz
			if quiz.ModuleID != nil {
				log.Printf("[socket.io] Host starting quiz for module_id=%d", *quiz.ModuleID)

				// Use mode-based question selection
				qs, err := questionSvc.GetQuestionsByModuleAndMode(*quiz.ModuleID, quiz.Mode)
				if err != nil {
					log.Printf("[socket.io] ERROR: Failed to get questions for module_id=%d mode=%s: %v", *quiz.ModuleID, quiz.Mode, err)
					// Fallback to old method
					qc := quiz.QuestionCount
					if qc <= 0 {
						qc = 10
					}
					qs, err = questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
					if err != nil {
						log.Printf("[socket.io] ERROR: Fallback also failed for module_id=%d: %v", *quiz.ModuleID, err)
						io.To(socket.Room(roomName(quiz.ID))).Emit(EventLobbyError, gin.H{
							"error":   "no_questions",
							"message": fmt.Sprintf("No questions available for module %d", *quiz.ModuleID),
						})
						return
					}
				}
				if qs == nil || len(*qs) == 0 {
					log.Printf("[socket.io] ERROR: No questions returned for module_id=%d mode=%s", *quiz.ModuleID, quiz.Mode)
					io.To(socket.Room(roomName(quiz.ID))).Emit(EventLobbyError, gin.H{
						"error":   "no_questions",
						"message": fmt.Sprintf("No questions available for module %d", *quiz.ModuleID),
					})
					return
				}

				log.Printf("[socket.io] Successfully retrieved %d questions for module_id=%d mode=%s", len(*qs), *quiz.ModuleID, quiz.Mode)

				s := getOrCreateSession(quiz)
				s.Questions = *qs
				if s.State != "running" { // prevent double start
					log.Printf("[socket.io] Emitting start_quiz event for quiz_id=%d with %d questions", quiz.ID, len(*qs))
					io.To(socket.Room(roomName(quiz.ID))).Emit(EventStartQuiz, gin.H{
						"quiz_id":         quiz.ID,
						"total_questions": len(*qs),
						"questions":       sanitizeQuestions(*qs),
						"module_id":       *quiz.ModuleID,
						"mode":            quiz.Mode,
					})
					log.Printf("[socket.io] Starting quiz session runner for quiz_id=%d", quiz.ID)
					go s.run(io)
				} else {
					log.Printf("[socket.io] Quiz session already running for quiz_id=%d, state=%s", quiz.ID, s.State)
				}
			} else {
				log.Printf("[socket.io] ERROR: Quiz has no module_id, cannot start quiz for quiz_id=%d", quiz.ID)
				io.To(socket.Room(roomName(quiz.ID))).Emit(EventLobbyError, gin.H{
					"error":   "no_module",
					"message": "Quiz has no module assigned",
				})
			}
		})

		// join_quiz
		client.On(EventJoinQuiz, func(args ...any) {
			if len(args) < 2 {
				return
			}
			qidf, ok1 := args[0].(float64)
			uidf, ok2 := args[1].(float64)
			if !ok1 || !ok2 {
				return
			}
			quizID := uint(qidf)
			userIDUint = uint(uidf)
			quiz, err := quizSvc.GetQuizByID(quizID)
			if err != nil || quiz == nil {
				return
			}
			room := roomName(quizID)
			client.Join(socket.Room(room))
			io.To(socket.Room(room)).Emit(EventUserJoin, gin.H{"user_id": userIDUint})
			if quiz.Mode == "multiplayer" {
				sessions.Lock()
				sess, exists := sessions.M[quizID]
				if !exists {
					log.Printf("[socket.io] join_quiz: creating new multiplayer session for quiz_id=%d", quizID)
					sessions.M[quizID] = NewMultiplayerSession(quiz, room, questionSvc, quizSvc, participantSvc, userSvc, nil)
					sess = sessions.M[quizID]
				}

				// Make sure the session is valid
				if sess == nil {
					sess = sessions.M[quizID]
					log.Printf("[socket.io] join_quiz: session was nil, retrieving from map for quiz_id=%d", quizID)
				}

				// Initialize players map if nil
				if sess.Players == nil {
					sess.Players = make(map[uint]*socket.Socket)
					log.Printf("[socket.io] join_quiz: players map was nil, initializing for quiz_id=%d", quizID)
				}

				// Add player to session
				sess.Players[userIDUint] = client
				players := len(sess.Players)
				sessions.Unlock()

				log.Printf("[socket.io] join_quiz: player added to multiplayer session, quiz_id=%d, user_id=%d, total_players=%d",
					quizID, userIDUint, players)

				if players == 2 && quiz.ModuleID != nil {
					log.Printf("[socket.io] join_quiz: Auto-starting quiz for module_id=%d with 2 players", *quiz.ModuleID)

					// Use mode-based question selection
					qs, err := questionSvc.GetQuestionsByModuleAndMode(*quiz.ModuleID, quiz.Mode)
					if err != nil {
						log.Printf("[socket.io] join_quiz: Failed to get questions by mode, falling back to random: %v", err)
						qc := quiz.QuestionCount
						if qc <= 0 {
							qc = 10
						}
						qs, err = questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
					}

					if err == nil && qs != nil && len(*qs) > 0 {
						log.Printf("[socket.io] join_quiz: Successfully retrieved %d questions for module_id=%d mode=%s", len(*qs), *quiz.ModuleID, quiz.Mode)
						sessions.Lock()
						if sessions.M[quizID] != nil {
							sessions.M[quizID].Questions = *qs
							sessions.Unlock()
							log.Printf("[socket.io] join_quiz: Emitting start_quiz event for quiz_id=%d", quiz.ID)
							io.To(socket.Room(room)).Emit(EventStartQuiz, gin.H{
								"quiz_id":         quiz.ID,
								"total_questions": len(*qs),
								"questions":       sanitizeQuestions(*qs),
								"module_id":       *quiz.ModuleID,
								"mode":            quiz.Mode,
							})
							log.Printf("[socket.io] join_quiz: Starting quiz session runner for quiz_id=%d", quiz.ID)
							go sessions.M[quizID].run(io)
						} else {
							sessions.Unlock()
							log.Printf("[socket.io] ERROR: session disappeared for quiz_id=%d", quizID)
						}
					} else {
						log.Printf("[socket.io] ERROR: failed to get questions for module_id=%d mode=%s: %v, questions count: %d", *quiz.ModuleID, quiz.Mode, err, func() int {
							if qs == nil {
								return 0
							}
							return len(*qs)
						}())
					}
				} else if players == 2 && quiz.ModuleID == nil {
					log.Printf("[socket.io] ERROR: join_quiz: Quiz has no module_id, cannot start quiz for quiz_id=%d", quizID)
				} else {
					log.Printf("[socket.io] join_quiz: Not starting quiz yet, players=%d, module_id present=%t", players, quiz.ModuleID != nil)
				}
				return
			}
			if quiz.ModuleID == nil {
				return
			}

			// Use mode-based question selection for single player
			qs, err := questionSvc.GetQuestionsByModuleAndMode(*quiz.ModuleID, quiz.Mode)
			if err != nil {
				log.Printf("[socket.io] Single player: Failed to get questions by mode, falling back to random: %v", err)
				qc := quiz.QuestionCount
				if qc <= 0 {
					qc = 10
				}
				qs, err = questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
			}

			if err != nil || qs == nil || len(*qs) == 0 {
				return
			}
			// Precompute boosters for determinism between FE/BE
			boosters := genBoosters(len(*qs))
			singleSess = NewSocketQuizSession(client, quiz, *qs, questionSvc, quizSvc, userSvc)
			singleSess.QuestionBoosters = make(map[uint]int)
			for i, qq := range *qs {
				idx := i
				if idx < len(boosters) {
					singleSess.QuestionBoosters[qq.ID] = boosters[idx]
				}
			}
			client.Emit(EventStartQuiz, gin.H{
				"quiz_id":         quiz.ID,
				"total_questions": len(singleSess.Questions),
				"questions":       sanitizeQuestionsWithBoosters(singleSess.Questions, boosters),
				"mode":            quiz.Mode,
			})
			singleSess.emitQuestion()
		})

		// join_quiz_by_module (single-player by module_id)
		client.On(EventJoinQuizByModule, func(args ...any) {
			// Supports: ({ module_id, user_id, question_count? }) | (module_id, user_id[, question_count]) | JSON string
			var moduleID uint
			var userID uint
			questionCount := 10
			parsed := false

			if len(args) == 1 {
				switch v := args[0].(type) {
				case map[string]any:
					if midf, ok := v["module_id"].(float64); ok { moduleID = uint(midf) }
					if uidf, ok := v["user_id"].(float64); ok { userID = uint(uidf) }
					if qc, ok := v["question_count"].(float64); ok && int(qc) > 0 { questionCount = int(qc) }
					parsed = moduleID != 0 && userID != 0
				case string:
					s := strings.TrimSpace(v)
					if strings.HasPrefix(s, "{") {
						var obj map[string]any
						if err := json.Unmarshal([]byte(s), &obj); err == nil {
							if midf, ok := obj["module_id"].(float64); ok { moduleID = uint(midf) }
							if uidf, ok := obj["user_id"].(float64); ok { userID = uint(uidf) }
							if qc, ok := obj["question_count"].(float64); ok && int(qc) > 0 { questionCount = int(qc) }
							parsed = moduleID != 0 && userID != 0
						}
					}
				}
			} else if len(args) >= 2 {
				if midf, ok := args[0].(float64); ok { moduleID = uint(midf) }
				if uidf, ok := args[1].(float64); ok { userID = uint(uidf) }
				if len(args) >= 3 {
					if qc, ok := args[2].(float64); ok && int(qc) > 0 { questionCount = int(qc) }
				}
				parsed = moduleID != 0 && userID != 0
			}

			if !parsed {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_args", "args": args})
				return
			}

			// Clean up any existing multiplayer sessions for this user
			sessions.Lock()
			for quizID, mpSession := range sessions.M {
				if mpSession != nil && mpSession.Players != nil {
					if _, exists := mpSession.Players[userID]; exists {
						log.Printf("[socket.io] Removing user %d from multiplayer session %d due to single-player join", userID, quizID)
						
						// Remove user from multiplayer session
						delete(mpSession.Players, userID)
						
						// Leave the multiplayer room
						client.Leave(socket.Room(roomName(quizID)))
						
						// If this was the last player, clean up the entire session
						if len(mpSession.Players) == 0 {
							log.Printf("[socket.io] Cleaning up empty multiplayer session %d", quizID)
							delete(sessions.M, quizID)
						} else {
							// Notify remaining players that opponent left
							mpSession.broadcast(io, EventUserLeave, gin.H{"user_id": userID, "message": "Opponent left for single-player mode"})
						}
						break // User can only be in one multiplayer session at a time
					}
				}
			}
			sessions.Unlock()

			// Prefer mode-specific question API
			var qs *[]models.Question
			var err error
			qs, err = questionSvc.GetQuestionsByModuleAndMode(moduleID, "single_player")
			if err != nil || qs == nil || len(*qs) == 0 {
				qs, err = questionSvc.GetRandomQuestionsByModule(moduleID, questionCount)
				if err != nil || qs == nil || len(*qs) == 0 {
					client.Emit(EventLobbyError, gin.H{"error": "no_questions", "module_id": moduleID})
					return
				}
			}

			// Create synthetic quiz metadata for session (not persisted)
			syntheticQuiz := &models.Quiz{Mode: "single_player"}
			syntheticQuiz.ModuleID = new(uint)
			*syntheticQuiz.ModuleID = moduleID

			// Precompute boosters for determinism
			boosters := genBoosters(len(*qs))
			singleSess = NewSocketQuizSession(client, syntheticQuiz, *qs, questionSvc, quizSvc, userSvc)
			singleSess.QuestionBoosters = make(map[uint]int)
			for i, qq := range *qs {
				if i < len(boosters) {
					singleSess.QuestionBoosters[qq.ID] = boosters[i]
				}
			}
			client.Emit(EventStartQuiz, gin.H{
				"quiz_id":         0, // no persisted quiz id
				"module_id":       moduleID,
				"total_questions": len(singleSess.Questions),
				"questions":       sanitizeQuestionsWithBoosters(singleSess.Questions, boosters),
				"mode":            "single_player",
			})
			singleSess.emitQuestion()
		})

		// submit_answer (supports either positional args or a single JSON object)
		client.On(EventAnswerSubmit, func(args ...any) {
			log.Printf("[socket.io] submit_answer raw=%#v", args)
			var quizID, questionID, userID uint
			var option string
			var selectedIndex *int
			parsed := false

			if len(args) == 1 {
				switch v := args[0].(type) {
				case map[string]any:
					if qidf, ok := v["quiz_id"].(float64); ok {
						quizID = uint(qidf)
					}
					if qnidf, ok := v["question_id"].(float64); ok {
						questionID = uint(qnidf)
					}
					if opt, ok := v["option"].(string); ok {
						option = opt
					}
					// Also allow nested options.selected_index or top-level selected_index
					if idxf, ok := v["selected_index"].(float64); ok {
						i := int(idxf)
						selectedIndex = &i
					}
					if optObj, ok := v["options"].(map[string]any); ok {
						if idxf, ok := optObj["selected_index"].(float64); ok {
							i := int(idxf)
							selectedIndex = &i
						}
					}
					if uidf, ok := v["user_id"].(float64); ok {
						userID = uint(uidf)
					}
					// Allow quiz_id == 0 for module mode
					parsed = questionID != 0 && (option != "" || selectedIndex != nil) && userID != 0
				case string:
					s := strings.TrimSpace(v)
					if strings.HasPrefix(s, "{") {
						var obj map[string]any
						if err := json.Unmarshal([]byte(s), &obj); err == nil {
							if qidf, ok := obj["quiz_id"].(float64); ok {
								quizID = uint(qidf)
							}
							if qnidf, ok := obj["question_id"].(float64); ok {
								questionID = uint(qnidf)
							}
							if opt, ok := obj["option"].(string); ok {
								option = opt
							}
							if idxf, ok := obj["selected_index"].(float64); ok {
								i := int(idxf)
								selectedIndex = &i
							}
							if optObj, ok := obj["options"].(map[string]any); ok {
								if idxf, ok := optObj["selected_index"].(float64); ok {
									i := int(idxf)
									selectedIndex = &i
								}
							}
							if uidf, ok := obj["user_id"].(float64); ok {
								userID = uint(uidf)
							}
							// Allow quiz_id == 0 for module mode
							parsed = questionID != 0 && (option != "" || selectedIndex != nil) && userID != 0
						}
					}
				}
			}
			if !parsed {
				if len(args) >= 4 {
					qidf, okQ := args[0].(float64)
					qnidf, okQN := args[1].(float64)
					// third arg can be option string or selected index (number)
					opt, okOpt := args[2].(string)
					if !okOpt {
						if idxf, ok := args[2].(float64); ok {
							i := int(idxf)
							selectedIndex = &i
							okOpt = true
						}
					}
					uidf, okUID := args[3].(float64)
					if okQ && okQN && okOpt && okUID {
						quizID = uint(qidf)
						questionID = uint(qnidf)
						option = opt
						userID = uint(uidf)
						// Allow quiz_id == 0 (module mode)
						parsed = questionID != 0 && (option != "" || selectedIndex != nil) && userID != 0
					}
				}
			}
			if !parsed {
				log.Printf("[socket.io] submit_answer: invalid args: %#v", args)
				return
			}

			// If only index is provided, map to option key
			if option == "" && selectedIndex != nil {
				norm := normalizeIndex(*selectedIndex)
				if norm >= 0 {
					option = optionKey(norm)
				}
			}

			// Check for multiplayer session
			sessions.Lock()
			mp, okSess := sessions.M[quizID]
			sessions.Unlock()

			if okSess && mp != nil {
				// Ensure AnswerCh exists
				if mp.AnswerCh == nil {
					mp.AnswerCh = make(chan AnswerSubmission, 10)
				}
				// Auto-start session if not running and we have/ can load questions
				if mp.State != "running" {
					// Load questions if missing
					if len(mp.Questions) == 0 {
						if quiz, err := quizSvc.GetQuizByID(quizID); err == nil && quiz != nil && quiz.ModuleID != nil {
							qc := quiz.QuestionCount
							if qc <= 0 {
								qc = 10
							}
							if qs, err := questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc); err == nil && qs != nil && len(*qs) > 0 {
								mp.Questions = *qs
							}
						}
					}
					if len(mp.Questions) > 0 {
						io.To(socket.Room(roomName(quizID))).Emit(EventStartQuiz, gin.H{"quiz_id": quizID, "total_questions": len(mp.Questions), "questions": sanitizeQuestions(mp.Questions)})
						go mp.run(io)
					} else {
						log.Printf("[socket.io] cannot auto-start: no questions for quiz_id=%d", quizID)
					}
				}

				log.Printf("[socket.io] submitting answer for multiplayer quiz_id=%d, user_id=%d", quizID, userID)
				mp.AnswerCh <- AnswerSubmission{UserID: userID, QuestionID: questionID, Option: option}

				return
			}

			// Check for single player session
			if singleSess != nil {
				log.Printf("[socket.io] submitting answer for single player quiz_id=%d, user_id=%d", quizID, userID)
				singleSess.handleAnswer(userID, option)
			} else {
				log.Printf("[socket.io] no valid session found for quiz_id=%d", quizID)
			}
		})

		// opponent_taunt - Send taunt/reaction to opponent in multiplayer
		client.On(EventOpponentTaunt, func(args ...any) {
			log.Printf("[socket.io] opponent_taunt received with args: %#v", args)
			
			var quizID uint
			var senderUserID uint
			var linkLottie string
			var tauntType string
			parsed := false
			
			// Parse arguments - expect {quiz_id, user_id, link_lottie, type?}
			if len(args) == 1 {
				if obj, ok := args[0].(map[string]any); ok {
					if qid, ok := obj["quiz_id"].(float64); ok {
						quizID = uint(qid)
					}
					if uid, ok := obj["user_id"].(float64); ok {
						senderUserID = uint(uid)
					}
					if link, ok := obj["link_lottie"].(string); ok {
						linkLottie = link
					}
					if tType, ok := obj["type"].(string); ok {
						tauntType = tType
					} else {
						tauntType = "reaction" // default type
					}
					parsed = quizID != 0 && senderUserID != 0 && linkLottie != ""
				}
			}
			
			if !parsed {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_taunt_args", "args": args})
				return
			}
			
			// Find the multiplayer session
			sessions.Lock()
			mpSession, exists := sessions.M[quizID]
			sessions.Unlock()
			
			if !exists || mpSession == nil {
				client.Emit(EventLobbyError, gin.H{"error": "session_not_found", "quiz_id": quizID})
				return
			}
			
			// Find the opponent (the other player in the session)
			var opponentSocket *socket.Socket
			var opponentID uint
			
			for uid, sock := range mpSession.Players {
				if uid != senderUserID && sock != nil {
					opponentSocket = sock
					opponentID = uid
					break
				}
			}
			
			if opponentSocket == nil {
				client.Emit(EventLobbyError, gin.H{"error": "opponent_not_found"})
				return
			}
			
			// Get sender user info for the taunt
			var senderName string
			if userSvc != nil {
				if user, err := userSvc.GetUserByIDUint(senderUserID); err == nil && user != nil {
					senderName = user.Name
				}
			}
			if senderName == "" {
				senderName = "Opponent" // fallback
			}
			
			log.Printf("[socket.io] Sending taunt from user %d (%s) to user %d", senderUserID, senderName, opponentID)
			
			// Send the taunt to the opponent
			opponentSocket.Emit(EventTauntReceived, gin.H{
				"quiz_id":     quizID,
				"sender_id":   senderUserID,
				"sender_name": senderName,
				"link_lottie": linkLottie,
				"type":        tauntType,
				"timestamp":   time.Now().Unix(),
			})
			
			// Confirm to sender that taunt was sent
			client.Emit("taunt_sent", gin.H{
				"quiz_id":     quizID,
				"target_id":   opponentID,
				"link_lottie": linkLottie,
				"type":        tauntType,
			})
		})

		client.On("disconnect", func(...any) {
			log.Printf("[socket.io] disconnected id=%s", client.Id())
			sessions.Lock()
			for _, s := range sessions.M {
				for uid, sock := range s.Players {
					if sock == client {
						delete(s.Players, uid)
						io.To(socket.Room(s.RoomName)).Emit(EventUserLeave, gin.H{"user_id": uid})
					}
				}
			}
			sessions.Unlock()
		})
	})

	router.GET("/socket.io/*any", gin.WrapH(io.ServeHandler(opts)))
	router.POST("/socket.io/*any", gin.WrapH(io.ServeHandler(opts)))
	router.OPTIONS("/socket.io/*any", gin.WrapH(io.ServeHandler(opts)))
	router.GET("/socket-io-health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
}

func roomName(id uint) string { return "quiz-" + fmt.Sprintf("%d", id) }

// optionIndex returns 0-based index for a/b/c/d, or -1 if unknown
func optionIndex(key string) int {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "a":
		return 0
	case "b":
		return 1
	case "c":
		return 2
	case "d":
		return 3
	default:
		return -1
	}
}

// getOptionText returns the text value of an option based on the key (A/B/C/D)
func getOptionText(q models.Question, optionKey string) string {
	switch strings.ToUpper(strings.TrimSpace(optionKey)) {
	case "A":
		return q.Options.OptionA
	case "B":
		return q.Options.OptionB
	case "C":
		return q.Options.OptionC
	case "D":
		return q.Options.OptionD
	default:
		return ""
	}
}

// optionKey returns canonical option key (a/b/c/d) for an index, or empty string if out of range
func optionKey(idx int) string {
	switch idx {
	case 0:
		return "a"
	case 1:
		return "b"
	case 2:
		return "c"
	case 3:
		return "d"
	default:
		return ""
	}
}

// normalizeIndex accepts either 0-based [0..3] or 1-based [1..4] index and returns 0-based, otherwise -1
func normalizeIndex(idx int) int {
	if idx >= 0 && idx <= 3 {
		return idx
	}
	if idx >= 1 && idx <= 4 {
		return idx - 1
	}
	return -1
}

// sanitizeQuestions returns a client-safe list of questions without revealing correct answers
func sanitizeQuestions(qs []models.Question) []map[string]any {
	out := make([]map[string]any, 0, len(qs))
	boosters := []int{2, 3, 5}
	for _, q := range qs {
		booster := 0
		if rand.Intn(10) < 3 { // ~30% chance to have a booster
			booster = boosters[rand.Intn(len(boosters))]
		}
		out = append(out, map[string]any{
			"id":          q.ID,
			"question_id": q.ID,
			"question":    q.Question,
			"options": map[string]any{
				"a": q.Options.OptionA,
				"b": q.Options.OptionB,
				"c": q.Options.OptionC,
				"d": q.Options.OptionD,
			},
			"read_time":      q.ReadTime,
			"answer_time":    q.AnswerTime,
			"module_id":      q.ModuleID,
			"booster":        booster,
			"question_type":  q.QuestionType, // "hots" or "regular"
			// Expose correct answer as requested by FE
			"correct_index":  optionIndex(q.CorrectAnswer),
			"correct_option": strings.ToLower(strings.TrimSpace(q.CorrectAnswer)),
			"explanation":    q.Explanation,
		})
	}
	return out
}

// genBoosters creates a deterministic booster slice for N questions.
// We base it on a simple repeating pattern for predictability: [0,2,0,3,0,5] ...
func genBoosters(n int) []int {
	pattern := []int{0, 2, 0, 3, 0, 5}
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = pattern[i%len(pattern)]
	}
	return out
}

// sanitizeQuestionsWithBoosters uses provided boosters instead of random assignment
func sanitizeQuestionsWithBoosters(qs []models.Question, boosters []int) []map[string]any {
	out := make([]map[string]any, 0, len(qs))
	for i, q := range qs {
		booster := 0
		if i < len(boosters) {
			booster = boosters[i]
		}
		out = append(out, map[string]any{
			"id":          q.ID,
			"question_id": q.ID,
			"question":    q.Question,
			"options": map[string]any{
				"a": q.Options.OptionA,
				"b": q.Options.OptionB,
				"c": q.Options.OptionC,
				"d": q.Options.OptionD,
			},
			"read_time":      q.ReadTime,
			"answer_time":    q.AnswerTime,
			"module_id":      q.ModuleID,
			"booster":        booster,
			"question_type":  q.QuestionType, // "hots" or "regular"
			"correct_index":  optionIndex(q.CorrectAnswer),
			"correct_option": strings.ToLower(strings.TrimSpace(q.CorrectAnswer)),
			"explanation":    q.Explanation,
		})
	}
	return out
}
