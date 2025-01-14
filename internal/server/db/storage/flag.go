package storage

import (
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

const (
	migrateFlag     = "m"
	migrateFlagFull = "migrate"
	migrateFlagHelp = "Migrate database to current version and exit"
)

func init() {
	fmt.Println("internal/server/db/storage/command.go")
	var f bool
	flag.BoolVarP(&f, migrateFlagFull, migrateFlag, false, migrateFlagHelp)
}

func (s *DBProxy) MigrateByFlag() (done bool, err error) {
	var f bool
	cl := flag.NewFlagSet(migrateFlagFull, flag.ContinueOnError)
	cl.BoolVarP(&f, migrateFlagFull, migrateFlag, done, "")
	cl.Parse(os.Args[1:])

	if !f {
		return
	}

	err = s.Bootstrap(context.Background())

	if err == nil {
		done = true
	}

	return
}
