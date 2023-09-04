package repository

import "github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"

type ContainerRepository interface {
	Create(container *domain.Container) error
	Update(container *domain.Container) error
	Delete(containerID int64) error
	GetByCode(code string) (*domain.Container, error)
	GetById(containerID int64) (*domain.Container, error)
	GetTotalCount(filter ContainerFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter ContainerFilterOptions) ([]*domain.Container, error)
}

type ContainerFilterOptions struct {
	Type   string
	Code   string
	Name   string
	Status string
}
