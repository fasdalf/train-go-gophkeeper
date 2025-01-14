package controller

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/contextkeys"
)

func (s *GophKeeperServer) ListSecrets(ctx context.Context, r *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	// Protected with NewValidateTokenInterceptor()
	userID := ctx.Value(contextkeys.UserIDKey).(int64)
	syncTs := time.Now().UnixNano()
	fromTs := int64(0)
	if r != nil {
		fromTs = r.FromTs
	}

	secrets, err := s.secretRepo.FindAfter(ctx, userID, fromTs)
	if err != nil {
		slog.Error("failed to fetch secrets", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to fetch secrets")

		return nil, gErr
	}

	resp := &pb.ListSecretsResponse{Secrets: make([]*pb.ListSecretsResponseSecret, len(secrets)), SyncTs: syncTs}
	for i, secret := range secrets {
		rs := pb.ListSecretsResponseSecret{
			Id:        secret.ID,
			Ts:        secret.UpdatedAt,
			IsDeleted: secret.IsDeleted,
		}
		if !secret.IsDeleted {
			rs.Data = secret.Data
		}
		resp.Secrets[i] = &rs
	}

	return resp, nil
}
