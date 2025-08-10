package services

import (
	"errors"
	"time"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type VideoQuizServiceInterface interface {
	CreateVideoQuiz(videoQuiz request.VideoQuizRequest) (*models.VideoQuiz, error)
	GetVideoQuizByID(id uint) (*models.VideoQuiz, error)
	GetVideoQuizzesByVideoMaterialID(videoMaterialID uint) ([]models.VideoQuiz, error)
	UpdateVideoQuiz(id uint, videoQuiz request.VideoQuizRequest) (*models.VideoQuiz, error)
	DeleteVideoQuiz(id uint) error
	
	// User Answer methods
	SubmitVideoQuizAnswer(userID uint, answer request.VideoQuizAnswerRequest) (*models.VideoQuizUserAnswer, error)
	GetUserVideoQuizAnswers(userID uint, videoMaterialID uint) ([]models.VideoQuizUserAnswer, error)
	GetAllUserVideoQuizAnswers(userID uint) ([]models.VideoQuizUserAnswer, error)
	HasUserAnsweredVideoQuiz(userID uint, videoQuizID uint) (bool, error)
}

type videoQuizService struct {
	videoQuizRepo repositories.VideoQuizRepositoryInterface
	userRepo      repositories.UserRepositoryInterface
}

func NewVideoQuizService(videoQuizRepo repositories.VideoQuizRepositoryInterface, userRepo repositories.UserRepositoryInterface) VideoQuizServiceInterface {
	return &videoQuizService{
		videoQuizRepo: videoQuizRepo,
		userRepo:      userRepo,
	}
}

func (s *videoQuizService) CreateVideoQuiz(videoQuiz request.VideoQuizRequest) (*models.VideoQuiz, error) {
	videoQuizModel := models.VideoQuiz{
		VideoMaterialID: videoQuiz.VideoMaterialID,
		Question:        videoQuiz.Question,
		TimestampStart:  videoQuiz.TimestampStart,
		TimestampEnd:    videoQuiz.TimestampEnd,
		Options:         models.Options(videoQuiz.Options),
		CorrectAnswer:   videoQuiz.CorrectAnswer,
		Explanation:     videoQuiz.Explanation,
		Order:           videoQuiz.Order,
	}

	return s.videoQuizRepo.CreateVideoQuiz(&videoQuizModel)
}

func (s *videoQuizService) GetVideoQuizByID(id uint) (*models.VideoQuiz, error) {
	return s.videoQuizRepo.GetVideoQuizByID(id)
}

func (s *videoQuizService) GetVideoQuizzesByVideoMaterialID(videoMaterialID uint) ([]models.VideoQuiz, error) {
	return s.videoQuizRepo.GetVideoQuizzesByVideoMaterialID(videoMaterialID)
}

func (s *videoQuizService) UpdateVideoQuiz(id uint, videoQuiz request.VideoQuizRequest) (*models.VideoQuiz, error) {
	videoQuizModel := models.VideoQuiz{
		VideoMaterialID: videoQuiz.VideoMaterialID,
		Question:        videoQuiz.Question,
		TimestampStart:  videoQuiz.TimestampStart,
		TimestampEnd:    videoQuiz.TimestampEnd,
		Options:         models.Options(videoQuiz.Options),
		CorrectAnswer:   videoQuiz.CorrectAnswer,
		Explanation:     videoQuiz.Explanation,
		Order:           videoQuiz.Order,
	}

	return s.videoQuizRepo.UpdateVideoQuiz(id, &videoQuizModel)
}

func (s *videoQuizService) DeleteVideoQuiz(id uint) error {
	return s.videoQuizRepo.DeleteVideoQuiz(id)
}

func (s *videoQuizService) SubmitVideoQuizAnswer(userID uint, answer request.VideoQuizAnswerRequest) (*models.VideoQuizUserAnswer, error) {
	// Check if user already answered this quiz
	hasAnswered, err := s.videoQuizRepo.HasUserAnsweredVideoQuiz(userID, answer.VideoQuizID)
	if err != nil {
		return nil, err
	}
	
	if hasAnswered {
		return nil, errors.New("user has already answered this video quiz")
	}

	// Get the video quiz to check correct answer
	videoQuiz, err := s.videoQuizRepo.GetVideoQuizByID(answer.VideoQuizID)
	if err != nil {
		return nil, err
	}

	// Check if answer is correct
	isCorrect := videoQuiz.CorrectAnswer == answer.SelectedAnswer

	userAnswer := models.VideoQuizUserAnswer{
		VideoQuizID:    answer.VideoQuizID,
		UserID:         userID,
		SelectedAnswer: answer.SelectedAnswer,
		IsCorrect:      isCorrect,
		AnsweredAt:     time.Now().Unix(),
		ResponseTime:   answer.ResponseTime,
	}

	return s.videoQuizRepo.CreateVideoQuizUserAnswer(&userAnswer)
}

func (s *videoQuizService) GetUserVideoQuizAnswers(userID uint, videoMaterialID uint) ([]models.VideoQuizUserAnswer, error) {
	return s.videoQuizRepo.GetUserVideoQuizAnswers(userID, videoMaterialID)
}

func (s *videoQuizService) GetAllUserVideoQuizAnswers(userID uint) ([]models.VideoQuizUserAnswer, error) {
	return s.videoQuizRepo.GetAllUserVideoQuizAnswers(userID)
}

func (s *videoQuizService) HasUserAnsweredVideoQuiz(userID uint, videoQuizID uint) (bool, error) {
	return s.videoQuizRepo.HasUserAnsweredVideoQuiz(userID, videoQuizID)
}
