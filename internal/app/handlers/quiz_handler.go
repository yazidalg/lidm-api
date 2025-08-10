package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/request"
	"github.com/yazidalg/lidm_backend/internal/app/services"
	"github.com/yazidalg/lidm_backend/internal/utils"
)

type QuizHandler struct {
	quizService        services.QuizServiceInterface
	participantService services.ParticipantServiceInterface
}

func NewQuizHandler(
	participantService services.ParticipantServiceInterface,
	quizService services.QuizServiceInterface,
) *QuizHandler {
	return &QuizHandler{
		quizService:        quizService,
		participantService: participantService,
	}
}

// CreateQuizLobby membuat lobi kuis baru dan mengembalikan invite code.
func (h *QuizHandler) CreateQuizLobby(c *gin.Context) {
	var req request.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	if req.ModuleID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ModuleID is required"})
		return
	}

	req.HostUserID = user.ID
	req.Mode = "multiplayer"
	req.Status = "pending"
	req.InviteCode = utils.GenerateInviteCode(6)

	newQuiz, err := h.quizService.CreateQuiz(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz lobby: " + err.Error()})
		return
	}

	// Otomatis buat partisipan untuk host
	_, err = h.participantService.CreateParticipant(request.CreateParticipantRequest{
		UserID: user.ID,
		QuizID: newQuiz.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create participant for host: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Quiz lobby created successfully",
		"quiz_id":     newQuiz.ID,
		"invite_code": newQuiz.InviteCode,
		"module_id":   newQuiz.ModuleID,
	})
}

// JoinQuizLobby memungkinkan user lain bergabung menggunakan invite code.
func (h *QuizHandler) JoinQuizLobby(c *gin.Context) {
	var req request.JoinQuizWithCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	quiz, err := h.quizService.GetQuizByInviteCode(req.InviteCode)
	if err != nil || quiz.Status != "pending" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lobby not found or already in progress"})
		return
	}

	for _, p := range quiz.Participants {
		if p.UserID == user.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You have already joined this lobby"})
			return
		}
	}

	if len(quiz.Participants) >= 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Lobby is full"})
		return
	}

	_, err = h.participantService.CreateParticipant(request.CreateParticipantRequest{
		UserID: user.ID,
		QuizID: quiz.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join lobby: " + err.Error()})
		return
	}

	// Jika setelah join jumlah pemain menjadi 2, ubah status kuis
	if len(quiz.Participants)+1 == 2 {
		_, err := h.quizService.UpdateQuiz(quiz.ID, request.UpdateQuizRequest{Status: "in_progress"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start quiz: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined the quiz lobby",
		"quiz_id": quiz.ID,
	})
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var request request.CreateQuizRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.quizService.CreateQuiz(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to create quiz",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz created successfully",
		"data":    result,
	})
}

func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to get quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to get quiz",
		})
		return
	}

	quiz, err := h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz retrieved successfully",
		"data":    quiz,
	})
}

// GetQuizzesByModule returns all quizzes for a given module id
func (h *QuizHandler) GetQuizzesByModule(c *gin.Context) {
	moduleIDParam := c.Param("module_id")
	if moduleIDParam == "" { c.JSON(http.StatusBadRequest, gin.H{"error":"module_id is required"}); return }
	idUint64, err := strconv.ParseUint(moduleIDParam, 10, 32)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid module_id"}); return }
	quizzes, err := h.quizService.GetQuizzesByModule(uint(idUint64))
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"data": quizzes})
}

func (h *QuizHandler) GetAllQuizzes(c *gin.Context) {
	quizzes, err := h.quizService.GetAllQuizzes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quizzes retrieved successfully",
		"data":    quizzes,
	})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	var request request.UpdateQuizRequest
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to update quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to update quiz",
		})
		return
	}

	// Check if quiz exists
	_, err = h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update quiz",
		})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.quizService.UpdateQuiz(uint(id), request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz updated successfully",
		"data":    result,
	})
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID is required",
			"message": "Failed to delete quiz",
		})
		return
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"message": "Failed to delete quiz",
		})
		return
	}

	// Check if quiz exists
	_, err = h.quizService.GetQuizByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete quiz",
		})
		return
	}

	err = h.quizService.DeleteQuiz(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to delete quiz",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quiz deleted successfully",
	})
}
