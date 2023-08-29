package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type StoreUseCase interface {
	CreateStore(store *domain.Store) error
	UpdateStore(store *domain.Store) error
	DeleteStore(storeID int64) error
	GetStoreByID(storeID int64) (*domain.Store, error)
	GetAllStores(page int, pageSize int, sort string, filter repository.StoreFilterOptions) ([]*domain.Store, int, error)
}
