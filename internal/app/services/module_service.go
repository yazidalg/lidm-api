package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type ModuleServiceInterface interface {
	GetModuleByID(id uint32) (*models.Module, error)
	GetAllModules() ([]models.Module, error)
	UpdateModule(id uint32, request request.UpdateModuleRequest) (*models.Module, error)
	DeleteModule(id uint32) error
}

type moduleService struct {
	moduleRepository repositories.ModuleRepositoryInterface
}

func NewModuleService(moduleRepository repositories.ModuleRepositoryInterface) *moduleService {
	return &moduleService{moduleRepository}
}

func (s *moduleService) GetModuleByID(id uint32) (*models.Module, error) {
	return s.moduleRepository.GetModuleByID(id)
}

func (s *moduleService) GetAllModules() ([]models.Module, error) {
	return s.moduleRepository.GetAllModules()
}

func (s *moduleService) UpdateModule(id uint32, request request.UpdateModuleRequest) (*models.Module, error) {
	// Get the existing module
	module := &models.Module{
		Title:       request.Title,
		Description: request.Description,
		SortOrder:   uint16(request.SortOrder),
	}

	// Delegate to repository for update
	return s.moduleRepository.UpdateModule(id, module)
}

func (s *moduleService) DeleteModule(id uint32) error {
	// Delegate to repository for deletion
	return s.moduleRepository.DeleteModule(id)
}
