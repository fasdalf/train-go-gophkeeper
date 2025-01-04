package controller

import (
	"context"
	"errors"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/cryptofacade"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/jwt"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
)

// UpdateMetrics updates metrics
func (s *GophKeeperServer) UserSignUp(ctx context.Context, r *pb.UserSignUpRequest) (*pb.UserSignUpResponse, error) {
	if r == nil || r.Login == "" || r.Password == "" {
		return nil, errors.New("login and password must be not empty")
	}

	// IRL use "golang.org/x/crypto/bcrypt" GenerateFromPassword() + CompareHashAndPassword()
	// https://medium.com/@rnp0728/secure-password-hashing-in-go-a-comprehensive-guide-5500e19e7c1f
	passHash := cryptofacade.Hash(&r.Password, s.key)
	user, err := s.userRepo.Create(ctx, r.Login, passHash)
	if err != nil {
		slog.Error("failed to create user", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to create user, see server logs")
		if errors.Is(err, repository.ErrUserExists) {
			gErr = status.Errorf(codes.AlreadyExists, "failed to create user: %v", err)
		}

		return nil, gErr
	}

	token, err := jwt.BuildJWTString(user.ID, s.key, s.exp)
	if err != nil {
		slog.Error("failed to build jwt string", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to respond with jwt, see server logs")

		return nil, gErr
	}

	return &pb.UserSignUpResponse{Token: token}, nil
}
