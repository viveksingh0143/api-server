package repository

import "github.com/vamika-digital/wms-api-server/internal/app/user/domain"

type RoleRepository interface {
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(roleID int) error
	GetById(roleID int) (*domain.Role, error)
	GetAll() ([]*domain.Role, error)
}
