package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/product/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/product/repository"
)

type ProductUseCaseImpl struct {
	Repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &ProductUseCaseImpl{Repo: repo}
}

func (u *ProductUseCaseImpl) CreateProduct(product *domain.Product) error {
	return u.Repo.Create(product)
}

func (u *ProductUseCaseImpl) UpdateProduct(product *domain.Product) error {
	// Check for an existing product with the specified ID
	// existingProduct, err := u.Repo.GetById(product.ID)
	_, err := u.Repo.GetById(product.ID)
	if err != nil {
		return err
	}
	return u.Repo.Update(product)
}

func (u *ProductUseCaseImpl) DeleteProduct(productID int64) error {
	return u.Repo.Delete(productID)
}

func (u *ProductUseCaseImpl) GetProductByID(productID int64) (*domain.Product, error) {
	return u.Repo.GetById(productID)
}

func (u *ProductUseCaseImpl) GetAllProducts(page int, pageSize int, sort string, filter repository.ProductFilterOptions) ([]*domain.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of products matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
