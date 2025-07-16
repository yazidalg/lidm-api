package services

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"github.com/yazidalg/lidm_backend/internal/app/repositories"
)

type RoleServiceInterface interface {
	GetRoleByName(name string) (models.Role, error)
	GetRoleById(id uint) (models.Role, error)
	GetAllRoles() ([]models.Role, error)
	CreateRole(role models.Role) (models.Role, error)
	UpdateRole(role models.Role) (models.Role, error)
	DeleteRole(id uint) error
}

type roleService struct {
	roleRepository repositories.RoleRepositoryInterface
}

func NewRoleService(roleRepository repositories.RoleRepositoryInterface) *roleService {
	return &roleService{roleRepository}
}

func (s *roleService) GetRoleByName(name string) (models.Role, error) {
	return s.roleRepository.GetByName(name)
}

func (s *roleService) GetRoleById(id uint) (models.Role, error) {
	return s.roleRepository.GetById(id)
}

func (s *roleService) GetAllRoles() ([]models.Role, error) {
	return s.roleRepository.GetAllRoles()
}

func (s *roleService) CreateRole(role models.Role) (models.Role, error) {
	return s.roleRepository.CreateRole(role)
}

func (s *roleService) UpdateRole(role models.Role) (models.Role, error) {
	return s.roleRepository.UpdateRole(role)
}

func (s *roleService) DeleteRole(id uint) error {
	return s.roleRepository.DeleteRole(id)
}
