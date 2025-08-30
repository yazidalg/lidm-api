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
	EventJoinQuiz       = "join_quiz"
	EventStartQuiz      = "start_quiz"
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
)

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
}

func NewSocketQuizSession(sock *socket.Socket, quiz *models.Quiz, questions []models.Question, qSvc services.QuestionServiceInterface, quizSvc services.QuizServiceInterface, userSvc services.UserServiceInterface) *SocketQuizSession {
	return &SocketQuizSession{Socket: sock, Quiz: quiz, Questions: questions, Answered: make(map[uint]bool), QuestionSvc: qSvc, QuizSvc: quizSvc, UserSvc: userSvc}
}

func (s *SocketQuizSession) emitQuestion() {
	if s.CurrentIdx >= len(s.Questions) {
		s.Socket.Emit(EventQuizCompleted, gin.H{"message": "Quiz finished"})
		return
	}
	q := s.Questions[s.CurrentIdx]
	// Do not expose full question (avoids leaking correct_answer). Frontend should use questions from start_quiz.
	s.Socket.Emit(EventQuestion, gin.H{"question_id": q.ID, "index": s.CurrentIdx})
}

func (s *SocketQuizSession) handleAnswer(userID uint, option string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if s.CurrentIdx >= len(s.Questions) {
		return
	}
	q := s.Questions[s.CurrentIdx]
	if s.Answered[q.ID] {
		return
	}
	// Compare by index to be case-insensitive and robust
	correct := optionIndex(option) == optionIndex(q.CorrectAnswer)
	if correct {
		_ = s.UserSvc.AddXP(userID, 10)
	} else {
		quiz, _ := s.QuizSvc.GetQuizByID(s.Quiz.ID)
		if quiz != nil && quiz.Mode == "single_player" {
			_ = s.UserSvc.DecrementLife(userID)
			if u, err := s.UserSvc.GetUserByIDUint(userID); err == nil && u != nil && u.Lives <= 0 {
				s.Socket.Emit(EventLivesExhausted, gin.H{"message": "Lives exhausted", "lives": u.Lives})
			}
		}
	}
	// Emit only the selected and correct option indexes for highlighting
	s.Socket.Emit(EventAnswerResult, gin.H{
		"question_id": q.ID,
		"is_correct":  correct,
		"options": gin.H{
			"selected_index": optionIndex(option),
			"correct_index":  optionIndex(q.CorrectAnswer),
		},
	})
	s.Answered[q.ID] = true
	// no question_ended event
	s.CurrentIdx++
	go func() { time.Sleep(time.Second); s.emitQuestion() }()
}

// Multiplayer session (simplified)
type MultiplayerSession struct {
	QuizID         uint
	RoomName       string
	Quiz           *models.Quiz
	Questions      []models.Question
	CurrentIdx     int
	State          string
	Players        map[uint]*socket.Socket
	Scores         map[uint]int
	Answers        map[uint]bool
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
	return &MultiplayerSession{QuizID: quiz.ID, RoomName: room, Quiz: quiz, Questions: questions, Players: make(map[uint]*socket.Socket), Scores: make(map[uint]int), Answers: make(map[uint]bool), QuestionSvc: qSvc, QuizSvc: quizSvc, ParticipantSvc: pSvc, UserSvc: uSvc, AnswerCh: make(chan AnswerSubmission, 10), TimerCancel: make(chan struct{}, 1)}
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

	// Compare by index to avoid case/key mismatches
	correct := optionIndex(ans.Option) == optionIndex(q.CorrectAnswer)
	if correct {
		s.Scores[ans.UserID] += 10
		if s.UserSvc != nil {
			_ = s.UserSvc.AddXP(ans.UserID, 10)
		} else {
			log.Printf("[socket.io] WARNING: UserService is nil, can't add XP")
		}
		log.Printf("[socket.io] Correct answer from user %d, new score: %d", ans.UserID, s.Scores[ans.UserID])
	} else {
		log.Printf("[socket.io] Incorrect answer from user %d, answer was: %s, correct: %s",
			ans.UserID, ans.Option, q.CorrectAnswer)
	}

	s.Answers[ans.UserID] = true
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
			oppScore := 0
			var oppID uint
			for id, sc := range s.Scores {
				if id != uid {
					oppScore = sc
					oppID = id
				}
			}
			sock.Emit(EventNextQuestion, gin.H{
				"current":  gin.H{"question_id": q.ID},
				"opponent": gin.H{"user_id": oppID, "score": oppScore},
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
		s.broadcast(io, EventQuizCompleted, gin.H{"scores": s.Scores, "winner": winnerName})
		return
	}

	quiz, err := s.QuizSvc.GetQuizByID(s.QuizID)
	if err != nil || quiz == nil {
		log.Printf("[socket.io] ERROR: Failed to get quiz %d: %v", s.QuizID, err)
		s.broadcast(io, EventQuizCompleted, gin.H{"scores": s.Scores, "winner": winnerName})
		return
	}

	// Check if ParticipantService is available
	if s.ParticipantSvc == nil {
		log.Printf("[socket.io] ERROR: ParticipantService is nil in finish()")
		s.broadcast(io, EventQuizCompleted, gin.H{"scores": s.Scores, "winner": winnerName})
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
	s.broadcast(io, EventQuizCompleted, gin.H{"scores": s.Scores, "winner": winnerName})
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

			// Immediately emit lobby_joined for the host (creator is considered joined)
			if u, err := userSvc.GetUserByIDUint(req.HostUserID); err == nil && u != nil {
				client.Emit(EventLobbyJoined, gin.H{
					"quiz_id":  quiz.ID,
					"user":     map[string]any{"id": u.ID, "name": u.Name, "point": u.Point, "profile_picture": u.ProfilePicture},
					"opponent": nil,
				})
			} else {
				client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user": map[string]any{"id": req.HostUserID}, "opponent": nil})
			}
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
						me = map[string]any{"id": u.ID, "name": u.Name, "point": u.Point, "profile_picture": u.ProfilePicture}
					}
					if opponent != nil {
						opp = map[string]any{"id": opponent.ID, "name": opponent.Name, "point": opponent.Point, "profile_picture": opponent.ProfilePicture}
					}
					client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user": me, "opponent": opp})
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
				// Emit to participant[0]
				if u1 != nil {
					if s1, ok := sess.Players[u1.ID]; ok && s1 != nil {
						s1.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID,
							"user": map[string]any{"id": u1.ID, "name": u1.Name, "point": u1.Point, "profile_picture": u1.ProfilePicture},
							"opponent": func() any {
								if u2 == nil {
									return nil
								}
								return map[string]any{"id": u2.ID, "name": u2.Name, "point": u2.Point, "profile_picture": u2.ProfilePicture}
							}(),
						})
					}
				}
				// Emit to participant[1]
				if u2 != nil {
					if s2, ok := sess.Players[u2.ID]; ok && s2 != nil {
						s2.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID,
							"user": map[string]any{"id": u2.ID, "name": u2.Name, "point": u2.Point, "profile_picture": u2.ProfilePicture},
							"opponent": func() any {
								if u1 == nil {
									return nil
								}
								return map[string]any{"id": u1.ID, "name": u1.Name, "point": u1.Point, "profile_picture": u1.ProfilePicture}
							}(),
						})
					}
				}
				// Emit random XP bonus to both
				bonus := []int{2, 3, 5}[time.Now().UnixNano()%3]
				io.To(socket.Room(roomName(quiz.ID))).Emit("bonus_multiplier", gin.H{"multiplier": bonus})
				// Start quiz automatically
				if quiz.ModuleID != nil {
					log.Printf("[socket.io] Auto-starting quiz for module_id=%d", *quiz.ModuleID)

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
							// Emit error to both players
							io.To(socket.Room(roomName(quiz.ID))).Emit(EventLobbyError, gin.H{
								"error":   "no_questions",
								"message": fmt.Sprintf("No questions available for module %d", *quiz.ModuleID),
							})
							return
						}
					}
					if qs == nil || len(*qs) == 0 {
						log.Printf("[socket.io] ERROR: No questions returned for module_id=%d mode=%s", *quiz.ModuleID, quiz.Mode)
						// Emit error to both players
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
				return
			}
			// If only one participant, emit lobby_joined with opponent nil and user details
			if u, err := userSvc.GetUserByIDUint(userID); err == nil && u != nil {
				client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user": map[string]any{"id": u.ID, "name": u.Name, "point": u.Point, "profile_picture": u.ProfilePicture}, "opponent": nil})
			} else {
				client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user": map[string]any{"id": userID}, "opponent": nil})
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
			singleSess = NewSocketQuizSession(client, quiz, *qs, questionSvc, quizSvc, userSvc)
			client.Emit(EventStartQuiz, gin.H{
				"quiz_id":         quiz.ID,
				"total_questions": len(singleSess.Questions),
				"questions":       sanitizeQuestions(singleSess.Questions),
				"mode":            quiz.Mode,
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
					parsed = quizID != 0 && questionID != 0 && (option != "" || selectedIndex != nil) && userID != 0
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
							parsed = quizID != 0 && questionID != 0 && (option != "" || selectedIndex != nil) && userID != 0
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
						parsed = true
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
			"read_time":   q.ReadTime,
			"answer_time": q.AnswerTime,
			"module_id":   q.ModuleID,
			"booster":     booster,
			// Expose correct answer as requested by FE
			"correct_index":  optionIndex(q.CorrectAnswer),
			"correct_option": strings.ToLower(strings.TrimSpace(q.CorrectAnswer)),
			"explanation":    q.Explanation,
		})
	}
	return out
}
