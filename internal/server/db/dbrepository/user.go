// Package dbrepository - DB repository, users file
package dbrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// typeCheckUserDBRepository ensures that UserDBRepository implements the UserRepository interface.
var typeCheckUserDBRepository repository.UserRepository = &UserDBRepository{}

type UserDBRepository struct {
	dp *dbstorage.DBProxy
}

// NewUserDBRepository init
func NewUserDBRepository(dp *dbstorage.DBProxy) *UserDBRepository {
	return &UserDBRepository{
		dp: dp,
	}
}

// FindById retrieves a user by their id.
func (r *UserDBRepository) FindById(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
        SELECT u.id, u.login, u.pass_hash FROM $prefix$user u
        WHERE u.id = @id
    `), pgx.NamedArgs{"id": id})
	err := row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, repository.ErrUserNotFound)
		}

		return nil, fmt.Errorf("error finding user by ID '%d': %w", id, err)
	}

	return &user, nil
}

// FindByLogin retrieves a user by their login.
func (r *UserDBRepository) FindByLogin(ctx context.Context, login string) (*entity.User, error) {
	var user entity.User
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
		SELECT u.id, u.login, u.pass_hash FROM $prefix$user u
		WHERE u.login = @login
	`), pgx.NamedArgs{"login": login})
	err := row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(repository.ErrUserNotFound, err)
		}

		return nil, fmt.Errorf("error finding user by Login '%s': %w", login, err)
	}

	return &user, nil
}

// Create inserts a new user into the database.
func (r *UserDBRepository) Create(ctx context.Context, login, passHash string) (*entity.User, error) {
	user := entity.User{Login: login, PassHash: passHash}
	row := r.dp.Db.QueryRowContext(ctx, r.dp.PrefixQuery(`
		INSERT INTO $prefix$user (login, pass_hash) VALUES (@login, @pass_hash)
		RETURNING id
	`), pgx.NamedArgs{"login": login, "pass_hash": passHash})
	err := row.Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgErr != nil && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			err = errors.Join(repository.ErrUserExists, err)
		}

		return nil, fmt.Errorf("error creating user with Login '%s': %w", login, err)
	}

	return &user, nil
}
