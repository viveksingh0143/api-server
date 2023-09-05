package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/machine/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlMachineRepository struct {
	conn database.Connection
}

func NewMachineRepository(conn database.Connection) MachineRepository {
	return &MySqlMachineRepository{conn: conn}
}

func (r *MySqlMachineRepository) Create(machine *domain.Machine) error {
	query := "INSERT INTO machines (name, code, status, last_updated_by) VALUES (?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, machine.Name, machine.Code, machine.Status, machine.LastUpdatedBy)
	return err
}

func (r *MySqlMachineRepository) Update(machine *domain.Machine) error {
	query := "UPDATE machines SET name=?, code=?, status=?, last_updated_by=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, machine.Name, machine.Code, machine.Status, machine.LastUpdatedBy, machine.ID)
	return err
}

func (r *MySqlMachineRepository) Delete(machineID int64) error {
	query := "DELETE FROM machines WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, machineID)
	return err
}

func (r *MySqlMachineRepository) GetById(machineID int64) (*domain.Machine, error) {
	query := "SELECT id, name, code, status, created_at, updated_at, last_updated_by FROM machines WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, machineID)
	machine := domain.NewMachineWithDefaults()
	err := row.Scan(&machine.ID, &machine.Name, &machine.Code, &machine.Status, &machine.CreatedAt, &machine.UpdatedAt, &machine.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return machine, nil
}

func (r *MySqlMachineRepository) GetTotalCount(filter MachineFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM machines", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlMachineRepository) GetAll(page int, pageSize int, sort string, filter MachineFilterOptions) ([]*domain.Machine, error) {
	query, args := r.buildFilterQuery("SELECT id, name, code, status, created_at, updated_at, last_updated_by FROM machines", filter)
	var allowedSortOrders = map[string]bool{
		"code ASC":    true,
		"code DESC":   true,
		"name ASC":    true,
		"name DESC":   true,
		"status ASC":  true,
		"status DESC": true,
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

	var machines []*domain.Machine
	for rows.Next() {
		machine := domain.NewMachineWithDefaults()
		if err := rows.Scan(&machine.ID, &machine.Name, &machine.Code, &machine.Status, &machine.CreatedAt, &machine.UpdatedAt, &machine.LastUpdatedBy); err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

func (r *MySqlMachineRepository) buildFilterQuery(baseQuery string, filter MachineFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Name != "" {
		filters = append(filters, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		filters = append(filters, "code LIKE ?")
		args = append(args, "%"+filter.Code+"%")
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
