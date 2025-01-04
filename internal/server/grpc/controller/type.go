package controller

import (
	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"time"
)

// GophKeeperServer поддерживает все необходимые методы сервера.
type GophKeeperServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedGophKeeperServer
	userRepo   repository.UserRepository
	key        *string
	exp        time.Duration
	secretRepo repository.SecretRepository
}

func NewGophKeeperServer(
	userRepo repository.UserRepository,
	key *string,
	exp time.Duration,
	secretRepo repository.SecretRepository,
) *GophKeeperServer {
	return &GophKeeperServer{
		UnimplementedGophKeeperServer: pb.UnimplementedGophKeeperServer{},
		userRepo:                      userRepo,
		key:                           key,
		exp:                           exp,
		secretRepo:                    secretRepo,
	}
}
