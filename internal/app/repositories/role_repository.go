package repositories

import (
	"github.com/yazidalg/lidm_backend/internal/app/models"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetByName(name string) (models.Role, error)
	GetById(id uint) (models.Role, error)
	GetAllRoles() ([]models.Role, error)
	CreateRole(role models.Role) (models.Role, error)
	UpdateRole(role models.Role) (models.Role, error)
	DeleteRole(id uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) GetByName(name string) (models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return role, err
}

func (r *roleRepository) GetById(id uint) (models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	return role, err
}

func (r *roleRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) CreateRole(role models.Role) (models.Role, error) {
	err := r.db.Create(&role).Error
	return role, err
}

func (r *roleRepository) UpdateRole(role models.Role) (models.Role, error) {
	err := r.db.Save(&role).Error
	return role, err
}

func (r *roleRepository) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
