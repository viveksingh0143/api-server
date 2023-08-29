package repository

import "github.com/vamika-digital/wms-api-server/internal/app/user/domain"

type UserRepository interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(userID int64) error
	GetById(userID int64) (*domain.User, error)
	GetTotalCount(filter UserFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter UserFilterOptions) ([]*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
}

type UserFilterOptions struct {
	Name     string
	Username string // Filter by username, using a LIKE query
	Email    string // Filter by email, using a LIKE query
	Status   string // Filter by user status, using an exact match
}
