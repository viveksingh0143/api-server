package usecase

import (
	"github.com/vamika-digital/wms-api-server/internal/app/user/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCaseImpl struct {
	Repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &UserUseCaseImpl{Repo: repo}
}

func (u *UserUseCaseImpl) CreateUser(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return u.Repo.Create(user)
}

func (u *UserUseCaseImpl) UpdateUser(user *domain.User) error {
	// Check for an existing user with the specified ID
	existingUser, err := u.Repo.GetById(user.ID)
	if err != nil {
		return err
	}

	// Hash the new password if it has been changed
	if user.PasswordHash != "" && user.PasswordHash != existingUser.PasswordHash {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}

	return u.Repo.Update(user)
}

func (u *UserUseCaseImpl) DeleteUser(userID int64) error {
	return u.Repo.Delete(userID)
}

func (u *UserUseCaseImpl) GetUserByID(userID int64) (*domain.User, error) {
	return u.Repo.GetById(userID)
}

func (u *UserUseCaseImpl) GetAllUsers(page int, pageSize int, sort string, filter repository.UserFilterOptions) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := u.Repo.GetAll((page-1)*pageSize, pageSize, sort, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the total count of users matching the filter
	total, err := u.Repo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UserUseCaseImpl) GetUserByUsername(username string) (*domain.User, error) {
	return u.Repo.FindByUsername(username)
}

func (u *UserUseCaseImpl) GetUserByEmail(email string) (*domain.User, error) {
	return u.Repo.FindByEmail(email)
}
