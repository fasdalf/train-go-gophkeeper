package repository

import (
	"context"
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
)

type UserRepository interface {
	FindById(ctx context.Context, Id int64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	Create(ctx context.Context, login, passHash string) (*entity.User, error)
}

type SecretRepository interface {
	FindById(ctx context.Context, id int64) (*entity.Secret, error)
	FindAfter(ctx context.Context, userId int64, updatedAt int64) ([]*entity.Secret, error)
	Create(ctx context.Context, userId int64, data []byte) (*entity.Secret, error)
	Update(ctx context.Context, id int64, data []byte, oldUpdatedAt int64) (*entity.Secret, error)
	Delete(ctx context.Context, id int64, oldUpdatedAt int64) error
}

var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrUserExists         = fmt.Errorf("user already exists")
	ErrSecretNotFound     = fmt.Errorf("secret not found")
	ErrSecretTimeNotMatch = fmt.Errorf("secret time does not match")
)
