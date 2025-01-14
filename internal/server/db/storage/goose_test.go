package storage

import (
	"context"
	"os"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	oldEnv := os.Getenv(testDBDSNENVKey)
	os.Setenv(testDBDSNENVKey, testDBDSNENVDefault)
	defer os.Setenv(testDBDSNENVKey, oldEnv)

	got, err := NewTestDbStorage()
	defer got.Teardown(context.Background())
	if err != nil {
		t.Errorf("NewTestDbStorage() error = %v", err)
		return
	}
	if got == nil {
		t.Errorf("NewTestDbStorage() returned nil")
		return
	}

	err = got.CheckVersion(context.Background())
	if err != nil {
		t.Errorf("CheckVersion() error = %v", err)
		return
	}

	err = got.Reset(context.Background())
	if err != nil {
		t.Errorf("Reset() error = %v", err)
		return
	}

	err = got.CheckVersion(context.Background())
	if err == nil {
		t.Errorf("CheckVersion() with no error on empty db")
		return
	}
}
