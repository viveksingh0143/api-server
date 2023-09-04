package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type InventoryUseCase interface {
	CreateInventoryForRawMaterial(inventory *domain.InventoryFormRawMaterial) error
	CreateInventory(inventory *domain.Inventory) error
	UpdateInventory(inventory *domain.Inventory) error
	DeleteInventory(inventoryID int64) error
	GetInventoryByID(inventoryID int64) (*domain.Inventory, error)
	GetAllInventories(page int, pageSize int, sort string, filter repository.InventoryFilterOptions) ([]*domain.Inventory, int, error)
}
