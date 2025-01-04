package dbrepository

import (
	"context"
	"errors"
	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"reflect"
	"strings"
	"testing"
)

func TestUserDBRepository_Create_AllCases(t *testing.T) {
	c := context.Background()
	dp, err := dbstorage.NewTestDbStorage()
	if err != nil {
		t.Errorf("NewTestDbStorage() error = %v", err)
		return
	}
	defer dp.Teardown(context.Background())
	rep := NewUserDBRepository(dp)

	// Check 1st creation, success
	user, err := rep.Create(c, "user1", "pass1")
	if err != nil {
		t.Errorf("1st Create() error = %v", err)
		return
	}
	if user == nil {
		t.Errorf("1st Create() returned nil")
		return
	}
	if user.ID == 0 {
		t.Errorf("1st Create() user id not set")
		return
	}

	// 2nd creation, user already exists
	user, err = rep.Create(c, "user1", "pass1")
	if err == nil {
		t.Errorf("2nd Create() did not set error")
		return
	}
	if !errors.Is(err, repository.ErrUserExists) {
		t.Errorf("2nd Create() error = %v but repository.ErrUserExists expectred", err)
		return
	}

	// 3rd creation, login is too long
	s := strings.Repeat("user1", 300/5)
	user, err = rep.Create(c, s, "pass1")
	if err == nil {
		t.Errorf("3rd Create() did not set error")
		return
	}
	if errors.Is(err, repository.ErrUserExists) {
		t.Errorf("3rd Create() error is %v", err)
		return
	}
}

func TestUserDBRepository_FindById(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{1},
			want: &entity.User{
				ID:       1,
				Login:    "user1",
				PassHash: "pass1",
			},
			wantErr: false,
		},
		{
			name:    "fail",
			args:    args{2},
			want:    nil,
			wantErr: true,
		},
	}

	// set up one user
	c := context.Background()
	dp, err := dbstorage.NewTestDbStorage()
	if err != nil {
		t.Errorf("NewTestDbStorage() error = %v", err)
		return
	}
	defer dp.Teardown(context.Background())
	r := NewUserDBRepository(dp)

	user, err := r.Create(c, "user1", "pass1")
	if err != nil {
		t.Errorf("1st Create() error = %v", err)
		return
	}
	if user == nil {
		t.Errorf("1st Create() returned nil")
		return
	}
	if user.ID != 1 {
		t.Errorf("1st Create() user id not set")
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindById(c, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDBRepository_FindByLogin(t *testing.T) {
	type args struct {
		login string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "success",
			args: args{"user1"},
			want: &entity.User{
				ID:       1,
				Login:    "user1",
				PassHash: "pass1",
			},
			wantErr: false,
		},
		{
			name:    "fail",
			args:    args{"user2"},
			want:    nil,
			wantErr: true,
		},
	}

	// set up one user
	c := context.Background()
	dp, err := dbstorage.NewTestDbStorage()
	if err != nil {
		t.Errorf("NewTestDbStorage() error = %v", err)
		return
	}
	defer dp.Teardown(context.Background())
	r := NewUserDBRepository(dp)

	user, err := r.Create(c, "user1", "pass1")
	if err != nil {
		t.Errorf("1st Create() error = %v", err)
		return
	}
	if user == nil {
		t.Errorf("1st Create() returned nil")
		return
	}
	if user.ID != 1 {
		t.Errorf("1st Create() user id not set")
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindByLogin(c, tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindById() got = %v, want %v", got, tt.want)
			}
		})
	}
}
