package grpcserver

import (
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/utils/localip"
	"github.com/fasdalf/train-go-gophkeeper/internal/server/db/dbrepository"
	"time"
)

func ExampleNewGrpcServer() {
	userRepo := dbrepository.UserDBRepository{}
	secretRepo := dbrepository.SecretDBRepository{}
	cKey := "mock"
	port, _ := localip.GetFreePort()
	addr := fmt.Sprintf("localhost:%d", port)
	gs := NewGrpcServer(&userRepo, &cKey, 3*time.Hour, &secretRepo)
	// covers listen error
	go ListenAndServe(gs, ":65999")
	// covers success
	go ListenAndServe(gs, addr)
	time.Sleep(50 * time.Millisecond)
	gs.GracefulStop()
	fmt.Println("interrupt signal handled")
	// covers serve error
	go ListenAndServe(gs, addr)

	// Output:
	// interrupt signal handled
}
