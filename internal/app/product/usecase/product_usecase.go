package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/product/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/product/repository"
)

type ProductUseCase interface {
	CreateProduct(product *domain.Product) error
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productID int64) error
	GetProductByID(productID int64) (*domain.Product, error)
	GetAllProducts(page int, pageSize int, sort string, filter repository.ProductFilterOptions) ([]*domain.Product, int, error)
}
