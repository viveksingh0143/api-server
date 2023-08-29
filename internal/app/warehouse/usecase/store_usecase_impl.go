package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
)

type StoreUseCaseImpl struct {
	Repo repository.StoreRepository
}

func NewStoreUseCase(repo repository.StoreRepository) StoreUseCase {
	return &StoreUseCaseImpl{Repo: repo}
}

func (u *StoreUseCaseImpl) CreateStore(store *domain.Store) error {
	return u.Repo.Create(store)
}

func (u *StoreUseCaseImpl) UpdateStore(store *domain.Store) error {
	// Check for an existing store with the specified ID
	// existingStore, err := u.Repo.GetById(store.ID)
	_, err := u.Repo.GetById(store.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(store)
}

func (u *StoreUseCaseImpl) DeleteStore(storeID int64) error {
	return u.Repo.Delete(storeID)
}

func (u *StoreUseCaseImpl) GetStoreByID(storeID int64) (*domain.Store, error) {
	return u.Repo.GetById(storeID)
}

func (u *StoreUseCaseImpl) GetAllStores(page int, pageSize int, sort string, filter repository.StoreFilterOptions) ([]*domain.Store, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	stores, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of stores matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}
