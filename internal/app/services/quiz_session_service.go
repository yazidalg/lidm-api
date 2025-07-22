package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/response"
	"gorm.io/gorm"
)

type QuizSessionServiceInterface interface {
	CreateQuizSession(userID uint, req request.CreateQuizSessionRequest) (*response.QuizSessionResponse, error)
	JoinQuiz(userID uint, req request.JoinQuizRequest) (*response.JoinQuizResponse, error)
	AnswerQuestion(userID uint, req request.AnswerQuestionRequest) (*models.QuizSession, error)
	GetQuizSession(quizID uint) (*response.QuizSessionResponse, error)
	GetQuizResult(quizID uint) (*response.QuizResultResponse, error)
	FinishQuiz(quizID, userID uint) error
}

type quizSessionService struct {
	quizSessionRepo repositories.QuizSessionRepositoryInterface
	quizRepo        repositories.QuizRepositoryInterface
	questionRepo    repositories.QuestionRepositoryInterface
	participantRepo repositories.ParticipantRepositoryInterface
	moduleRepo      repositories.ModuleRepositoryInterface
}

func NewQuizSessionService(
	quizSessionRepo repositories.QuizSessionRepositoryInterface,
	quizRepo repositories.QuizRepositoryInterface,
	questionRepo repositories.QuestionRepositoryInterface,
	participantRepo repositories.ParticipantRepositoryInterface,
	moduleRepo repositories.ModuleRepositoryInterface,
) *quizSessionService {
	return &quizSessionService{
		quizSessionRepo: quizSessionRepo,
		quizRepo:        quizRepo,
		questionRepo:    questionRepo,
		participantRepo: participantRepo,
		moduleRepo:      moduleRepo,
	}
}

func (s *quizSessionService) CreateQuizSession(userID uint, req request.CreateQuizSessionRequest) (*response.QuizSessionResponse, error) {
	// Validate module exists
	_, err := s.moduleRepo.GetModuleByID(uint32(req.ModuleID))
	if err != nil {
		return nil, fmt.Errorf("module not found: %v", err)
	}

	// Get random questions from the module
	questionCount := 5 // Default
	questions, err := s.questionRepo.GetRandomQuestionsByModule(req.ModuleID, questionCount)
	if err != nil || len(*questions) == 0 {
		return nil, fmt.Errorf("no questions found for this module")
	}

	// Create quiz
	quiz := &models.Quiz{
		Status:        "pending",
		Mode:          req.Mode,
		ModuleID:      &req.ModuleID,
		HostUserID:    userID,
		QuestionCount: questionCount,
		Questions:     *questions,
	}

	createdQuiz, err := s.quizRepo.CreateQuiz(quiz)
	if err != nil {
		return nil, fmt.Errorf("failed to create quiz: %v", err)
	}

	// Create participant for the creator
	participant := &models.Participant{
		QuizID: createdQuiz.ID,
		UserID: userID,
	}

	_, err = s.participantRepo.CreateParticipant(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to create participant: %v", err)
	}

	// Build response
	resp := &response.QuizSessionResponse{
		QuizID:         createdQuiz.ID,
		QuestionNumber: 0,
		TotalQuestions: len(*questions),
		Phase:          "waiting", // Waiting for participants (multiplayer) or ready to start (single_player)
	}

	// Add participant status
	resp.Participants = []response.ParticipantStatus{
		{
			UserID:       userID,
			Username:     createdQuiz.Host.Name,
			IsReady:      true,
			CurrentScore: 0,
			IsFinished:   false,
		},
	}

	return resp, nil
}

func (s *quizSessionService) JoinQuiz(userID uint, req request.JoinQuizRequest) (*response.JoinQuizResponse, error) {
	// Find quiz by invite code
	quiz, err := s.quizRepo.GetQuizByInviteCode(req.InviteCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid invite code")
		}
		return nil, fmt.Errorf("failed to find quiz: %v", err)
	}

	// Check if quiz is still pending
	if quiz.Status != "pending" {
		return nil, fmt.Errorf("quiz has already started or finished")
	}

	// Check if user is not already a participant
	for _, participant := range quiz.Participants {
		if participant.UserID == userID {
			return nil, fmt.Errorf("user is already a participant")
		}
	}

	// Check room capacity for multiplayer
	if quiz.Mode == "multiplayer" && len(quiz.Participants) >= 2 {
		return nil, fmt.Errorf("quiz room is full")
	}

	// Create participant
	participant := &models.Participant{
		QuizID: quiz.ID,
		UserID: userID,
	}

	_, err = s.participantRepo.CreateParticipant(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to join quiz: %v", err)
	}

	// Build response
	response := &response.JoinQuizResponse{
		QuizID:       quiz.ID,
		InviteCode:   quiz.InviteCode,
		ModuleName:   quiz.Module.Title,
		HostUsername: quiz.Host.Name,
		Status:       quiz.Status,
		Mode:         quiz.Mode,
		Message:      "Successfully joined the quiz",
	}

	return response, nil
}

func (s *quizSessionService) AnswerQuestion(userID uint, req request.AnswerQuestionRequest) (*models.QuizSession, error) {
	// Get question
	question, err := s.questionRepo.GetQuestionByID(int32(req.QuestionID))
	if err != nil {
		return nil, fmt.Errorf("question not found: %v", err)
	}

	// Find participant
	participants, err := s.participantRepo.GetParticipantsByUserID(userID)
	if err != nil || len(participants) == 0 {
		return nil, fmt.Errorf("participant not found")
	}

	// Assume the latest active quiz
	var participant *models.Participant
	for _, p := range participants {
		if p.Quiz.Status == "in_progress" {
			participant = &p
			break
		}
	}

	if participant == nil {
		return nil, fmt.Errorf("no active quiz found for user")
	}

	// Check if user already answered this question
	_, err = s.quizSessionRepo.GetQuizSession(participant.QuizID, participant.ID, req.QuestionID)
	if err == nil {
		return nil, fmt.Errorf("question already answered")
	}

	// Calculate points based on correctness and response time
	isCorrect := question.CorrectAnswer == req.UserAnswer
	points := s.calculatePoints(isCorrect, req.ResponseTime, question.AnswerTime)

	// Create quiz session
	session := &models.QuizSession{
		QuizID:        participant.QuizID,
		ParticipantID: participant.ID,
		QuestionID:    req.QuestionID,
		UserAnswer:    req.UserAnswer,
		IsCorrect:     isCorrect,
		ResponseTime:  req.ResponseTime,
		PointsEarned:  points,
	}

	createdSession, err := s.quizSessionRepo.CreateQuizSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to save answer: %v", err)
	}

	// Update participant stats
	s.updateParticipantStats(participant, isCorrect, points)

	return createdSession, nil
}

func (s *quizSessionService) calculatePoints(isCorrect bool, responseTime int32, maxTime int32) int {
	if !isCorrect {
		return 0
	}

	// Base points for correct answer
	basePoints := 100

	// Time bonus: faster response = more points
	timeBonus := int(float64(maxTime-responseTime) / float64(maxTime) * 50)
	if timeBonus < 0 {
		timeBonus = 0
	}

	return basePoints + timeBonus
}

func (s *quizSessionService) updateParticipantStats(participant *models.Participant, isCorrect bool, points int) error {
	participant.TotalScore += points

	if isCorrect {
		participant.CorrectAnswers++
		participant.CurrentStreak++
		if participant.CurrentStreak > participant.ConsecutiveCorrect {
			participant.ConsecutiveCorrect = participant.CurrentStreak
		}
	} else {
		participant.WrongAnswers++
		participant.CurrentStreak = 0
	}

	_, err := s.participantRepo.UpdateParticipant(int32(participant.ID), participant)
	return err
}

func (s *quizSessionService) GetQuizSession(quizID uint) (*response.QuizSessionResponse, error) {
	quiz, err := s.quizRepo.GetQuizByID(quizID)
	if err != nil {
		return nil, fmt.Errorf("quiz not found: %v", err)
	}

	// Build participants status
	var participants []response.ParticipantStatus
	for _, p := range quiz.Participants {
		participants = append(participants, response.ParticipantStatus{
			UserID:        p.UserID,
			Username:      p.User.Name,
			CurrentScore:  p.TotalScore,
			CurrentStreak: p.CurrentStreak,
			IsFinished:    p.IsFinished,
		})
	}

	response := &response.QuizSessionResponse{
		QuizID:         quiz.ID,
		TotalQuestions: len(quiz.Questions),
		Participants:   participants,
		Phase:          quiz.Status,
	}

	return response, nil
}

func (s *quizSessionService) GetQuizResult(quizID uint) (*response.QuizResultResponse, error) {
	quiz, err := s.quizRepo.GetQuizByID(quizID)
	if err != nil {
		return nil, fmt.Errorf("quiz not found: %v", err)
	}

	var participantResults []response.ParticipantResult
	var winner *response.ParticipantResult

	for _, p := range quiz.Participants {
		result := response.ParticipantResult{
			UserID:             p.UserID,
			Username:           p.User.Name,
			TotalScore:         p.TotalScore,
			CorrectAnswers:     p.CorrectAnswers,
			WrongAnswers:       p.WrongAnswers,
			ConsecutiveCorrect: p.ConsecutiveCorrect,
			IsFinished:         p.IsFinished,
		}

		if p.FinishedAt != nil && p.FinishedAt.Valid {
			result.FinishedAt = &p.FinishedAt.Time
		}

		participantResults = append(participantResults, result)

		// Determine winner (highest score)
		if winner == nil || result.TotalScore > winner.TotalScore {
			winner = &result
		}
	}

	var moduleName string
	if quiz.Module != nil {
		moduleName = quiz.Module.Title
	}

	response := &response.QuizResultResponse{
		QuizID:             quiz.ID,
		Mode:               quiz.Mode,
		Status:             quiz.Status,
		ModuleName:         moduleName,
		TotalQuestions:     len(quiz.Questions),
		ParticipantResults: participantResults,
		Winner:             winner,
		CreatedAt:          *quiz.CreatedAt,
	}

	if quiz.UpdatedAt != nil {
		completedAt := *quiz.UpdatedAt
		response.CompletedAt = &completedAt
	}

	return response, nil
}

func (s *quizSessionService) FinishQuiz(quizID, userID uint) error {
	// Find participant
	participants, err := s.participantRepo.GetParticipantsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to find participant: %v", err)
	}

	var participant *models.Participant
	for _, p := range participants {
		if p.QuizID == quizID {
			participant = &p
			break
		}
	}

	if participant == nil {
		return fmt.Errorf("participant not found in this quiz")
	}

	// Mark participant as finished
	participant.IsFinished = true
	now := gorm.DeletedAt{Time: time.Now(), Valid: true}
	participant.FinishedAt = &now

	_, err = s.participantRepo.UpdateParticipant(int32(participant.ID), participant)
	if err != nil {
		return fmt.Errorf("failed to update participant: %v", err)
	}

	// Check if all participants are finished to update quiz status
	quiz, err := s.quizRepo.GetQuizByID(quizID)
	if err != nil {
		return fmt.Errorf("failed to get quiz: %v", err)
	}

	allFinished := true
	for _, p := range quiz.Participants {
		if !p.IsFinished {
			allFinished = false
			break
		}
	}

	if allFinished {
		quiz.Status = "completed"
		// Determine winner
		var winnerID *uint
		maxScore := -1
		for _, p := range quiz.Participants {
			if p.TotalScore > maxScore {
				maxScore = p.TotalScore
				winnerID = &p.UserID
			}
		}
		quiz.WinnerID = winnerID

		_, err = s.quizRepo.UpdateQuiz(quiz.ID, quiz)
		if err != nil {
			return fmt.Errorf("failed to update quiz status: %v", err)
		}
	}

	return nil
}
