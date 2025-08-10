package socketio

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	EventQuestionEnded  = "question_ended"
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
	s.Socket.Emit(EventQuestion, q)
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
	correct := option == q.CorrectAnswer
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
	s.Socket.Emit(EventAnswerResult, gin.H{"is_correct": correct})
	s.Answered[q.ID] = true
	s.Socket.Emit(EventQuestionEnded, gin.H{"correct_answer": q.CorrectAnswer})
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
	s.State = "running"
	for s.CurrentIdx < len(s.Questions) {
		s.runQuestion(io)
		s.CurrentIdx++
	}
	s.finish(io)
}
func (s *MultiplayerSession) runQuestion(io *socket.Server) {
	q := s.Questions[s.CurrentIdx]
	s.Answers = make(map[uint]bool)
	s.broadcast(io, EventQuestion, q)
	totalTime := time.Duration(q.AnswerTime+q.ReadTime) * time.Second
	if totalTime <= 0 {
		totalTime = 15 * time.Second
	}
	timer := time.NewTimer(totalTime)
	defer timer.Stop()
questionLoop:
	for {
		select {
		case ans := <-s.AnswerCh:
			if ans.QuestionID != q.ID {
				continue
			}
			s.processAnswer(io, q, ans)
			if len(s.Answers) == len(s.Players) {
				break questionLoop
			}
		case <-timer.C:
			break questionLoop
		case <-s.TimerCancel:
			break questionLoop
		}
	}
	s.broadcast(io, EventQuestionEnded, gin.H{"correct_answer": q.CorrectAnswer, "question_id": q.ID})
	time.Sleep(time.Second)
}
func (s *MultiplayerSession) processAnswer(io *socket.Server, q models.Question, ans AnswerSubmission) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if s.Answers[ans.UserID] {
		return
	}
	correct := ans.Option == q.CorrectAnswer
	if correct {
		s.Scores[ans.UserID] += 10
		_ = s.UserSvc.AddXP(ans.UserID, 10)
	}
	s.Answers[ans.UserID] = true
	io.To(socket.Room(s.RoomName)).Emit(EventAnswerResult, gin.H{"user_id": ans.UserID, "question_id": q.ID, "is_correct": correct, "score": s.Scores[ans.UserID]})
}
func (s *MultiplayerSession) finish(io *socket.Server) {
	s.State = "finished"
	winnerName := "Seri"
	var winnerID *uint
	maxScore := -1
	quiz, err := s.QuizSvc.GetQuizByID(s.QuizID)
	if err == nil && quiz != nil {
		for _, part := range quiz.Participants {
			score := s.Scores[part.UserID]
			updateReq := request.UpdateParticipantRequest{TotalScore: score}
			_, _ = s.ParticipantSvc.UpdateParticipant(int32(part.ID), updateReq)
			if score > maxScore {
				maxScore = score
				winnerName = part.User.Name
				tempID := part.UserID
				winnerID = &tempID
			} else if score == maxScore && maxScore != -1 {
				winnerName = "Seri"
				winnerID = nil
			}
		}
		updateQuizReq := request.UpdateQuizRequest{Status: "completed", WinnerID: winnerID}
		_, _ = s.QuizSvc.UpdateQuiz(s.QuizID, updateQuizReq)
	}
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
			client.Emit(EventLobbyCreated, gin.H{"quiz_id": quiz.ID, "invite_code": quiz.InviteCode, "module_id": quiz.ModuleID, "question_count": quiz.QuestionCount})
		})

		// join_lobby
		client.On(EventJoinLobby, func(args ...any) {
			var code, token string
			var okC, okT bool
			var userID uint
			if len(args) == 1 {
				switch v := args[0].(type) {
				case map[string]any:
					if vv, ok := v["invite_code"].(string); ok {
						code, okC = vv, true
					}
					if vv, ok := v["token"].(string); ok {
						token, okT = vv, true
					}
				case string:
					s := strings.TrimSpace(v)
					if strings.HasPrefix(s, "{") {
						var obj map[string]any
						if err := json.Unmarshal([]byte(s), &obj); err == nil {
							if vv, ok := obj["invite_code"].(string); ok {
								code, okC = vv, true
							}
							if vv, ok := obj["token"].(string); ok {
								token, okT = vv, true
							}
						}
					}
				}
			} else if len(args) >= 2 {
				code, okC = args[0].(string)
				token, okT = args[1].(string)
			}
			if !okC || !okT {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_args"})
				return
			}
			uid, errTok := utils.ParseToken(token)
			if errTok != nil || uid == 0 {
				client.Emit(EventLobbyError, gin.H{"error": "invalid_token"})
				return
			}
			userID = uid
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
			// Check if already joined
			for _, p := range quiz.Participants {
				if p.UserID == userID {
					// Find opponent info
					var opponent *models.User
					for _, op := range quiz.Participants {
						if op.UserID != userID {
							opponent = &op.User
						}
					}
					client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user_id": userID, "opponent": opponent})
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
			// Refresh quiz participants
			quiz, _ = quizSvc.GetQuizByInviteCode(code)
			if len(quiz.Participants) == 2 {
				_, _ = quizSvc.UpdateQuiz(quiz.ID, request.UpdateQuizRequest{Status: "in_progress"})
				// Emit lobby_joined to both clients with opponent info
				for _, p := range quiz.Participants {
					var opponent *models.User
					for _, op := range quiz.Participants {
						if op.UserID != p.UserID {
							opponent = &op.User
						}
					}
					io.To(socket.Room(roomName(quiz.ID))).Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user_id": p.UserID, "opponent": opponent})
				}
				// Emit random XP bonus to both
				bonus := []int{2, 3, 5}[time.Now().UnixNano()%3]
				io.To(socket.Room(roomName(quiz.ID))).Emit("bonus_multiplier", gin.H{"multiplier": bonus})
				// Start quiz automatically
				if quiz.ModuleID != nil {
					qc := quiz.QuestionCount
					if qc <= 0 {
						qc = 10
					}
					qs, err := questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
					if err == nil && qs != nil && len(*qs) > 0 {
						sessions.Lock()
						sessions.M[quiz.ID].Questions = *qs
						sessions.Unlock()
						io.To(socket.Room(roomName(quiz.ID))).Emit(EventStartQuiz, gin.H{"quiz_id": quiz.ID, "total_questions": len(*qs)})
						go sessions.M[quiz.ID].run(io)
					}
				}
				return
			}
			// If only one participant, emit lobby_joined with opponent nil
			client.Emit(EventLobbyJoined, gin.H{"quiz_id": quiz.ID, "user_id": userID, "opponent": nil})
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
					sessions.M[quizID] = NewMultiplayerSession(quiz, room, questionSvc, quizSvc, participantSvc, userSvc, nil)
					sess = sessions.M[quizID]
				}
				sess.Players[userIDUint] = client
				players := len(sess.Players)
				sessions.Unlock()
				if players == 2 && quiz.ModuleID != nil {
					qc := quiz.QuestionCount
					if qc <= 0 {
						qc = 10
					}
					qs, err := questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
					if err == nil && qs != nil && len(*qs) > 0 {
						sessions.Lock()
						sessions.M[quizID].Questions = *qs
						sessions.Unlock()
						io.To(socket.Room(room)).Emit(EventStartQuiz, gin.H{"quiz_id": quiz.ID, "total_questions": len(*qs)})
						go sessions.M[quizID].run(io)
					}
				}
				return
			}
			if quiz.ModuleID == nil {
				return
			}
			qc := quiz.QuestionCount
			if qc <= 0 {
				qc = 10
			}
			qs, err := questionSvc.GetRandomQuestionsByModule(*quiz.ModuleID, qc)
			if err != nil || qs == nil || len(*qs) == 0 {
				return
			}
			singleSess = NewSocketQuizSession(client, quiz, *qs, questionSvc, quizSvc, userSvc)
			client.Emit(EventStartQuiz, gin.H{"quiz_id": quiz.ID, "total_questions": len(singleSess.Questions)})
			singleSess.emitQuestion()
		})

		// submit_answer
		client.On(EventAnswerSubmit, func(args ...any) {
			if len(args) < 4 {
				return
			}
			qidf, okQ := args[0].(float64)
			qnidf, okQN := args[1].(float64)
			option, okOpt := args[2].(string)
			uidf, okUID := args[3].(float64)
			if !okQ || !okQN || !okOpt || !okUID {
				return
			}
			quizID := uint(qidf)
			questionID := uint(qnidf)
			sessions.Lock()
			mp, okSess := sessions.M[quizID]
			sessions.Unlock()
			if okSess && mp != nil && mp.State == "running" {
				mp.AnswerCh <- AnswerSubmission{UserID: uint(uidf), QuestionID: questionID, Option: option}
				return
			}
			if singleSess != nil {
				singleSess.handleAnswer(uint(uidf), option)
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
