package usecase

import "github.com/vamika-digital/wms-api-server/internal/app/user/domain"

type RoleUseCase interface {
	CreateRole(role *domain.Role) error
	UpdateRole(role *domain.Role) error
	DeleteRole(roleID int) error
	GetRoleByID(roleID int) (*domain.Role, error)
	GetAllRoles() ([]*domain.Role, error)
}
