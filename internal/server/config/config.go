package config

import (
	goflag "flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"

	"github.com/fasdalf/train-go-gophkeeper/internal/common/configfile"
)

const (
	defaultCryptoKey = "please-set-in-env"
	defaultGRPCAddr  = ":8090"
	defaultDBPrefix  = "gk_"
	defaultTokenExp  = 3 * time.Hour
)

type Config struct {
	GRPCAddr        string        `env:"GRPC_ADDRESS" json:"grpc_address"`
	StorageDBDSN    string        `env:"DATABASE_URI" json:"database_uri"`
	StorageDBPrefix string        `env:"DATABASE_TABLE_PREFIX" json:"database_table_prefix"`
	CryptoKey       string        `env:"KEY" json:"key"`
	TokenExp        time.Duration `env:"TOKEN_EXP" json:"token_exp"`
}

var config *Config = &Config{
	GRPCAddr:        defaultGRPCAddr,
	CryptoKey:       defaultCryptoKey,
	StorageDBPrefix: defaultDBPrefix,
	TokenExp:        defaultTokenExp,
}

func GetConfig() Config {
	return *config
}

func init() {
	// Env. variables.
	if err := env.Parse(config); err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Config file overwrites environment
	configfile.ParseFile(config)

	// Flag overwrites config file
	flag.StringVarP(&config.GRPCAddr, "grpcaddress", "g", config.GRPCAddr, "The address to listen on for GRPC requests. Required.")
	flag.StringVarP(&config.StorageDBDSN, "databasedsn", "d", config.StorageDBDSN, "Postgres PGX DSN to use DB storage. Required.")
	flag.StringVarP(&config.StorageDBPrefix, "databaseprefix", "p", config.StorageDBPrefix, "Postgres tables prefix.")
	flag.StringVarP(&config.CryptoKey, "key", "k", config.CryptoKey, "Key for authorizations cryptography.")
	flag.DurationVarP(&config.TokenExp, "tokenexp", "e", config.TokenExp, "JWT token lifetime")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	// pflag handles --help itself.

	// configure logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
}
