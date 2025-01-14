// Package dbrepository - DB repository, secrets file
package dbrepository

import (
	"context"
	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	"testing"
)

func TestSecretDBRepository_Update(t *testing.T) {
	var oldAt, newAt int64
	type args struct {
		id        int64
		updatedAt *int64
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "fail ts",
			args:    args{1, &oldAt},
			want:    false,
			wantErr: true,
		},
		{
			name:    "success",
			args:    args{1, &newAt},
			want:    true,
			wantErr: false,
		},
		{
			name:    "fail id",
			args:    args{2, &oldAt},
			want:    false,
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

	ur := NewUserDBRepository(dp)
	user, err := ur.Create(c, "user1", "pass1")
	if err != nil {
		t.Errorf("user Create() error = %v", err)
		return
	}
	r := NewSecretDBRepository(dp)
	secret, err := r.Create(c, user.ID, []byte("pass1"))
	if err != nil {
		t.Errorf("secret Create() error = %v", err)
		return
	}
	if secret == nil {
		t.Errorf("secret Create() returned nil")
		return
	}
	newAt = secret.UpdatedAt

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Update(c, tt.args.id, []byte("payload"), *tt.args.updatedAt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}
