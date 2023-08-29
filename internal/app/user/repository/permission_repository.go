package repository

import "github.com/vamika-digital/wms-api-server/internal/app/user/domain"

type PermissionRepository interface {
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	Delete(permissionID int) error
	GetById(permissionID int) (*domain.Permission, error)
	GetAll() ([]*domain.Permission, error)
}
