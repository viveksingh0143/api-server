package repository

import "github.com/vamika-digital/wms-api-server/internal/app/customer/domain"

type CustomerRepository interface {
	Create(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(customerID int64) error
	GetById(customerID int64) (*domain.Customer, error)
	GetTotalCount(filter CustomerFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter CustomerFilterOptions) ([]*domain.Customer, error)
}

type CustomerFilterOptions struct {
	Name          string
	ContactPerson string
	Code          string
	Status        string
}
