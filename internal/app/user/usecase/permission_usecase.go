package usecase

import "github.com/vamika-digital/wms-api-server/internal/app/user/domain"

type PermissionUseCase interface {
	CreatePermission(permission *domain.Permission) error
	UpdatePermission(permission *domain.Permission) error
	DeletePermission(permissionID int) error
	GetPermissionByID(permissionID int) (*domain.Permission, error)
	GetAllPermissions() ([]*domain.Permission, error)
}
