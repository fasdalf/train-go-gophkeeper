// Package dbstorage - DB storage
package storage

import (
	"context"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	goose "github.com/pressly/goose/v3"
	"os"
	"sync"
)

const (
	gooseMinimalVersion = 20250113204423
	goosePrefixEnv      = "DATABASE_MIGRATION_TABLE_PREFIX"
)

var (
	//go:embed migrations/*.sql
	gooseEmbedMigrations   embed.FS
	gooseOriginalTableName string
	gooseLock              sync.Mutex
	ErrDBSchemaIsOld       = fmt.Errorf("database schema is older than %d, please migrate", gooseMinimalVersion)
)

func init() {
	gooseOriginalTableName = goose.TableName()
	goose.SetBaseFS(gooseEmbedMigrations)
	_ = goose.SetDialect("postgres")
}

func (s *DBProxy) CheckVersion(ctx context.Context) (err error) {
	gooseLock.Lock()
	defer gooseLock.Unlock()
	setGoosePrefix(s.Prefix)

	current, err := goose.GetDBVersionContext(ctx, s.Db)
	if err == nil && current < gooseMinimalVersion {
		err = ErrDBSchemaIsOld
	}

	return
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s *DBProxy) Bootstrap(ctx context.Context) error {
	gooseLock.Lock()
	defer gooseLock.Unlock()
	setGoosePrefix(s.Prefix)

	return goose.UpContext(ctx, s.Db, "migrations", goose.WithAllowMissing())
}

// Teardown Удаляет таблицы из БД
func (s *DBProxy) Teardown(ctx context.Context) (err error) {
	const sqlTemplate = `DROP TABLE IF EXISTS $prefix$%s;`
	sql := s.PrefixQuery(fmt.Sprintf(sqlTemplate, gooseOriginalTableName))

	err = s.Reset(ctx)

	if err == nil {
		_, err = s.Db.ExecContext(ctx, sql)
	}

	return
}

func (s *DBProxy) Reset(ctx context.Context) (err error) {
	gooseLock.Lock()
	defer gooseLock.Unlock()
	_ = os.Setenv(goosePrefixEnv, s.Prefix)
	goose.SetBaseFS(gooseEmbedMigrations)

	err = goose.ResetContext(ctx, s.Db, "migrations")

	return
}

func setGoosePrefix(prefix string) {
	_ = os.Setenv(goosePrefixEnv, prefix)
	goose.SetTableName(prefix + gooseOriginalTableName)
}
