// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"github.com/jackc/pgx/v5"
	"time"
)

const (
	deleteSelect = `
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s WHERE s.id = @id FOR UPDATE
    `
	deleteUpdate = `
        UPDATE $prefix$secret SET updated_at = @updatedAt, DATA = @data, is_deleted = true WHERE id = @id
    `
)

func (r *SecretDBRepository) Delete(ctx context.Context, id int64, oldUpdatedAt int64) error {
	tx, err := r.dp.Db.BeginTx(ctx, nil)
	defer tx.Rollback()

	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(deleteSelect), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		return fmt.Errorf("error fetching secret: %w", err)
	}
	if s.UpdatedAt != oldUpdatedAt {
		return repository.ErrSecretTimeNotMatch
	}

	s.UpdatedAt = time.Now().UnixNano()
	s.IsDeleted = true

	_, err = r.dp.Db.ExecContext(ctx, r.dp.PrefixQuery(deleteUpdate), pgx.NamedArgs{
		"id":        s.ID,
		"updatedAt": s.UpdatedAt,
		"data":      s.Data,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}
