package config

import (
	goflag "flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/configfile"
	flag "github.com/spf13/pflag"
)

const (
	defaultGRPCAddr = "localhost:8090"
)

type Config struct {
	GRPCAddr string `env:"CLIENT_GRPC_ADDRESS" json:"client_grpc_address"`
}

var config *Config = &Config{
	GRPCAddr: defaultGRPCAddr,
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
	flag.StringVarP(&config.GRPCAddr, "grpcaddress", "g", config.GRPCAddr, "The address to listen on for GRPC requests.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
	// pflag handles --help itself.
}
