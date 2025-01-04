package grpcserver

import (
	"github.com/fasdalf/train-go-gophkeeper/internal/server/model/repository"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"

	pb "github.com/fasdalf/train-go-gophkeeper/internal/common/proto/gophkeeperserver"
	gc "github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/controller"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/interceptors"
)

func NewGrpcServer(
	userRepo repository.UserRepository,
	key *string,
	exp time.Duration,
	secretRepo repository.SecretRepository,
) *grpc.Server {
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.NewSlogInterceptor()),
		grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor()),
		// TODO: ##@@ add auth interceptor and throw all bellow away
		// IRL use github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/realip
		// or at least peer.FromContext(ctx).Addr
		//grpc.ChainUnaryInterceptor(interceptors.NewValidateIPInterceptor(tr)),
		//grpc.ChainUnaryInterceptor(interceptors.NewDecryptBodyInterceptor(decryptionKey)),
		//grpc.ChainUnaryInterceptor(interceptors.NewValidateHashInterceptor(key)),
		//grpc.ChainUnaryInterceptor(interceptors.NewRespondWithHashInterceptor(key)),
		//// TODO: ##@@ add gzip+rsa home-made compressor and validate metadata for "grpc-accept-encoding": ["gziprsa"]
		//// This one grpc.ForceServerCodecV2(),
		//// or use interceptor with decrypt from message body
	}
	s := grpc.NewServer(options...)
	mServer := gc.NewGophKeeperServer(userRepo, key, exp, secretRepo)
	pb.RegisterGophKeeperServer(s, mServer)

	return s
}

// ListenAndServe intended to use with errgroup.
func ListenAndServe(server *grpc.Server, addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if err = server.Serve(listen); err != nil {
		return err
	}
	return nil
}
