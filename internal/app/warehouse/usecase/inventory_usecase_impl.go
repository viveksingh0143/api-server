package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type InventoryUseCaseImpl struct {
	Repo          repository.InventoryRepository
	RepoContainer repository.ContainerRepository
}

func NewInventoryUseCase(repo repository.InventoryRepository) InventoryUseCase {
	return &InventoryUseCaseImpl{Repo: repo}
}

func (u *InventoryUseCaseImpl) CreateInventoryForRawMaterial(inventory *domain.InventoryFormRawMaterial) error {
	// palletContainer, err := u.RepoContainer.GetByCode(inventory.Pallet)
	// if err != nil {
	// 	palletContainer, err = u.RepoContainer.Create(&domain.Container{
	// 		Type:   domain.PALLET_TYPE,
	// 		Status: "active",
	// 		Code:   customtypes.NullableString(inventory.Pallet),
	// 		Name:   customtypes.NullableString(inventory.Pallet),
	// 	})
	// }
	// return u.Repo.Create(inventory)
	return nil
}

func (u *InventoryUseCaseImpl) CreateInventory(inventory *domain.Inventory) error {
	return u.Repo.Create(inventory)
}

func (u *InventoryUseCaseImpl) UpdateInventory(inventory *domain.Inventory) error {
	// Check for an existing inventory with the specified ID
	// existingInventory, err := u.Repo.GetById(inventory.ID)
	_, err := u.Repo.GetById(inventory.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(inventory)
}

func (u *InventoryUseCaseImpl) DeleteInventory(inventoryID int64) error {
	return u.Repo.Delete(inventoryID)
}

func (u *InventoryUseCaseImpl) GetInventoryByID(inventoryID int64) (*domain.Inventory, error) {
	return u.Repo.GetById(inventoryID)
}

func (u *InventoryUseCaseImpl) GetAllInventories(page int, pageSize int, sort string, filter repository.InventoryFilterOptions) ([]*domain.Inventory, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	inventories, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of inventories matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}
