// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
)

const createInsert = `
        INSERT INTO $prefix$secret(user_id, updated_at, data) VALUES (@id, @updatedAt, @data) RETURNING id 
    `

func (r *SecretDBRepository) Create(ctx context.Context, userId int64, data []byte) (*entity.Secret, error) {
	s := &entity.Secret{
		ID:        0,
		UserId:    userId,
		UpdatedAt: time.Now().UnixNano(),
		IsDeleted: false,
		Data:      data,
	}
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(createInsert), pgx.NamedArgs{
		"id":        s.UserId,
		"updatedAt": s.UpdatedAt,
		"data":      s.Data,
	})
	if err := row.Scan(&s.ID); err != nil {
		return nil, fmt.Errorf("error creating secret: %w", err)
	}
	return s, nil
}
