package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
	"github.com/yazidalg/lidm_backend/internal/app/request"
)

type ModuleServiceInterface interface {
	CreateModule(request request.ModuleRequest) (*models.Module, error)
	GetModuleByID(id uint32) (*models.Module, error)
	GetAllModules() ([]models.Module, error)
	GetAllModulesWithProgress(userID uint) ([]map[string]interface{}, error)
	UpdateModule(id uint32, request request.ModuleRequest) (*models.Module, error)
	DeleteModule(id uint32) error
}

type moduleService struct {
	moduleRepository     repositories.ModuleRepositoryInterface
	progressRepo         repositories.ProgressRepositoryInterface
	videoQuizRepo        repositories.VideoQuizRepositoryInterface
	prequizRepo          repositories.PrequizRepositoryInterface
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
			"ID":          module.ID,
			"Title":       module.Title,
			"Description": module.Description,
			"Thumbnail":   module.Thumbnail,
			"Icon":        module.Icon,
			"OffsetX":     module.OffsetX,
			"OffsetY":     module.OffsetY,
				"CreatedAt":   module.CreatedAt,
				"UpdatedAt":   module.UpdatedAt,
		}

		// Calculate module progress
		totalVideoQuizzes := 0
		answeredVideoQuizzes := 0
		totalPrequizzes := 0
		answeredPrequizzes := 0
		subMaterialsWithStatus := make([]map[string]interface{}, 0)

	for _, sm := range module.SubMaterials {
			// FORCE CLEAR any existing prequizzes data
			sm.Prequizzes = nil
			
			smData := map[string]interface{}{
				"ID":          sm.ID,
				"Title":       sm.Title,
				"Description": sm.Description,
				"Order":       sm.Order,
				"CreatedAt":   sm.CreatedAt,
				"UpdatedAt":   sm.UpdatedAt,
			}

			// Video quizzes status
			videoQuizCount := 0
			videoQuizAnswered := 0
			if sm.VideoMaterial != nil {
				videoQuizCount = len(sm.VideoMaterial.VideoQuizzes)
				totalVideoQuizzes += videoQuizCount

				// Check answered video quizzes
				for _, vq := range sm.VideoMaterial.VideoQuizzes {
					hasAnswered, _ := s.videoQuizRepo.HasUserAnsweredVideoQuiz(userID, vq.ID)
					if hasAnswered {
						videoQuizAnswered++
						answeredVideoQuizzes++
					}
				}
				
				smData["video_material"] = map[string]interface{}{
					"ID":           sm.VideoMaterial.ID,
					"Title":        sm.VideoMaterial.Title,
					"YoutubeLink":  sm.VideoMaterial.YoutubeLink,
					"Duration":     sm.VideoMaterial.Duration,
					"video_quizzes": sm.VideoMaterial.VideoQuizzes,
				}
			}

	    // Prequizzes status (limit to 3)
	    limitedPrequizzes, _ := s.prequizRepo.GetPrequizzesBySubMaterial(sm.ID, 3)
			prequizCount := len(limitedPrequizzes)
			prequizAnswered := 0
			totalPrequizzes += prequizCount

	    // answeredPrequizMap already built above

			// COMPLETELY CUSTOM prequizzes - ignore GORM models
			customPrequizzes := []map[string]interface{}{}
			for _, pq := range limitedPrequizzes {
				isAnswered, exists := answeredPrequizMap[pq.ID]
				if !exists {
					isAnswered = false
				}
				if isAnswered {
					prequizAnswered++
					answeredPrequizzes++
				}
				
					customPrequiz := map[string]interface{}{
						"ID":               pq.ID,
						"SubMaterialID":    pq.SubMaterialID,
						"Question":         pq.Question,
						"Options":          pq.Options,
						"CorrectAnswer":    pq.CorrectAnswer,
						"Explanation":      pq.Explanation,
						"isAlreadyAnswered": isAnswered,
					}
				customPrequizzes = append(customPrequizzes, customPrequiz)
			}

			smData["prequizzes"] = customPrequizzes
			smData["ar_experiment"] = sm.ARExperiment
			smData["video_quiz_status"] = map[string]interface{}{
				"total":     videoQuizCount,
				"answered":  videoQuizAnswered,
				"completed": videoQuizCount > 0 && videoQuizAnswered == videoQuizCount,
			}
			smData["prequiz_status"] = map[string]interface{}{
				"total":     prequizCount,
				"answered":  prequizAnswered,
				"completed": prequizCount > 0 && prequizAnswered == prequizCount,
			}

			subMaterialsWithStatus = append(subMaterialsWithStatus, smData)
		}

		moduleData["sub_materials"] = subMaterialsWithStatus
		
		// Calculate module status - true if all sub-materials are completed
		allSubMaterialsCompleted := true
		for _, sm := range subMaterialsWithStatus {
			videoStatus := sm["video_quiz_status"].(map[string]interface{})
			prequizStatus := sm["prequiz_status"].(map[string]interface{})
			
			videoCompleted := videoStatus["completed"].(bool)
			prequizCompleted := prequizStatus["completed"].(bool)
			
			// Sub-material is completed if both video quizzes and prequizzes are completed
			// (or if no video quizzes/prequizzes exist for that sub-material)
			if !videoCompleted || !prequizCompleted {
				allSubMaterialsCompleted = false
				break
			}
		}
		
		moduleData["module_status"] = allSubMaterialsCompleted

		result = append(result, moduleData)
	}

	return result, nil
}

func (s *moduleService) UpdateModule(id uint32, request request.ModuleRequest) (*models.Module, error) {
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