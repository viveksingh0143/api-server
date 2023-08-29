package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/user/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/user/repository"
)

type UserUseCase interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(userID int64) error
	GetUserByID(userID int64) (*domain.User, error)
	GetAllUsers(page int, pageSize int, sort string, filter repository.UserFilterOptions) ([]*domain.User, int, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}
