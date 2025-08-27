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
	CheckAndUnlockNextModule(userID, currentModuleID uint) error
	// GetUnlockEligibility returns whether the module has met the unlock criteria
	// (all prequizzes answered) regardless of video quizzes
	GetUnlockEligibility(userID, moduleID uint) (bool, error)
	// CalculateModuleProgress calculates the current progress percentage for a module
	CalculateModuleProgress(userID, moduleID uint) (float32, error)
}

type moduleProgressService struct {
	moduleProgressRepo repositories.ModuleProgressRepositoryInterface
	userRepo           repositories.UserRepositoryInterface
	moduleRepo         repositories.ModuleRepositoryInterface
	prequizRepo        repositories.PrequizRepositoryInterface
	videoQuizRepo      repositories.VideoQuizRepositoryInterface
}

func NewModuleProgressService(
	moduleProgressRepo repositories.ModuleProgressRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
	moduleRepo repositories.ModuleRepositoryInterface,
	prequizRepo repositories.PrequizRepositoryInterface,
	videoQuizRepo repositories.VideoQuizRepositoryInterface,
) *moduleProgressService {
	return &moduleProgressService{
		moduleProgressRepo: moduleProgressRepo,
		userRepo:           userRepo,
		moduleRepo:         moduleRepo,
		prequizRepo:        prequizRepo,
		videoQuizRepo:      videoQuizRepo,
	}
}

func (s *moduleProgressService) GetUserModuleProgress(userID, moduleID uint) (*models.ModuleProgress, error) {
	return s.moduleProgressRepo.GetByUserAndModule(userID, moduleID)
}

// CalculateModuleProgress calculates the current progress percentage for a module
func (s *moduleProgressService) CalculateModuleProgress(userID, moduleID uint) (float32, error) {
	// Get module with all its content
	module, err := s.moduleRepo.GetModuleByID(uint32(moduleID))
	if err != nil {
		return 0, err
	}

	// Get all prequizzes for this module
	allPrequizzes, err := s.prequizRepo.GetPrequizzesByModule(moduleID, 0) // 0 = no limit
	if err != nil {
		return 0, err
	}

	// If no prequizzes exist, return 0
	if len(allPrequizzes) == 0 {
		return 0, nil
	}

	// Get user's prequiz answers
	prequizAnswers, err := s.prequizRepo.GetUserPrequizAnswers(userID)
	if err != nil {
		return 0, err
	}

	// Count answered prequizzes
	answeredPrequizMap := make(map[uint]bool)
	for _, answer := range prequizAnswers {
		answeredPrequizMap[answer.PrequizID] = true
	}

	answeredPrequizzes := 0
	for _, prequiz := range allPrequizzes {
		if answeredPrequizMap[prequiz.ID] {
			answeredPrequizzes++
		}
	}

	// Check if all prequizzes are answered
	allPrequizzesAnswered := answeredPrequizzes == len(allPrequizzes)

	// Check video quizzes if they exist
	if module.VideoMaterial != nil && len(module.VideoMaterial.VideoQuizzes) > 0 {
		// There are video quizzes - check if they're all answered
		totalVideoQuizzes := len(module.VideoMaterial.VideoQuizzes)

		// Get user's video quiz answers
		videoAnswers, err := s.videoQuizRepo.GetAllUserVideoQuizAnswers(userID)
		if err != nil {
			return 0, err
		}

		// Count answered video quizzes for this module
		answeredVideoMap := make(map[uint]bool)
		for _, answer := range videoAnswers {
			answeredVideoMap[answer.VideoQuizID] = true
		}

		answeredVideoQuizzes := 0
		for _, videoQuiz := range module.VideoMaterial.VideoQuizzes {
			if answeredVideoMap[videoQuiz.ID] {
				answeredVideoQuizzes++
			}
		}

		// Check if all video quizzes are answered
		allVideoQuizzesAnswered := answeredVideoQuizzes == totalVideoQuizzes

		// Progress calculation with video quizzes:
		// - If all prequizzes AND all video quizzes are answered → 100%
		// - Otherwise, calculate based on combined progress
		if allPrequizzesAnswered && allVideoQuizzesAnswered {
			return 100.0, nil
		}

		// Calculate partial progress
		totalQuizzes := len(allPrequizzes) + totalVideoQuizzes
		answeredQuizzes := answeredPrequizzes + answeredVideoQuizzes
		return float32(answeredQuizzes) / float32(totalQuizzes) * 100.0, nil
	} else {
		// No video quizzes - progress depends only on prequizzes
		// If all prequizzes are answered → 100%
		if allPrequizzesAnswered {
			return 100.0, nil
		}

		// Calculate partial progress based on prequizzes only
		return float32(answeredPrequizzes) / float32(len(allPrequizzes)) * 100.0, nil
	}
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

		firstModuleID := modules[0].ID          // Assume modules are ordered by ID
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

	// Calculate current progress percentage
	newProgress, err := s.CalculateModuleProgress(userID, moduleID)
	if err != nil {
		return progress, err // Return existing progress if calculation fails
	}

	// Update progress values
	progress.Progress = newProgress

	// Check if module is now completed
	isCompleted, err := s.isModuleCompleted(userID, moduleID)
	if err == nil && isCompleted && !progress.IsComplete {
		progress.IsComplete = true
		progress.Progress = 100.0
		progress.MarkAsCompleted()
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

// CheckAndUnlockNextModule checks if current module is completed and unlocks the next module
func (s *moduleProgressService) CheckAndUnlockNextModule(userID, currentModuleID uint) error {
	// Unlock policy:
	// - Unlock next module when ALL PREQUIZZES are answered (video quizzes optional)
	// - Mark current module complete (is_complete=true, progress=100) only when
	//   both prequizzes AND (if any) video quizzes are answered.

	// 1) Check unlock eligibility (all prequizzes answered)
	canUnlock, err := s.GetUnlockEligibility(userID, currentModuleID)
	if err != nil {
		return err
	}
	if !canUnlock {
		return nil // not eligible to unlock next yet
	}

	// 2) Independently check full completion for marking complete/trigger
	isCompleted, _ := s.isModuleCompleted(userID, currentModuleID)

	// Mark current module as completed and update progress to 100%
	currentProgress, err := s.GetUserModuleProgress(userID, currentModuleID)
	if err == nil && currentProgress != nil {
		if isCompleted && !currentProgress.IsComplete {
			currentProgress.IsComplete = true
			currentProgress.Progress = 100.0
			currentProgress.MarkAsCompleted()
			s.moduleProgressRepo.UpdateProgress(currentProgress)
		}
	}

	// Get all modules to find the next module
	modules, err := s.moduleRepo.GetAllModules()
	if err != nil {
		return err
	}

	// Find next module (by ID order)
	var nextModule *models.Module
	for i, module := range modules {
		if module.ID == currentModuleID && i+1 < len(modules) {
			nextModule = &modules[i+1]
			break
		}
	}

	if nextModule == nil {
		return nil // No next module to unlock
	}

	// Safely unlock the next module (create record if it doesn't exist)
	return s.safeUnlockModule(userID, nextModule.ID)
}

// safeUnlockModule unlocks a module without generating "record not found" errors
func (s *moduleProgressService) safeUnlockModule(userID, moduleID uint) error {
	// Try to get existing progress without logging errors
	progress, err := s.moduleProgressRepo.GetByUserAndModule(userID, moduleID)

	if err != nil {
		// Record doesn't exist, create new unlocked progress
		newProgress := &models.ModuleProgress{
			UserID:     userID,
			ModuleID:   moduleID,
			IsUnlocked: true,
			IsComplete: false,
			Progress:   0,
		}
		newProgress.MarkAsStarted()
		_, err = s.moduleProgressRepo.CreateModuleProgress(newProgress)
		return err
	}

	// Record exists, just unlock it if not already unlocked
	if !progress.IsUnlocked {
		progress.IsUnlocked = true
		if progress.StartedAt == nil {
			progress.MarkAsStarted()
		}
		_, err = s.moduleProgressRepo.UpdateProgress(progress)
	}

	return err
}

// isModuleCompleted checks if all prequizzes and video quizzes in a module are answered
func (s *moduleProgressService) isModuleCompleted(userID, moduleID uint) (bool, error) {
	// Get module with all its content
	module, err := s.moduleRepo.GetModuleByID(uint32(moduleID))
	if err != nil {
		return false, err
	}

	// Get all prequizzes for this module
	allPrequizzes, err := s.prequizRepo.GetPrequizzesByModule(moduleID, 0) // 0 = no limit
	if err != nil {
		return false, err
	}

	// If no prequizzes exist, module cannot be completed through quiz answers
	if len(allPrequizzes) == 0 {
		return false, nil
	}

	// Get user's prequiz answers
	prequizAnswers, err := s.prequizRepo.GetUserPrequizAnswers(userID)
	if err != nil {
		return false, err
	}

	// Create map of answered prequizzes
	answeredPrequizMap := make(map[uint]bool)
	for _, answer := range prequizAnswers {
		answeredPrequizMap[answer.PrequizID] = true
	}

	// Check if all prequizzes are answered
	for _, prequiz := range allPrequizzes {
		if !answeredPrequizMap[prequiz.ID] {
			return false, nil // Found unanswered prequiz
		}
	}

	// Check video quizzes ONLY if video material exists AND has video quizzes
	if module.VideoMaterial != nil && len(module.VideoMaterial.VideoQuizzes) > 0 {
		// Get user's video quiz answers
		videoAnswers, err := s.videoQuizRepo.GetAllUserVideoQuizAnswers(userID)
		if err != nil {
			return false, err
		}

		// Create map of answered video quizzes
		answeredVideoMap := make(map[uint]bool)
		for _, answer := range videoAnswers {
			answeredVideoMap[answer.VideoQuizID] = true
		}

		// Check if all video quizzes are answered
		for _, videoQuiz := range module.VideoMaterial.VideoQuizzes {
			if !answeredVideoMap[videoQuiz.ID] {
				return false, nil // Found unanswered video quiz
			}
		}
	}
	// If no video material or no video quizzes, we only need prequizzes to be completed

	return true, nil // All required quizzes completed
}

// GetUnlockEligibility returns true if all prequizzes for the module are answered by the user
// This is used to unlock the next module without requiring video quizzes to be completed.
func (s *moduleProgressService) GetUnlockEligibility(userID, moduleID uint) (bool, error) {
	// Get all prequizzes for this module
	allPrequizzes, err := s.prequizRepo.GetPrequizzesByModule(moduleID, 0)
	if err != nil {
		return false, err
	}
	if len(allPrequizzes) == 0 {
		// If no prequizzes, don't unlock based on prequizzes alone
		return false, nil
	}

	// Get user's prequiz answers
	prequizAnswers, err := s.prequizRepo.GetUserPrequizAnswers(userID)
	if err != nil {
		return false, err
	}

	answered := make(map[uint]bool)
	for _, a := range prequizAnswers {
		answered[a.PrequizID] = true
	}
	for _, pq := range allPrequizzes {
		if !answered[pq.ID] {
			return false, nil
		}
	}
	return true, nil
}
