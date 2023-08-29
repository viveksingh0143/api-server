package repository

import (
	"strconv"

	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
)

type StoreRepository interface {
	Create(store *domain.Store) error
	Update(store *domain.Store) error
	Delete(storeID int64) error
	GetById(storeID int64) (*domain.Store, error)
	GetTotalCount(filter StoreFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter StoreFilterOptions) ([]*domain.Store, error)
}

type StoreFilterOptions struct {
	Name     string
	Location string
	Status   string
	OwnerID  int64
}

func (f *StoreFilterOptions) SetOwnerID(ownerIDStr string) {
	if ownerIDStr != "" {
		ownerId, err := strconv.ParseInt(ownerIDStr, 10, 64)
		if err != nil {
			return
		}
		f.OwnerID = ownerId
	} else {
		f.OwnerID = 0
	}
}
