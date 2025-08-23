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
	UpdateModule(id uint32, request request.UpdateModuleRequest) (*models.Module, error)
	DeleteModule(id uint32) error
}

type moduleService struct {
	moduleRepository repositories.ModuleRepositoryInterface
	progressRepo     repositories.ProgressRepositoryInterface
	videoQuizRepo    repositories.VideoQuizRepositoryInterface
	prequizRepo      repositories.PrequizRepositoryInterface
}

func NewModuleService(
	moduleRepository repositories.ModuleRepositoryInterface,
	progressRepo repositories.ProgressRepositoryInterface,
	videoQuizRepo repositories.VideoQuizRepositoryInterface,
	prequizRepo repositories.PrequizRepositoryInterface,
) *moduleService {
	return &moduleService{
		moduleRepository: moduleRepository,
		progressRepo:     progressRepo,
		videoQuizRepo:    videoQuizRepo,
		prequizRepo:      prequizRepo,
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
	if len(request.VideoMaterial) > 0 {
		module.VideoMaterial = request.VideoMaterial
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
		module.VideoMaterial = []models.VideoMaterial{*request.VideoMaterial}
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
