package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/cryptofacade"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/jwt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
)

func (s *GophKeeperServer) UserSignIn(ctx context.Context, r *pb.UserSignInRequest) (*pb.UserSignInResponse, error) {
	if r == nil || r.Login == "" || r.Password == "" {
		return nil, errors.New("login and password must be not empty")
	}

	// IRL use "golang.org/x/crypto/bcrypt" GenerateFromPassword() + CompareHashAndPassword()
	// https://medium.com/@rnp0728/secure-password-hashing-in-go-a-comprehensive-guide-5500e19e7c1f
	passHash := cryptofacade.Hash(&r.Password, s.key)
	user, err := s.userRepo.FindByLogin(ctx, r.Login)
	if err != nil || user.PassHash != passHash {
		if user != nil && user.PassHash != passHash {
			err = fmt.Errorf("password mismatch")
		}
		slog.Error("failed to sign in user", "err", err)
		gErr := status.Errorf(codes.Internal, "failed to sign in user, see server logs")

		return nil, gErr
	}

	token := jwt.BuildJWTString(user.ID, s.key, s.exp)
	return &pb.UserSignInResponse{Token: token}, nil
}
