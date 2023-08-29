package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type ContainerUseCase interface {
	CreateContainer(container *domain.Container) error
	UpdateContainer(container *domain.Container) error
	DeleteContainer(containerID int64) error
	GetContainerByID(containerID int64) (*domain.Container, error)
	GetAllContainers(page int, pageSize int, sort string, filter repository.ContainerFilterOptions) ([]*domain.Container, int, error)
}
