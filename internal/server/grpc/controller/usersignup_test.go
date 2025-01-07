package controller

import (
	"context"
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/db/dbrepository"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"testing"
	"time"
)

func TestGophKeeperServer_UserSignUp(t *testing.T) {
	key := "the key"
	type fields struct {
		key           string
		userRepoUser  entity.User
		userRepoError error
	}
	type args struct {
		ctx context.Context
		r   *pb.UserSignUpRequest
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
					ID: 10,
				},
				userRepoError: nil,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignUpRequest{
					Login:    "login",
					Password: "123456",
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
					ID: 10,
				},
				userRepoError: nil,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignUpRequest{
					Login:    "login",
					Password: "",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "name taken",
			fields: fields{
				key:           key,
				userRepoUser:  entity.User{},
				userRepoError: repository.ErrUserExists,
			},
			args: args{
				ctx: context.Background(),
				r: &pb.UserSignUpRequest{
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
			got, err := s.UserSignUp(tt.args.ctx, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserSignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.want {
				t.Errorf("UserSignUp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
