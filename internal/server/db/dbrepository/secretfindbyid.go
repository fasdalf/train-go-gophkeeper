// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

const findByIdSelect = `
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s
        WHERE s.id = @id
    `

func (r *SecretDBRepository) FindById(ctx context.Context, id int64) (*entity.Secret, error) {
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(findByIdSelect), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(repository.ErrSecretNotFound, err)
		}

		return nil, fmt.Errorf("error finding secret by id '%d': %w", id, err)
	}

	return s, nil
}
