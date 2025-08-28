package repositories

import (
	"math/rand"

	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type QuestionRepositoryInterface interface {
	CreateQuestion(question *models.Question) (*models.Question, error)
	GetQuestionByID(id int32) (*models.Question, error)
	GetAllQuestions() ([]models.Question, error)
	UpdateQuestion(id int32, question *models.Question) (*models.Question, error)
	DeleteQuestion(id int32) error
	GetRandomQuestion(count int) (*[]models.Question, error)
	GetRandomQuestionsByModule(moduleID uint, count int) (*[]models.Question, error)
	GetQuestionsByModuleAndMode(moduleID uint, mode string) (*[]models.Question, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *questionRepository {
	return &questionRepository{db}
}

func (r *questionRepository) CreateQuestion(question *models.Question) (*models.Question, error) {
	err := r.db.Create(question).Error
	return question, err
}

func (r *questionRepository) GetQuestionByID(id int32) (*models.Question, error) {
	var question models.Question
	err := r.db.First(&question, id).Error
	return &question, err
}

func (r *questionRepository) GetAllQuestions() ([]models.Question, error) {
	var questions []models.Question
	err := r.db.Find(&questions).Error
	return questions, err
}

func (r *questionRepository) UpdateQuestion(id int32, question *models.Question) (*models.Question, error) {
	err := r.db.Model(&models.Question{}).Where("id = ?", id).Updates(question).Error
	return question, err
}

func (r *questionRepository) DeleteQuestion(id int32) error {
	return r.db.Delete(&models.Question{}, id).Error
}

func (r *questionRepository) GetRandomQuestion(count int) (*[]models.Question, error) {
	var question []models.Question
	err := r.db.Order("RAND()").Limit(count).Find(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetRandomQuestionsByModule(moduleID uint, count int) (*[]models.Question, error) {
	var questions []models.Question
	err := r.db.Where("module_id = ?", moduleID).
		Order("RAND()").
		Limit(count).
		Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return &questions, nil
}

func (r *questionRepository) GetQuestionsByModuleAndMode(moduleID uint, mode string) (*[]models.Question, error) {
	var questions []models.Question
	var hotsCount, regularCount int

	// Determine question mix based on mode
	if mode == "multiplayer" {
		hotsCount = 3
		regularCount = 7
	} else { // solo or single_player
		hotsCount = 4
		regularCount = 6
	}

	// Get HOTS questions
	var hotsQuestions []models.Question
	err := r.db.Where("module_id = ? AND question_type = ?", moduleID, "hots").
		Order("RAND()").
		Limit(hotsCount).
		Find(&hotsQuestions).Error
	if err != nil {
		return nil, err
	}

	// Get regular questions
	var regularQuestions []models.Question
	err = r.db.Where("module_id = ? AND (question_type = ? OR question_type IS NULL OR question_type = '')", moduleID, "regular").
		Order("RAND()").
		Limit(regularCount).
		Find(&regularQuestions).Error
	if err != nil {
		return nil, err
	}

	// Combine and shuffle
	questions = append(questions, hotsQuestions...)
	questions = append(questions, regularQuestions...)

	// Shuffle the combined questions
	for i := len(questions) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		questions[i], questions[j] = questions[j], questions[i]
	}

	return &questions, nil
}
