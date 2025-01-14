package dbrepository

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"reflect"
	"testing"
)

func TestUserMockRepository_Create(t *testing.T) {
	want := &entity.User{
		ID:       10,
		Login:    "L",
		PassHash: "P",
	}
	r := &UserMockRepository{
		User:  want,
		Error: nil,
	}
	got, err := r.Create(context.Background(), "A", "B")
	if err != nil {
		t.Errorf("Create() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Create() got = %v, want %v", got, want)
	}
}

func TestUserMockRepository_FindById(t *testing.T) {
	want := &entity.User{
		ID:       10,
		Login:    "L",
		PassHash: "P",
	}
	r := &UserMockRepository{
		User:  want,
		Error: nil,
	}
	got, err := r.FindById(context.Background(), 10)
	if err != nil {
		t.Errorf("FindById() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindById() got = %v, want %v", got, want)
	}
}

func TestUserMockRepository_FindByLogin(t *testing.T) {
	want := &entity.User{
		ID:       10,
		Login:    "L",
		PassHash: "P",
	}
	r := &UserMockRepository{
		User:  want,
		Error: nil,
	}
	got, err := r.FindByLogin(context.Background(), "K")
	if err != nil {
		t.Errorf("FindByLogin() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FindByLogin() got = %v, want %v", got, want)
	}
}
