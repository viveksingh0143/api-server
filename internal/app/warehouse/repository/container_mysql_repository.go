package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlContainerRepository struct {
	conn database.Connection
}

func NewContainerRepository(conn database.Connection) ContainerRepository {
	return &MySqlContainerRepository{conn: conn}
}

func (r *MySqlContainerRepository) Create(container *domain.Container) error {
	query := "INSERT INTO containers (type, code, name, address, status, last_updated_by) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, container.Type, container.Code, container.Name, container.Address, container.Status, container.LastUpdatedBy)
	return err
}

func (r *MySqlContainerRepository) Update(container *domain.Container) error {
	query := "UPDATE containers SET type=?, code=?, name=?, address=?, status=?, last_updated_by=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, container.Type, container.Code, container.Name, container.Address, container.Status, container.LastUpdatedBy, container.ID)
	return err
}

func (r *MySqlContainerRepository) Delete(containerID int64) error {
	query := "DELETE FROM containers WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, containerID)
	return err
}

func (r *MySqlContainerRepository) GetById(containerID int64) (*domain.Container, error) {
	query := "SELECT id, type, code, name, address, status, created_at, updated_at, last_updated_by FROM containers WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, containerID)
	container := &domain.Container{}
	err := row.Scan(&container.ID, &container.Type, &container.Code, &container.Name, &container.Address, &container.Status, &container.CreatedAt, &container.UpdatedAt, &container.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return container, nil
}

func (r *MySqlContainerRepository) GetTotalCount(filter ContainerFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM containers", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlContainerRepository) GetAll(page int, pageSize int, sort string, filter ContainerFilterOptions) ([]*domain.Container, error) {
	query, args := r.buildFilterQuery("SELECT id, type, code, name, address, status, created_at, updated_at, last_updated_by FROM containers", filter)
	var allowedSortOrders = map[string]bool{
		"type ASC":    true,
		"type DESC":   true,
		"name ASC":    true,
		"name DESC":   true,
		"code ASC":    true,
		"code DESC":   true,
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

	var containers []*domain.Container
	for rows.Next() {
		container := &domain.Container{}
		if err := rows.Scan(&container.ID, &container.Type, &container.Code, &container.Name, &container.Address, &container.Status, &container.CreatedAt, &container.UpdatedAt, &container.LastUpdatedBy); err != nil {
			return nil, err
		}
		containers = append(containers, container)
	}
	return containers, nil
}

func (r *MySqlContainerRepository) buildFilterQuery(baseQuery string, filter ContainerFilterOptions) (string, []interface{}) {
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
