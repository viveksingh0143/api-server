package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/customer/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/repository"
)

type CustomerUseCaseImpl struct {
	Repo repository.CustomerRepository
}

func NewCustomerUseCase(repo repository.CustomerRepository) CustomerUseCase {
	return &CustomerUseCaseImpl{Repo: repo}
}

func (u *CustomerUseCaseImpl) CreateCustomer(customer *domain.Customer) error {
	return u.Repo.Create(customer)
}

func (u *CustomerUseCaseImpl) UpdateCustomer(customer *domain.Customer) error {
	_, err := u.Repo.GetById(customer.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(customer)
}

func (u *CustomerUseCaseImpl) DeleteCustomer(customerID int64) error {
	return u.Repo.Delete(customerID)
}

func (u *CustomerUseCaseImpl) GetCustomerByID(customerID int64) (*domain.Customer, error) {
	return u.Repo.GetById(customerID)
}

func (u *CustomerUseCaseImpl) GetAllCustomers(page int, pageSize int, sort string, filter repository.CustomerFilterOptions) ([]*domain.Customer, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	customers, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of customers matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
