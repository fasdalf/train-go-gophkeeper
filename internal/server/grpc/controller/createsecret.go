package controller

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/contextkeys"
)

func (s *GophKeeperServer) CreateSecret(ctx context.Context, r *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	// Protected with NewValidateTokenInterceptor()
	userID := ctx.Value(contextkeys.UserIDKey).(int64)

	secret, err := s.secretRepo.Create(ctx, userID, r.GetData())
	if err != nil {
		slog.Error("failed to create secret", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to create secret")
		return nil, gErr
	}

	resp := &pb.CreateSecretResponse{
		Id: secret.ID,
		Ts: secret.UpdatedAt,
	}
	return resp, nil
}
