package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type ModuleProgressServiceInterface interface {
	GetUserModuleProgress(userID, moduleID uint) (*models.ModuleProgress, error)
	UpdateUserProgress(userID, moduleID uint) (*models.ModuleProgress, error)
	GetAllUserProgress(userID uint) ([]models.ModuleProgress, error)
	InitializeUserProgress(userID uint) error
	GetUnlockedModulesForUser(userID uint) ([]models.ModuleProgress, error)
	CheckModuleAccess(userID, moduleID uint) (bool, error)
}

type moduleProgressService struct {
	moduleProgressRepo repositories.ModuleProgressRepositoryInterface
	userRepo           repositories.UserRepositoryInterface
	moduleRepo         repositories.ModuleRepositoryInterface
}

func NewModuleProgressService(
	moduleProgressRepo repositories.ModuleProgressRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	moduleRepo repositories.ModuleRepositoryInterface,
) *moduleProgressService {
	return &moduleProgressService{
		moduleProgressRepo: moduleProgressRepo,
		userRepo:           userRepo,
		moduleRepo:         moduleRepo,
	}
}

func (s *moduleProgressService) GetUserModuleProgress(userID, moduleID uint) (*models.ModuleProgress, error) {
	return s.moduleProgressRepo.GetByUserAndModule(userID, moduleID)
}

func (s *moduleProgressService) UpdateUserProgress(userID, moduleID uint) (*models.ModuleProgress, error) {
	// Get or create progress entry
	progress, err := s.moduleProgressRepo.GetByUserAndModule(userID, moduleID)
	if err != nil {
		// If not found, create new progress (but don't unlock unless it's module 1)
		modules, err := s.moduleRepo.GetAllModules()
		if err != nil || len(modules) == 0 {
			return nil, err
		}
		
		firstModuleID := modules[0].ID // Assume modules are ordered by ID
		isUnlocked := moduleID == firstModuleID // Only unlock if it's the first module
		
		progress = &models.ModuleProgress{
			UserID:     userID,
			ModuleID:   moduleID,
			IsUnlocked: isUnlocked,
			IsComplete: false,
			Progress:   0,
		}
		
		if isUnlocked {
			progress.MarkAsStarted()
		}
		
		return s.moduleProgressRepo.CreateModuleProgress(progress)
	}

	// Update existing progress
	return s.moduleProgressRepo.UpdateProgress(progress)
}

func (s *moduleProgressService) GetAllUserProgress(userID uint) ([]models.ModuleProgress, error) {
	return s.moduleProgressRepo.GetAllByUser(userID)
}

func (s *moduleProgressService) InitializeUserProgress(userID uint) error {
	return s.moduleProgressRepo.InitializeFirstModuleForUser(userID)
}

func (s *moduleProgressService) GetUnlockedModulesForUser(userID uint) ([]models.ModuleProgress, error) {
	return s.moduleProgressRepo.GetUnlockedModulesForUser(userID)
}

func (s *moduleProgressService) CheckModuleAccess(userID, moduleID uint) (bool, error) {
	progress, err := s.moduleProgressRepo.GetByUserAndModule(userID, moduleID)
	if err != nil {
		// If no progress record, check if it's the first module
		modules, err := s.moduleRepo.GetAllModules()
		if err != nil || len(modules) == 0 {
			return false, err
		}
		
		firstModuleID := modules[0].ID // Assume modules are ordered by ID
		return moduleID == firstModuleID, nil
	}

	return progress.IsUnlocked, nil
}
