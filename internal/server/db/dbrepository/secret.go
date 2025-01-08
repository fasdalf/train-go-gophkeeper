// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r *SecretDBRepository) FindById(ctx context.Context, id int64) (*entity.Secret, error) {
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s
        WHERE s.id = @id
    `), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(repository.ErrSecretNotFound, err)
		}

		return nil, fmt.Errorf("error finding secret by id '%d': %w", id, err)
	}

	return s, nil
}
func (r *SecretDBRepository) FindAfter(ctx context.Context, userId int64, updatedAt int64) ([]*entity.Secret, error) {
	rows, err := r.dp.Db.QueryContext(ctx, r.dp.PrefixQuery(`
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s
        WHERE s.user_id = @id AND s.updated_at > @updatedAt
    `), pgx.NamedArgs{"id": userId, "updatedAt": updatedAt})
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
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *SecretDBRepository) Create(ctx context.Context, userId int64, data []byte) (*entity.Secret, error) {
	s := &entity.Secret{
		ID:        0,
		UserId:    userId,
		UpdatedAt: time.Now().UnixNano(),
		IsDeleted: false,
		Data:      data,
	}
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
        INSERT INTO $prefix$secret(user_id, updated_at, data) VALUES (@id, @updatedAt, @data) RETURNING id 
    `), pgx.NamedArgs{
		"id":        s.UserId,
		"updatedAt": s.UpdatedAt,
		"data":      s.Data,
	})
	if err := row.Scan(&s.ID); err != nil {
		return nil, fmt.Errorf("error creating secret: %w", err)
	}
	return s, nil
}

func (r *SecretDBRepository) Update(ctx context.Context, id int64, data []byte, oldUpdatedAt int64) (*entity.Secret, error) {
	tx, err := r.dp.Db.BeginTx(ctx, nil)
	defer tx.Rollback()

	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s WHERE s.id = @id FOR UPDATE
    `), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		return nil, fmt.Errorf("error fetching secret: %w", err)
	}
	if s.UpdatedAt != oldUpdatedAt {
		return nil, repository.ErrSecretTimeNotMatch
	}

	s.UpdatedAt = time.Now().UnixNano()
	s.Data = data

	_, err = r.dp.Db.ExecContext(ctx, r.dp.PrefixQuery(`
        UPDATE $prefix$secret SET updated_at = @updatedAt, DATA = @data WHERE id = @id
    `), pgx.NamedArgs{
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
func (r *SecretDBRepository) Delete(ctx context.Context, id int64, oldUpdatedAt int64) error {
	tx, err := r.dp.Db.BeginTx(ctx, nil)
	defer tx.Rollback()

	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
        SELECT s.id, s.user_id, s.updated_at, s.is_deleted, s.data FROM $prefix$secret s WHERE s.id = @id FOR UPDATE
    `), pgx.NamedArgs{"id": id})
	s, err := r.scanSecret(row)
	if err != nil {
		return fmt.Errorf("error fetching secret: %w", err)
	}
	if s.UpdatedAt != oldUpdatedAt {
		return repository.ErrSecretTimeNotMatch
	}

	s.UpdatedAt = time.Now().UnixNano()
	s.IsDeleted = true

	_, err = r.dp.Db.ExecContext(ctx, r.dp.PrefixQuery(`
        UPDATE $prefix$secret SET updated_at = @updatedAt, DATA = @data, is_deleted = true WHERE id = @id
    `), pgx.NamedArgs{
		"id":        s.ID,
		"updatedAt": s.UpdatedAt,
		"data":      s.Data,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SecretDBRepository) scanSecret(row pgx.Row) (*entity.Secret, error) {
	var s entity.Secret
	if err := row.Scan(&s.ID, &s.UserId, &s.UpdatedAt, &s.IsDeleted, &s.Data); err != nil {
		return nil, err
	}
	return &s, nil
}
