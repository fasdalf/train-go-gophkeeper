// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"github.com/jackc/pgx/v5"

	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

var typeCheckSecretDBRepository repository.SecretRepository = &SecretDBRepository{}

type SecretDBRepository struct {
	dp *dbstorage.DBProxy
}

// NewSecretDBRepository init
func NewSecretDBRepository(dp *dbstorage.DBProxy) *SecretDBRepository {
	return &SecretDBRepository{
		dp: dp,
	}
}

func (r *SecretDBRepository) scanSecret(row pgx.Row) (*entity.Secret, error) {
	var s entity.Secret
	if err := row.Scan(&s.ID, &s.UserId, &s.UpdatedAt, &s.IsDeleted, &s.Data); err != nil {
		return nil, err
	}
	return &s, nil
}
