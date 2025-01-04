package repository

import (
	"context"
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
)

type UserRepository interface {
	FindById(ctx context.Context, Id uint64) (*entity.User, error)
	FindByLogin(ctx context.Context, login string) (*entity.User, error)
	Create(ctx context.Context, login, passHash string) (*entity.User, error)
}

type SecretRepository interface {
	FindById(Id uint64) (*entity.Secret, error)
	FindAfter(UserId uint64, UpdatedAt uint64) ([]*entity.Secret, error)
	Create(UserId uint64, Data []byte) (*entity.Secret, error)
	Update(Id uint64, Data []byte, OldUpdatedAt uint64) (*entity.Secret, error)
	Delete(Id uint64, OldUpdatedAt uint64) (*entity.Secret, error)
}

var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrUserExists         = fmt.Errorf("user already exists")
	ErrSecretNotFound     = fmt.Errorf("secret not found")
	ErrSecretTimeNotMatch = fmt.Errorf("secret time does not match")
)
