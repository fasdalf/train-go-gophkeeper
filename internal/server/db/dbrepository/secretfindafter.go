// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
)

const findAfterSelect = `
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s
        WHERE s.user_id = @id AND s.updated_at > @updatedAt
    `

func (r *SecretDBRepository) FindAfter(ctx context.Context, userId int64, updatedAt int64) ([]*entity.Secret, error) {
	rows, err := r.dp.Db.QueryContext(ctx, r.dp.PrefixQuery(findAfterSelect), pgx.NamedArgs{"id": userId, "updatedAt": updatedAt})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*entity.Secret{}
	for rows.Next() {
		s, err := r.scanSecret(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secrets: %w", err)
		}
		results = append(results, s)
	}
	err = rows.Err()

	return results, err
}
