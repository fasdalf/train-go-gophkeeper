// Package dbstorage - DB storage
package dbstorage

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s *DBProxy) Bootstrap(ctx context.Context) error {
	migrations := []string{
		// Пользователи
		`
		CREATE TABLE IF NOT EXISTS "$prefix$user" (
			id bigserial NOT NULL,
			login varchar(250) NOT NULL,
			pass_hash varchar(250) NOT NULL,
			CONSTRAINT "$prefix$user_pk" PRIMARY KEY (id),
			CONSTRAINT "$prefix$user_login_udx" UNIQUE (login)
		);`,
		// Пароли
		`
		CREATE TABLE IF NOT EXISTS "$prefix$secret" (
			id bigserial NOT NULL,
			user_id bigint NULL,
			updated_at bigint NOT NULL,
			is_deleted boolean NOT NULL DEFAULT FALSE,
			data bytea NOT NULL,
			CONSTRAINT "$prefix$secret_pk" PRIMARY KEY (id),
			CONSTRAINT "$prefix$secret_user_fk" FOREIGN KEY (user_id) REFERENCES "$prefix$user" (id)		    
		);`,
	}

	return s.runMigrations(ctx, migrations)
}

// Teardown Удаляет таблицы из БД
func (s *DBProxy) Teardown(ctx context.Context) error {
	migrations := []string{
		// Пароли
		`DROP TABLE IF EXISTS "$prefix$secret"`,
		// Пользователи
		`DROP TABLE IF EXISTS "$prefix$user"`,
	}

	return s.runMigrations(ctx, migrations)
}

func (s *DBProxy) runMigrations(ctx context.Context, migrations []string) error {
	// запускаем транзакцию
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		query := s.PrefixQuery(migration)

		if _, err = tx.ExecContext(ctx, query); err != nil {
			// откатываем транзакцию в случае ошибки
			_ = tx.Rollback()
			return fmt.Errorf("query execution error %w with query \"%s\"", err, query)
		}
	}

	// коммитим транзакцию
	return tx.Commit()
}
