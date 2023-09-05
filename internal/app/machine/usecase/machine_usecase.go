package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/machine/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/repository"
)

type MachineUseCase interface {
	CreateMachine(store *domain.Machine) error
	UpdateMachine(store *domain.Machine) error
	DeleteMachine(storeID int64) error
	GetMachineByID(storeID int64) (*domain.Machine, error)
	GetAllMachines(page int, pageSize int, sort string, filter repository.MachineFilterOptions) ([]*domain.Machine, int, error)
}
