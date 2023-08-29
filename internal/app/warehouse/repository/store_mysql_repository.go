package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlStoreRepository struct {
	conn database.Connection
}

func NewStoreRepository(conn database.Connection) StoreRepository {
	return &MySqlStoreRepository{conn: conn}
}

func (r *MySqlStoreRepository) Create(store *domain.Store) error {
	query := "INSERT INTO stores (name, location, status, last_updated_by, owner_id) VALUES (?, ?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, store.Name, store.Location, store.Status, store.LastUpdatedBy, store.Owner.ID)
	return err
}

func (r *MySqlStoreRepository) Update(store *domain.Store) error {
	query := "UPDATE stores SET name=?, location=?, status=?, last_updated_by=?, owner_id=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, store.Name, store.Location, store.Status, store.LastUpdatedBy, store.Owner.ID, store.ID)
	return err
}

func (r *MySqlStoreRepository) Delete(storeID int64) error {
	query := "DELETE FROM stores WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, storeID)
	return err
}

func (r *MySqlStoreRepository) GetById(storeID int64) (*domain.Store, error) {
	query := "SELECT id, name, location, status, created_at, updated_at, last_updated_by, owner_id FROM stores WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, storeID)
	store := domain.NewStoreWithDefaults()
	err := row.Scan(&store.ID, &store.Name, &store.Location, &store.Status, &store.CreatedAt, &store.UpdatedAt, &store.LastUpdatedBy, &store.Owner.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	if store.Owner.ID <= 0 {
		store.Owner = nil
	}
	return store, nil
}

func (r *MySqlStoreRepository) GetTotalCount(filter StoreFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM stores", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlStoreRepository) GetAll(page int, pageSize int, sort string, filter StoreFilterOptions) ([]*domain.Store, error) {
	query, args := r.buildFilterQuery("SELECT id, name, location, status, created_at, updated_at, last_updated_by, owner_id FROM stores", filter)
	var allowedSortOrders = map[string]bool{
		"location ASC":  true,
		"location DESC": true,
		"name ASC":      true,
		"name DESC":     true,
		"code ASC":      true,
		"code DESC":     true,
		"status ASC":    true,
		"status DESC":   true,
		"owner_id ASC":  true,
		"owner_id DESC": true,
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

	var stores []*domain.Store
	for rows.Next() {
		store := domain.NewStoreWithDefaults()
		if err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Status, &store.CreatedAt, &store.UpdatedAt, &store.LastUpdatedBy, &store.Owner.ID); err != nil {
			return nil, err
		}
		if store.Owner.ID <= 0 {
			store.Owner = nil
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (r *MySqlStoreRepository) buildFilterQuery(baseQuery string, filter StoreFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Name != "" {
		filters = append(filters, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.Location != "" {
		filters = append(filters, "location LIKE ?")
		args = append(args, "%"+filter.Location+"%")
	}
	if filter.Status != "" {
		filters = append(filters, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.OwnerID > 0 {
		filters = append(filters, "owner_id = ?")
		args = append(args, filter.OwnerID)
	}

	if len(filters) > 0 {
		return baseQuery + " WHERE " + strings.Join(filters, " AND "), args
	}

	return baseQuery, args
}
