package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type VideoQuizRepositoryInterface interface {
	CreateVideoQuiz(videoQuiz *models.VideoQuiz) (*models.VideoQuiz, error)
	GetVideoQuizByID(id uint) (*models.VideoQuiz, error)
	GetVideoQuizzesByVideoMaterialID(videoMaterialID uint) ([]models.VideoQuiz, error)
	UpdateVideoQuiz(id uint, videoQuiz *models.VideoQuiz) (*models.VideoQuiz, error)
	DeleteVideoQuiz(id uint) error

	// User Answer methods
	CreateVideoQuizUserAnswer(userAnswer *models.VideoQuizUserAnswer) (*models.VideoQuizUserAnswer, error)
	GetUserVideoQuizAnswers(userID uint, videoMaterialID uint) ([]models.VideoQuizUserAnswer, error)
	GetAllUserVideoQuizAnswers(userID uint) ([]models.VideoQuizUserAnswer, error)
	HasUserAnsweredVideoQuiz(userID uint, videoQuizID uint) (bool, error)
}

type videoQuizRepository struct {
	db *gorm.DB
}

func NewVideoQuizRepository(db *gorm.DB) VideoQuizRepositoryInterface {
	return &videoQuizRepository{db}
}

func (r *videoQuizRepository) CreateVideoQuiz(videoQuiz *models.VideoQuiz) (*models.VideoQuiz, error) {
	if err := r.db.Create(videoQuiz).Error; err != nil {
		return nil, err
	}
	return videoQuiz, nil
}

func (r *videoQuizRepository) GetVideoQuizByID(id uint) (*models.VideoQuiz, error) {
	var videoQuiz models.VideoQuiz
	if err := r.db.Preload("VideoMaterial").First(&videoQuiz, id).Error; err != nil {
		return nil, err
	}
	return &videoQuiz, nil
}

func (r *videoQuizRepository) GetVideoQuizzesByVideoMaterialID(videoMaterialID uint) ([]models.VideoQuiz, error) {
	var videoQuizzes []models.VideoQuiz
	if err := r.db.Where("video_material_id = ?", videoMaterialID).Order("timestamp_start ASC").Find(&videoQuizzes).Error; err != nil {
		return nil, err
	}
	return videoQuizzes, nil
}

func (r *videoQuizRepository) UpdateVideoQuiz(id uint, videoQuiz *models.VideoQuiz) (*models.VideoQuiz, error) {
	existingVideoQuiz, err := r.GetVideoQuizByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingVideoQuiz.Question = videoQuiz.Question
	existingVideoQuiz.TimestampStart = videoQuiz.TimestampStart
	existingVideoQuiz.TimestampEnd = videoQuiz.TimestampEnd
	existingVideoQuiz.Options = videoQuiz.Options
	existingVideoQuiz.CorrectAnswer = videoQuiz.CorrectAnswer
	existingVideoQuiz.Explanation = videoQuiz.Explanation
	existingVideoQuiz.Order = videoQuiz.Order

	if err := r.db.Save(existingVideoQuiz).Error; err != nil {
		return nil, err
	}

	return existingVideoQuiz, nil
}

func (r *videoQuizRepository) DeleteVideoQuiz(id uint) error {
	return r.db.Delete(&models.VideoQuiz{}, id).Error
}

func (r *videoQuizRepository) CreateVideoQuizUserAnswer(userAnswer *models.VideoQuizUserAnswer) (*models.VideoQuizUserAnswer, error) {
	if err := r.db.Create(userAnswer).Error; err != nil {
		return nil, err
	}
	return userAnswer, nil
}

func (r *videoQuizRepository) GetUserVideoQuizAnswers(userID uint, videoMaterialID uint) ([]models.VideoQuizUserAnswer, error) {
	var userAnswers []models.VideoQuizUserAnswer

	err := r.db.
		Joins("JOIN video_quizzes ON video_quiz_user_answers.video_quiz_id = video_quizzes.id").
		Where("video_quiz_user_answers.user_id = ? AND video_quizzes.video_material_id = ?", userID, videoMaterialID).
		Preload("VideoQuiz").
		Find(&userAnswers).Error

	if err != nil {
		return nil, err
	}

	return userAnswers, nil
}

func (r *videoQuizRepository) GetAllUserVideoQuizAnswers(userID uint) ([]models.VideoQuizUserAnswer, error) {
	var userAnswers []models.VideoQuizUserAnswer

	err := r.db.Preload("User").Preload("VideoQuiz").Where("user_id = ?", userID).Find(&userAnswers).Error
	if err != nil {
		return nil, err
	}

	return userAnswers, nil
}

func (r *videoQuizRepository) HasUserAnsweredVideoQuiz(userID uint, videoQuizID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.VideoQuizUserAnswer{}).
		Where("user_id = ? AND video_quiz_id = ?", userID, videoQuizID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
