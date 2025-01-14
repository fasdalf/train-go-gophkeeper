package controller

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/contextkeys"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
)

func (s *GophKeeperServer) UpdateSecret(ctx context.Context, r *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	// Protected with NewValidateTokenInterceptor()
	userID := ctx.Value(contextkeys.UserIDKey).(int64)

	secret, err := s.secretRepo.FindById(ctx, r.GetId())
	if err != nil || secret.UserId != userID {
		if err == nil {
			err = errors.New("bad boy")
		}
		slog.Error("failed to create secret", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to create secret")
		return nil, gErr
	}

	if secret.UpdatedAt != r.GetOldTs() {
		err = repository.ErrSecretTimeNotMatch
	}

	if err == nil {
		secret, err = s.secretRepo.Update(ctx, r.GetId(), r.GetData(), r.GetOldTs())
	}

	if errors.Is(err, repository.ErrSecretTimeNotMatch) {
		return &pb.UpdateSecretResponse{Status: pb.TimeStampStatus_TS_MISMATCH}, nil
	}

	resp := &pb.UpdateSecretResponse{
		Status: pb.TimeStampStatus_SUCCESS,
		NewTs:  secret.UpdatedAt,
	}
	return resp, nil
}
