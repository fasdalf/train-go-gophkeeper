package controller

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/cryptofacade"
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/db/dbrepository"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"testing"
	"time"
)

func TestGophKeeperServer_UserSignIn(t *testing.T) {
	key := "the key"
	password := "123456"
	passHash := cryptofacade.Hash(&password, &key)
	type fields struct {
		key           string
		userRepoUser  entity.User
		userRepoError error
	}
	type args struct {
		ctx context.Context
		r   *pb.UserSignInRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success 1",
			fields: fields{
				key: key,
				userRepoUser: entity.User{
					ID:       10,
					Login:    "login",
					PassHash: passHash,
				},
				userRepoError: nil,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignInRequest{
					Login:    "login",
					Password: password,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "empty password",
			fields: fields{
				key: key,
				userRepoUser: entity.User{
					ID:       10,
					Login:    "login",
					PassHash: passHash,
				},
				userRepoError: nil,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignInRequest{
					Login:    "login",
					Password: "",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "wrong password",
			fields: fields{
				key: key,
				userRepoUser: entity.User{
					ID:       10,
					Login:    "login",
					PassHash: passHash,
				},
				userRepoError: nil,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignInRequest{
					Login:    "login",
					Password: "HACK",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := dbrepository.UserMockRepository{
				&tt.fields.userRepoUser,
				tt.fields.userRepoError,
			}
			s := NewGophKeeperServer(
				&ur,
				&tt.fields.key,
				3*time.Second,
				nil,
			)
			got, err := s.UserSignIn(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserSignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("UserSignIn() got = %v, want %v", got, tt.want)
			}
		})
	}
}
