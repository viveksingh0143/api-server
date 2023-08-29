package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/user/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlUserRepository struct {
	conn database.Connection
}

func NewUserRepository(conn database.Connection) UserRepository {
	return &MySqlUserRepository{conn: conn}
}

func (r *MySqlUserRepository) Create(user *domain.User) error {
	query := "INSERT INTO users (username, password_hash, name, staff_id, email, email_confirmation, status, last_updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, user.Username, user.PasswordHash, user.Name, user.StaffID, user.Email, user.EmailConfirmation, user.Status, user.LastUpdatedBy)
	return err
}

func (r *MySqlUserRepository) Update(user *domain.User) error {
	query := "UPDATE users SET username=?, password_hash=?, name=?, staff_id=?, email=?, email_confirmation=?, status=?, last_updated_by=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, user.Username, user.PasswordHash, user.Name, user.StaffID, user.Email, user.EmailConfirmation, user.Status, user.LastUpdatedBy, user.ID)
	return err
}

func (r *MySqlUserRepository) Delete(userID int64) error {
	query := "DELETE FROM users WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, userID)
	return err
}

func (r *MySqlUserRepository) GetById(userID int64) (*domain.User, error) {
	query := "SELECT id, username, password_hash, name, staff_id, email, email_confirmation, status, created_at, updated_at, last_updated_by FROM users WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, userID)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.StaffID, &user.Email, &user.EmailConfirmation, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *MySqlUserRepository) GetTotalCount(filter UserFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM users", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlUserRepository) GetAll(page int, pageSize int, sort string, filter UserFilterOptions) ([]*domain.User, error) {
	query, args := r.buildFilterQuery("SELECT id, username, password_hash, name, staff_id, email, email_confirmation, status, created_at, updated_at, last_updated_by FROM users", filter)

	var allowedSortOrders = map[string]bool{
		"username ASC":  true,
		"username DESC": true,
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

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.StaffID, &user.Email, &user.EmailConfirmation, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastUpdatedBy); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *MySqlUserRepository) FindByUsername(username string) (*domain.User, error) {
	query := "SELECT id, username, password_hash, name, staff_id, email, email_confirmation, status, created_at, updated_at, last_updated_by FROM users WHERE username = ?"
	row := r.conn.GetDB().QueryRow(query, username)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.StaffID, &user.Email, &user.EmailConfirmation, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *MySqlUserRepository) FindByEmail(email string) (*domain.User, error) {
	query := "SELECT id, username, password_hash, name, staff_id, email, email_confirmation, status, created_at, updated_at, last_updated_by FROM users WHERE email = ?"
	row := r.conn.GetDB().QueryRow(query, email)
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Name, &user.StaffID, &user.Email, &user.EmailConfirmation, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *MySqlUserRepository) buildFilterQuery(baseQuery string, filter UserFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Name != "" {
		filters = append(filters, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.Username != "" {
		filters = append(filters, "username LIKE ?")
		args = append(args, "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		filters = append(filters, "email LIKE ?")
		args = append(args, "%"+filter.Email+"%")
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
