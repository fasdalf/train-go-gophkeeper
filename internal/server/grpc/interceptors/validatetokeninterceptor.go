package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/fasdalf/train-go-gophkeeper/internal/server/jwt"
)

const (
	methodPrefix = "/gophkeeperserver.GophKeeper/User"
)

var (
	userIdKey = struct{}{}
)

func NewValidateTokenInterceptor(key *string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !strings.HasPrefix(info.FullMethod, methodPrefix) {
			token := strings.TrimPrefix(getFirstMd(ctx, jwt.AuthHeader), jwt.AuthPrefix)
			userID, err := jwt.GetUserID(&token, key)
			if err != nil {
				slog.Error("token is empty or invalid", "token", token)
				msg := fmt.Sprintf("%s token is empty or invalid: %s", jwt.AuthHeader, token)
				return nil, status.Error(codes.Unauthenticated, msg)
			}
			ctx = context.WithValue(ctx, userIdKey, userID)
		}
		return handler(ctx, req)
	}
}

// getFirstMd silently gets first metadata value as string
func getFirstMd(ctx context.Context, key string) (s string) {
	values := metadata.ValueFromIncomingContext(ctx, key)
	if len(values) > 0 {
		s = values[0]
	}
	return
}
