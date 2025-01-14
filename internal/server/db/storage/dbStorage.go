// Package storage - DB storage
package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"math/rand"
	"os"
	"strings"
)

const (
	TablePrefix         = "$prefix$"
	testDBDSNENVKey     = "DATABASE_URI_TEST"
	testDBDSNENVDefault = "host=train-go-gophkeeper_db user=postgres password=postgresP@SS dbname=postgres_gophkeeper sslmode=disable"
)

// DBProxy store metrics in DB
type DBProxy struct {
	Db     *sql.DB
	Prefix string
}

func NewDBStorage(ctx context.Context, dsn string, prefix string) (s *DBProxy) {
	// err is possible only when no driver imported
	db, _ := sql.Open("pgx", dsn)
	s = &DBProxy{
		Db:     db,
		Prefix: prefix,
	}
	return s
}

// NewTestDbStorage creates new DBProxy with random prefix for parallel tests run.
func NewTestDbStorage() (s *DBProxy, err error) {
	dsn := testDBDSNENVDefault
	envDSN := os.Getenv(testDBDSNENVKey)
	if envDSN != "" {
		dsn = envDSN
	}

	ctx := context.Background()
	s = NewDBStorage(ctx, dsn, fmt.Sprintf("go_test_%06d_", rand.Intn(999998)+1))

	if err = s.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("failed to bootstrap database: %w", err)
	}

	return s, nil
}

// PrefixQuery adds prefix to table name where TablePrefix used
func (s *DBProxy) PrefixQuery(query string) string {
	query = strings.ReplaceAll(query, TablePrefix, s.Prefix)
	return query
}
