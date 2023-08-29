package repository

import "github.com/vamika-digital/wms-api-server/internal/app/product/domain"

type ProductRepository interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(productID int64) error
	GetById(productID int64) (*domain.Product, error)
	GetTotalCount(filter ProductFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter ProductFilterOptions) ([]*domain.Product, error)
}

type ProductFilterOptions struct {
	Type    string
	Code    string
	RawCode string
	Name    string
	Status  string
}
