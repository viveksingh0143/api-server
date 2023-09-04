package repository

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customtypes"
)

type InventoryRepository interface {
	Create(inventory *domain.Inventory) error
	Update(inventory *domain.Inventory) error
	Delete(inventoryID int64) error
	GetById(inventoryID int64) (*domain.Inventory, error)
	GetTotalCount(filter InventoryFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter InventoryFilterOptions) ([]*domain.Inventory, error)
}

type InventoryFilterOptions struct {
	Status    domain.InventoryType
	ProductID customtypes.NullableInt64
	BinID     customtypes.NullableInt64
	RackID    customtypes.NullableInt64
	StoreID   customtypes.NullableInt64
}
