// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"github.com/jackc/pgx/v5"
	"time"
)

const (
	updateSelect = `
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s WHERE s.id = @id FOR UPDATE
    `
	updateUpdate = `
        UPDATE $prefix$secret SET updated_at = @updatedAt, DATA = @data WHERE id = @id
    `
)

func (r *SecretDBRepository) Update(ctx context.Context, id int64, data []byte, oldUpdatedAt int64) (*entity.Secret, error) {
	tx, err := r.dp.Db.BeginTx(ctx, nil)
	defer tx.Rollback()

	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(updateSelect), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		return nil, fmt.Errorf("error fetching secret: %w", err)
	}
	if s.UpdatedAt != oldUpdatedAt {
		return nil, repository.ErrSecretTimeNotMatch
	}

	s.UpdatedAt = time.Now().UnixNano()
	s.Data = data

	_, err = r.dp.Db.ExecContext(ctx, r.dp.PrefixQuery(updateUpdate), pgx.NamedArgs{
		"id":        s.ID,
		"updatedAt": s.UpdatedAt,
		"data":      s.Data,
	})
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		s = nil
	}

	return s, err
}
