package config

import (
	flag "github.com/spf13/pflag"
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	InitConfig()
	cfg := GetConfig()

	if cfg.GRPCAddr == "" {
		t.Errorf("Expected ADDRESS to be %s, but got %s", defaultGRPCAddr, cfg.GRPCAddr)
	}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	oldEnv := os.Getenv("TOKEN_EXP")
	defer os.Setenv("TOKEN_EXP", oldEnv)
	os.Setenv("TOKEN_EXP", "invalid")
	InitConfig()
}
