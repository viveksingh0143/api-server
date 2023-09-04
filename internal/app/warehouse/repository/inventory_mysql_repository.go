package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlInventoryRepository struct {
	conn database.Connection
}

func NewInventoryRepository(conn database.Connection) InventoryRepository {
	return &MySqlInventoryRepository{conn: conn}
}

func (r *MySqlInventoryRepository) Create(inventory *domain.Inventory) error {
	query := `INSERT INTO inventories (status, pallet_id, bin_id, rack_id, store_id, product_id, batch, machine, shift, supervisor, quantity, unit, stockin_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.conn.GetDB().Exec(query, inventory.Status, inventory.PalletID, inventory.BinID, inventory.RackID, inventory.StoreID, inventory.ProductID, inventory.Batch, inventory.Machine, inventory.Shift, inventory.Supervisor, inventory.Quantity, inventory.Unit, inventory.StockInAt)
	return err
}

func (r *MySqlInventoryRepository) Update(inventory *domain.Inventory) error {
	query := `UPDATE inventories SET status=?, pallet_id=?, bin_id=?, rack_id=?, store_id=?, product_id=?, batch=?, machine=?, shift=?, supervisor=?, quantity=?, unit=?, stockout_at=? WHERE id=?`
	_, err := r.conn.GetDB().Exec(query, inventory.Status, inventory.PalletID, inventory.BinID, inventory.RackID, inventory.StoreID, inventory.ProductID, inventory.Batch, inventory.Machine, inventory.Shift, inventory.Supervisor, inventory.Quantity, inventory.Unit, inventory.StockOutAt, inventory.ID)
	return err
}

func (r *MySqlInventoryRepository) Delete(inventoryID int64) error {
	query := "DELETE FROM inventories WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, inventoryID)
	return err
}

func (r *MySqlInventoryRepository) GetById(inventoryID int64) (*domain.Inventory, error) {
	query := `SELECT id, status, pallet_id, bin_id, rack_id, store_id, product_id, batch, machine, shift, supervisor, quantity, unit, stockin_at, stockout_at FROM inventories WHERE id = ?`
	row := r.conn.GetDB().QueryRow(query, inventoryID)
	inventory := domain.NewInventoryWithDefaults()
	err := row.Scan(&inventory.ID, &inventory.Status, &inventory.PalletID, &inventory.BinID, &inventory.RackID, &inventory.StoreID, &inventory.ProductID, &inventory.Batch, &inventory.Machine, &inventory.Shift, &inventory.Supervisor, &inventory.Quantity, &inventory.Unit, &inventory.StockInAt, &inventory.StockOutAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return inventory, nil
}

func (r *MySqlInventoryRepository) GetTotalCount(filter InventoryFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM inventories", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlInventoryRepository) GetAll(page int, pageSize int, sort string, filter InventoryFilterOptions) ([]*domain.Inventory, error) {
	query, args := r.buildFilterQuery("SELECT id, status, pallet_id, bin_id, rack_id, store_id, product_id, batch, machine, shift, supervisor, quantity, unit, stockin_at, stockout_at FROM inventories", filter)
	var allowedSortOrders = map[string]bool{
		"stockin_at ASC":   true,
		"stockin_at DESC":  true,
		"stockin_out ASC":  true,
		"stockin_out DESC": true,
		"status ASC":       true,
		"status DESC":      true,
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

	var inventories []*domain.Inventory
	for rows.Next() {
		inventory := domain.NewInventoryWithDefaults()
		if err := rows.Scan(&inventory.ID, &inventory.Status, &inventory.PalletID, &inventory.BinID, &inventory.RackID, &inventory.StoreID, &inventory.ProductID, &inventory.Batch, &inventory.Machine, &inventory.Shift, &inventory.Supervisor, &inventory.Quantity, &inventory.Unit, &inventory.StockInAt, &inventory.StockOutAt); err != nil {
			return nil, err
		}
		inventories = append(inventories, inventory)
	}
	return inventories, nil
}

func (r *MySqlInventoryRepository) buildFilterQuery(baseQuery string, filter InventoryFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Status != "" {
		filters = append(filters, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.ProductID > 0 {
		filters = append(filters, "product_id = ?")
		args = append(args, filter.ProductID)
	}
	if filter.BinID > 0 {
		filters = append(filters, "bin_id = ?")
		args = append(args, filter.BinID)
	}
	if filter.RackID > 0 {
		filters = append(filters, "rack_id = ?")
		args = append(args, filter.RackID)
	}
	if filter.StoreID > 0 {
		filters = append(filters, "store_id = ?")
		args = append(args, filter.StoreID)
	}

	if len(filters) > 0 {
		return baseQuery + " WHERE " + strings.Join(filters, " AND "), args
	}

	return baseQuery, args
}
