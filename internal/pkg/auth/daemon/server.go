package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/tremendouscan/bifrost/api/protobuf-spec/authpb"
	"github.com/tremendouscan/bifrost/internal/pkg/auth/config"
	"github.com/tremendouscan/bifrost/internal/pkg/auth/endpoint"
	"github.com/tremendouscan/bifrost/internal/pkg/auth/logging"
	"github.com/tremendouscan/bifrost/internal/pkg/auth/service"
	"github.com/tremendouscan/bifrost/internal/pkg/auth/transport"
)

// authSvc service.AuthService.
var isInit bool

func ServerRun() error {
	if !isInit {
		return errors.New("service related configuration not initialized")
	}

	Log(DEBUG, "Listening system call signal")
	go ListenSignal(signalChan)
	Log(DEBUG, "Listened system call signal")

	ctx := context.Background()

	var svc service.Service
	svc = AuthConf.AuthService
	svc = logging.LoggingMiddleware(config.KitLogger)(svc)

	eps := endpoint.MakeAuthEndpoints(svc)

	handler := transport.NewAuthServer(ctx, eps)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", AuthConf.AuthService.Port))
	if lisErr != nil {
		return lisErr
	}
	defer lis.Close()

	gRPCServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(gRPCServer, handler)
	svrErrChan := make(chan error, 1)
	go func() {
		svrErr := gRPCServer.Serve(lis)
		Log(NOTICE, "bifrost-auth service is running on %s", lis.Addr())
		svrErrChan <- svrErr
	}()

	var stopErr error
	select {
	case s := <-signalChan:
		if s == 9 {
			// fmt.Println("stopping...")
			Log(DEBUG, "stopping...")
			break
		}
		Log(DEBUG, "stop signal error")
	case stopErr = <-svrErrChan:
		break
	}
	// fmt.Println("gRPC Server stopping...")
	gRPCServer.Stop()
	return stopErr
}
