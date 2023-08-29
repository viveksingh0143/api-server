package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/product/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlProductRepository struct {
	conn database.Connection
}

func NewProductRepository(conn database.Connection) ProductRepository {
	return &MySqlProductRepository{conn: conn}
}

func (r *MySqlProductRepository) Create(product *domain.Product) error {
	query := "INSERT INTO products (type, code, raw_code, name, description, unit, status, last_updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, product.Type, product.Code, product.RawCode, product.Name, product.Description, product.Unit, product.Status, product.LastUpdatedBy)
	return err
}

func (r *MySqlProductRepository) Update(product *domain.Product) error {
	query := "UPDATE products SET type=?, code=?, raw_code=?, name=?, description=?, unit=?, status=?, last_updated_by=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, product.Type, product.Code, product.RawCode, product.Name, product.Description, product.Unit, product.Status, product.LastUpdatedBy, product.ID)
	return err
}

func (r *MySqlProductRepository) Delete(productID int64) error {
	query := "DELETE FROM products WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, productID)
	return err
}

func (r *MySqlProductRepository) GetById(productID int64) (*domain.Product, error) {
	query := "SELECT id, type, code, raw_code, name, description, unit, status, created_at, updated_at, last_updated_by FROM products WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, productID)
	product := &domain.Product{}
	err := row.Scan(&product.ID, &product.Type, &product.Code, &product.RawCode, &product.Name, &product.Description, &product.Unit, &product.Status, &product.CreatedAt, &product.UpdatedAt, &product.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *MySqlProductRepository) GetTotalCount(filter ProductFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM products", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlProductRepository) GetAll(page int, pageSize int, sort string, filter ProductFilterOptions) ([]*domain.Product, error) {
	query, args := r.buildFilterQuery("SELECT id, type, code, raw_code, name, description, unit, status, created_at, updated_at, last_updated_by FROM products", filter)
	var allowedSortOrders = map[string]bool{
		"name ASC":        true,
		"name DESC":       true,
		"code ASC":        true,
		"code DESC":       true,
		"status ASC":      true,
		"status DESC":     true,
		"updated_at ASC":  true,
		"updated_at DESC": true,
	}

	if sort != "" {
		if _, ok := allowedSortOrders[sort]; ok {
			query += " ORDER BY " + sort
		} else {
			return nil, errors.New("invalid sort order")
		}
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, page)

	rows, err := r.conn.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		if err := rows.Scan(&product.ID, &product.Type, &product.Code, &product.RawCode, &product.Name, &product.Description, &product.Unit, &product.Status, &product.CreatedAt, &product.UpdatedAt, &product.LastUpdatedBy); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *MySqlProductRepository) buildFilterQuery(baseQuery string, filter ProductFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Type != "" {
		filters = append(filters, "type = ?")
		args = append(args, filter.Type)
	}
	if filter.Code != "" {
		filters = append(filters, "code LIKE ?")
		args = append(args, "%"+filter.Code+"%")
	}
	if filter.RawCode != "" {
		filters = append(filters, "raw_code LIKE ?")
		args = append(args, "%"+filter.RawCode+"%")
	}
	if filter.Name != "" {
		filters = append(filters, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		filters = append(filters, "status = ?")
		args = append(args, filter.Status)
	}

	if len(filters) > 0 {
		return baseQuery + " WHERE " + strings.Join(filters, " AND "), args
	}

	return baseQuery, args
}
