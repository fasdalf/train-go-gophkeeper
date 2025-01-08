// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

var typeCheckSecretMockRepository repository.SecretRepository = &SecretMockRepository{}

type SecretMockRepository struct {
	Secret  *entity.Secret
	Secrets []*entity.Secret
	Error   error
}

func (r *SecretMockRepository) FindById(ctx context.Context, Id int64) (*entity.Secret, error) {
	return r.Secret, r.Error
}
func (r *SecretMockRepository) FindAfter(ctx context.Context, UserId int64, UpdatedAt int64) ([]*entity.Secret, error) {
	return r.Secrets, r.Error
}
func (r *SecretMockRepository) Create(ctx context.Context, UserId int64, Data []byte) (*entity.Secret, error) {
	return r.Secret, r.Error
}
func (r *SecretMockRepository) Update(ctx context.Context, Id int64, Data []byte, OldUpdatedAt int64) (*entity.Secret, error) {
	return r.Secret, r.Error
}
func (r *SecretMockRepository) Delete(ctx context.Context, Id int64, OldUpdatedAt int64) error {
	return r.Error
}
