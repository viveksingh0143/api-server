package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/machine/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/repository"
)

type MachineUseCaseImpl struct {
	Repo repository.MachineRepository
}

func NewMachineUseCase(repo repository.MachineRepository) MachineUseCase {
	return &MachineUseCaseImpl{Repo: repo}
}

func (u *MachineUseCaseImpl) CreateMachine(machine *domain.Machine) error {
	return u.Repo.Create(machine)
}

func (u *MachineUseCaseImpl) UpdateMachine(machine *domain.Machine) error {
	_, err := u.Repo.GetById(machine.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(machine)
}

func (u *MachineUseCaseImpl) DeleteMachine(machineID int64) error {
	return u.Repo.Delete(machineID)
}

func (u *MachineUseCaseImpl) GetMachineByID(machineID int64) (*domain.Machine, error) {
	return u.Repo.GetById(machineID)
}

func (u *MachineUseCaseImpl) GetAllMachines(page int, pageSize int, sort string, filter repository.MachineFilterOptions) ([]*domain.Machine, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	machines, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of machines matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return machines, total, nil
}
