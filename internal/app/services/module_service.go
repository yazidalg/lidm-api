package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type ModuleServiceInterface interface {
	CreateModule(request request.ModuleRequest) (*models.Module, error)
	CreateModuleWithVideo(request request.CreateModuleWithVideoRequest) (*models.Module, error)
	GetModuleByID(id uint32) (*models.Module, error)
	GetAllModules() ([]models.Module, error)
	GetAllModulesWithProgress(userID uint) ([]map[string]interface{}, error)
	GetAllModulesWithUnlockStatus(userID uint) ([]map[string]interface{}, error)
	UpdateModule(id uint32, request request.UpdateModuleRequest) (*models.Module, error)
	DeleteModule(id uint32) error
}

type moduleService struct {
	moduleRepository        repositories.ModuleRepositoryInterface
	progressRepo            repositories.ProgressRepositoryInterface
	videoQuizRepo           repositories.VideoQuizRepositoryInterface
	prequizRepo             repositories.PrequizRepositoryInterface
	moduleProgressRepo      repositories.ModuleProgressRepositoryInterface
	moduleProgressService   ModuleProgressServiceInterface
}

func NewModuleService(
	moduleRepository repositories.ModuleRepositoryInterface,
	progressRepo repositories.ProgressRepositoryInterface,
	videoQuizRepo repositories.VideoQuizRepositoryInterface,
	prequizRepo repositories.PrequizRepositoryInterface,
	moduleProgressRepo repositories.ModuleProgressRepositoryInterface,
	moduleProgressService ModuleProgressServiceInterface,
) *moduleService {
	return &moduleService{
		moduleRepository:      moduleRepository,
		progressRepo:          progressRepo,
		videoQuizRepo:         videoQuizRepo,
		prequizRepo:           prequizRepo,
		moduleProgressRepo:    moduleProgressRepo,
		moduleProgressService: moduleProgressService,
	}
}

func (s *moduleService) CreateModule(request request.ModuleRequest) (*models.Module, error) {
	module := &models.Module{
		Title:       request.Title,
		Description: request.Description,
	}

	// Handle optional fields
	if request.OffsetX != nil {
		module.OffsetX = uint16(*request.OffsetX)
	}
	if request.OffsetY != nil {
		module.OffsetY = uint16(*request.OffsetY)
	}
	if request.Icon != nil {
		module.Icon = *request.Icon
	}
	if request.Thumbnail != nil {
		module.Thumbnail = *request.Thumbnail
	}

	result, err := s.moduleRepository.CreateModule(module)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *moduleService) CreateModuleWithVideo(request request.CreateModuleWithVideoRequest) (*models.Module, error) {
	module := &models.Module{
		Title:       request.Title,
		Description: request.Description,
	}

	// Handle optional fields
	if request.OffsetX != nil {
		module.OffsetX = uint16(*request.OffsetX)
	}
	if request.OffsetY != nil {
		module.OffsetY = uint16(*request.OffsetY)
	}
	if request.Icon != nil {
		module.Icon = *request.Icon
	}
	if request.Thumbnail != nil {
		module.Thumbnail = *request.Thumbnail
	}

	// Handle VideoMaterial
	if request.VideoMaterial != nil && len(request.VideoMaterial) > 0 {
		// Take only the first VideoMaterial if multiple provided
		module.VideoMaterial = &request.VideoMaterial[0]
	}

	result, err := s.moduleRepository.CreateModule(module)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *moduleService) GetModuleByID(id uint32) (*models.Module, error) {
	result, err := s.moduleRepository.GetModuleByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *moduleService) GetAllModules() ([]models.Module, error) {
	result, err := s.moduleRepository.GetAllModules()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *moduleService) GetAllModulesWithProgress(userID uint) ([]map[string]interface{}, error) {
	modules, err := s.moduleRepository.GetAllModules()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	// Fetch all user prequiz answers once and build a lookup map for fast checks
	userAnswers, _ := s.prequizRepo.GetUserPrequizAnswers(userID)
	answeredPrequizMap := make(map[uint]bool, len(userAnswers))
	for _, answer := range userAnswers {
		answeredPrequizMap[answer.PrequizID] = true
	}

	for _, module := range modules {
		moduleData := map[string]interface{}{
			"ID":            module.ID,
			"Title":         module.Title,
			"Description":   module.Description,
			"Thumbnail":     module.Thumbnail,
			"Icon":          module.Icon,
			"OffsetX":       module.OffsetX,
			"OffsetY":       module.OffsetY,
			"CreatedAt":     module.CreatedAt,
			"UpdatedAt":     module.UpdatedAt,
			"VideoMaterial": module.VideoMaterial,
		}

		result = append(result, moduleData)
	}

	return result, nil
}

func (s *moduleService) UpdateModule(id uint32, request request.UpdateModuleRequest) (*models.Module, error) {
	module := &models.Module{
		Title:       request.Title,
		Description: request.Description,
	}

	// Handle optional fields
	if request.OffsetX != nil {
		module.OffsetX = uint16(*request.OffsetX)
	}
	if request.OffsetY != nil {
		module.OffsetY = uint16(*request.OffsetY)
	}
	if request.Icon != nil {
		module.Icon = *request.Icon
	}
	if request.Thumbnail != nil {
		module.Thumbnail = *request.Thumbnail
	}

	// Handle VideoMaterial if provided
	if request.VideoMaterial != nil {
		// Set the ModuleID for the VideoMaterial
		request.VideoMaterial.ModuleID = uint(id)
		// Note: You'll need to implement VideoMaterial creation in repository
		// For now, we'll just store the reference
		module.VideoMaterial = request.VideoMaterial
	}

	result, err := s.moduleRepository.UpdateModule(id, module)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *moduleService) DeleteModule(id uint32) error {
	err := s.moduleRepository.DeleteModule(id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllModulesWithUnlockStatus returns all modules with unlock status and progress for a user
func (s *moduleService) GetAllModulesWithUnlockStatus(userID uint) ([]map[string]interface{}, error) {
	modules, err := s.moduleRepository.GetAllModules()
	if err != nil {
		return nil, err
	}

	// Initialize first module for user if no progress exists
	s.moduleProgressService.InitializeUserProgress(userID)

	result := make([]map[string]interface{}, 0)

	for _, module := range modules {
		// Get module progress for this user
		moduleProgress, err := s.moduleProgressService.GetUserModuleProgress(userID, module.ID)
		
		isUnlocked := false
		isComplete := false
		progress := float32(0)
		var startedAt, completedAt interface{}

		if err == nil && moduleProgress != nil {
			isUnlocked = moduleProgress.IsUnlocked
			isComplete = moduleProgress.IsComplete
			progress = moduleProgress.Progress
			startedAt = moduleProgress.StartedAt
			completedAt = moduleProgress.CompletedAt
		} else {
			// Check if this is the first module (should be unlocked by default)
			hasAccess, _ := s.moduleProgressService.CheckModuleAccess(userID, module.ID)
			isUnlocked = hasAccess
		}

		// Count total content items
		totalPrequizzes := len(module.Prequizzes)
		totalVideoQuizzes := 0
		if module.VideoMaterial != nil {
			totalVideoQuizzes = len(module.VideoMaterial.VideoQuizzes)
		}
		totalFlashcards := len(module.Flashcards)

		// Limit prequizzes to 3 for display
		limitedPrequizzes := module.Prequizzes
		if len(module.Prequizzes) > 3 {
			limitedPrequizzes = module.Prequizzes[:3]
		}

		// Get completion counts if module is unlocked
		answeredPrequizzes := 0
		answeredVideoQuizzes := 0
		if isUnlocked {
			userAnswers, _ := s.prequizRepo.GetUserPrequizAnswers(userID)
			answeredPrequizMap := make(map[uint]bool)
			for _, answer := range userAnswers {
				answeredPrequizMap[answer.PrequizID] = true
			}
			for _, prequiz := range limitedPrequizzes {
				if answeredPrequizMap[prequiz.ID] {
					answeredPrequizzes++
				}
			}

			// Count video quiz answers
			videoAnswers, _ := s.videoQuizRepo.GetAllUserVideoQuizAnswers(userID)
			answeredVideoMap := make(map[uint]bool)
			for _, answer := range videoAnswers {
				answeredVideoMap[answer.VideoQuizID] = true
			}
			if module.VideoMaterial != nil {
				for _, quiz := range module.VideoMaterial.VideoQuizzes {
					if answeredVideoMap[quiz.ID] {
						answeredVideoQuizzes++
					}
				}
			}
		}

		moduleData := map[string]interface{}{
			"id":          module.ID,
			"title":       module.Title,
			"description": module.Description,
			"thumbnail":   module.Thumbnail,
			"icon":        module.Icon,
			"offset_x":    module.OffsetX,
			"offset_y":    module.OffsetY,
			"created_at":  module.CreatedAt,
			"updated_at":  module.UpdatedAt,
			
			// Progress and unlock status
			"is_unlocked":            isUnlocked,
			"is_complete":            isComplete,
			"progress":               progress,
			"started_at":             startedAt,
			"completed_at":           completedAt,
			
			// Content counts
			"total_prequizzes":       totalPrequizzes,
			"answered_prequizzes":    answeredPrequizzes,
			"total_video_quizzes":    totalVideoQuizzes,
			"answered_video_quizzes": answeredVideoQuizzes,
			"total_flashcards":       totalFlashcards,
			
			// Include content based on unlock status
			"video_material": func() interface{} {
				if isUnlocked {
					return module.VideoMaterial
				}
				return nil
			}(),
			"ar_experiment": func() interface{} {
				if isUnlocked {
					return module.ARExperiment
				}
				return nil
			}(),
			"prequizzes": limitedPrequizzes, // Show only 3 prequizzes regardless of unlock status
			"flashcards": func() interface{} {
				if isUnlocked {
					return module.Flashcards
				}
				return []interface{}{}
			}(),
		}

		result = append(result, moduleData)
	}

	return result, nil
}
