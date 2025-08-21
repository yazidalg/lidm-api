package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type SubMaterialProgressServiceInterface interface {
	CreateOrUpdateProgress(userID uint, subMaterialID uint, videoCompleted, prequizzesCompleted bool) (*models.SubMaterialProgress, error)
	GetUserProgressForModule(userID uint, moduleID uint) ([]models.SubMaterialProgress, error)
	MarkVideoCompleted(userID uint, subMaterialID uint) error
	MarkPrequizzesCompleted(userID uint, subMaterialID uint) error
	GetSubMaterialProgress(userID uint, subMaterialID uint) (*models.SubMaterialProgress, error)
}

type subMaterialProgressService struct {
	repo repositories.SubMaterialProgressRepositoryInterface
}

func NewSubMaterialProgressService(repo repositories.SubMaterialProgressRepositoryInterface) SubMaterialProgressServiceInterface {
	return &subMaterialProgressService{repo: repo}
}

func (s *subMaterialProgressService) CreateOrUpdateProgress(userID uint, subMaterialID uint, videoCompleted, prequizzesCompleted bool) (*models.SubMaterialProgress, error) {
	progress := &models.SubMaterialProgress{
		UserID:              userID,
		SubMaterialID:       subMaterialID,
		VideoCompleted:      videoCompleted,
		PrequizzesCompleted: prequizzesCompleted,
		Completed:           videoCompleted && prequizzesCompleted,
	}
	
	return s.repo.CreateOrUpdateProgress(progress)
}

func (s *subMaterialProgressService) GetUserProgressForModule(userID uint, moduleID uint) ([]models.SubMaterialProgress, error) {
	return s.repo.GetByUserAndModule(userID, moduleID)
}

func (s *subMaterialProgressService) MarkVideoCompleted(userID uint, subMaterialID uint) error {
	return s.repo.UpdateVideoCompleted(userID, subMaterialID, true)
}

func (s *subMaterialProgressService) MarkPrequizzesCompleted(userID uint, subMaterialID uint) error {
	return s.repo.UpdatePrequizzesCompleted(userID, subMaterialID, true)
}

func (s *subMaterialProgressService) GetSubMaterialProgress(userID uint, subMaterialID uint) (*models.SubMaterialProgress, error) {
	return s.repo.GetByUserAndSubMaterial(userID, subMaterialID)
}
