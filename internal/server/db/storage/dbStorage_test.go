package storage

import (
	"context"
	"os"
	"testing"
)

func TestNewTestDbStorage(t *testing.T) {
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

	err = got.Teardown(context.Background())
	if err != nil {
		t.Errorf("Teardown() error = %v", err)
		return
	}

	os.Setenv(testDBDSNENVKey, "fail"+testDBDSNENVDefault)
	got, err = NewTestDbStorage()
	if err == nil {
		t.Errorf("NewTestDbStorage() did not set error")
		return
	}
	if got != nil {
		t.Errorf("NewTestDbStorage() did not return nil")
		return
	}
}
