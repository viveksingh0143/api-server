package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type ContainerUseCaseImpl struct {
	Repo repository.ContainerRepository
}

func NewContainerUseCase(repo repository.ContainerRepository) ContainerUseCase {
	return &ContainerUseCaseImpl{Repo: repo}
}

func (u *ContainerUseCaseImpl) CreateContainer(container *domain.Container) error {
	return u.Repo.Create(container)
}

func (u *ContainerUseCaseImpl) UpdateContainer(container *domain.Container) error {
	// Check for an existing container with the specified ID
	// existingContainer, err := u.Repo.GetById(container.ID)
	_, err := u.Repo.GetById(container.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(container)
}

func (u *ContainerUseCaseImpl) DeleteContainer(containerID int64) error {
	return u.Repo.Delete(containerID)
}

func (u *ContainerUseCaseImpl) GetContainerByID(containerID int64) (*domain.Container, error) {
	return u.Repo.GetById(containerID)
}

func (u *ContainerUseCaseImpl) GetAllContainers(page int, pageSize int, sort string, filter repository.ContainerFilterOptions) ([]*domain.Container, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	containers, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of containers matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return containers, total, nil
}
