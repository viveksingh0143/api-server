package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/customer/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/repository"
)

type CustomerUseCase interface {
	CreateCustomer(store *domain.Customer) error
	UpdateCustomer(store *domain.Customer) error
	DeleteCustomer(storeID int64) error
	GetCustomerByID(storeID int64) (*domain.Customer, error)
	GetAllCustomers(page int, pageSize int, sort string, filter repository.CustomerFilterOptions) ([]*domain.Customer, int, error)
}
