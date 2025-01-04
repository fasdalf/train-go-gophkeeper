package config

import "testing"

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()

	if cfg.GRPCAddr == "" {
		t.Errorf("Expected ADDRESS to be %s, but got %s", defaultGRPCAddr, cfg.GRPCAddr)
	}
}
