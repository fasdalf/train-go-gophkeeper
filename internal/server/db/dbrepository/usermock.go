// Package dbrepository - DB repository, users file
package dbrepository

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

// typeCheckUserMockRepository ensures that UserDBRepository implements the UserRepository interface.
var typeCheckUserMockRepository repository.UserRepository = &UserMockRepository{}

type UserMockRepository struct {
	User  *entity.User
	Error error
}

// FindById retrieves a user by their id.
func (r *UserMockRepository) FindById(ctx context.Context, id uint64) (*entity.User, error) {
	return r.User, r.Error
}

// FindByLogin retrieves a user by their login.
func (r *UserMockRepository) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	return r.User, r.Error
}

// Create inserts a new user into the database.
func (r *UserMockRepository) Create(ctx context.Context, login, passHash string) (*entity.User, error) {
	return r.User, r.Error
}
