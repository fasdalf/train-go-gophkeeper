package interceptors

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/jwt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
)

func TestNewValidateHashInterceptor(t *testing.T) {
	key := "the key"
	validToken := jwt.BuildJWTString(100, &key, 3*time.Hour)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		wantResp any
		wantErr  bool
	}{
		{
			name: "valid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{jwt.AuthHeader: validToken}),
				)},
			wantResp: 10,
			wantErr:  false,
		},
		{
			name: "invalid",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.New(map[string]string{jwt.AuthHeader: "mock"}),
				),
			},
			wantResp: nil,
			wantErr:  true,
		},
	}

	usi := grpc.UnaryServerInfo{}
	nh := func(ctx context.Context, req any) (any, error) {
		return 10, nil
	}
	req := pb.CreateSecretRequest{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewValidateTokenInterceptor(&key)
			gotResp, err := got(tt.args.ctx, &req, &usi, nh)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTokenInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ValidateTokenInterceptor() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
