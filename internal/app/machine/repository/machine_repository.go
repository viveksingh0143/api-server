package repository

import (
	"github.com/vamika-digital/wms-api-server/internal/app/machine/domain"
)

type MachineRepository interface {
	Create(machine *domain.Machine) error
	Update(machine *domain.Machine) error
	Delete(machineID int64) error
	GetById(machineID int64) (*domain.Machine, error)
	GetTotalCount(filter MachineFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter MachineFilterOptions) ([]*domain.Machine, error)
}

type MachineFilterOptions struct {
	Name   string
	Code   string
	Status string
}
