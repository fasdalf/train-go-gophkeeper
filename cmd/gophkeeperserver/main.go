package main

import (
	"context"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/printbuild"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/config"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/db/dbrepository"
	dbstorage "github.com/fasdalf/train-go-gophkeeper/internal/server/db/storage"
	grpcserver "github.com/fasdalf/train-go-gophkeeper/internal/server/grpc/server"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	(&printbuild.Data{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}).Print()
	// TODO: ##@@ improve N
	// * Make printbuild work with --version arg only
	// * quit after output
	// * move to printbuild package
	// * Write script to fill them on build.

	// TODO: ##@@ improve N+1
	// * Write cross-platform build script for client and server.

	ctx := context.Background()
	config.InitConfig()
	cfg := config.GetConfig()
	db := dbstorage.NewDBStorage(ctx, cfg.StorageDBDSN, cfg.StorageDBPrefix)
	if done, err := db.MigrateByFlag(); err != nil || done {
		if err != nil {
			slog.Error("can not use the DB", "error", err)
			panic(err)
		}

		if done {
			os.Exit(0)
		}
	}
	if err := db.CheckVersion(ctx); err != nil {
		slog.Error("can not use the DB", "error", err)

		panic(err)
	}
	userRepo := dbrepository.NewUserDBRepository(db)
	secretRepo := dbrepository.NewSecretDBRepository(db)
	gs := grpcserver.NewGrpcServer(userRepo, &cfg.CryptoKey, cfg.TokenExp, secretRepo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	slog.Info("starting server", "addr", cfg.GRPCAddr)
	go grpcserver.ListenAndServe(gs, cfg.GRPCAddr)

	<-quit
	slog.Info("interrupt signal received")
	signal.Stop(quit)
	gs.GracefulStop()
}
