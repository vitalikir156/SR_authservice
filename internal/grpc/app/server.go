package grpcsrv

import (
	"fmt"
	"log/slog"
	"net"

	// Сгенерированный код

	//"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	authgrpc "github.com/vitalikir156/SR_authservice/internal/grpc/auth"

	protogenerated "github.com/vitalikir156/SR_authservice/internal/grpc/protogenerated"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)
type UserServiceServer struct {
    protogenerated.UnimplementedAuthServiceServer
}
/*
// Метод GetUser
func (s *UserServiceServer) Login(ctx context.Context, req *protogenerated.LoginRequest) (*protogenerated.LoginResponse, error) {
    log.Printf("Получен запрос на пользователя с ID: %s", req.Username)

    return &protogenerated.LoginResponse{
        Token:  "Testtoken123",
    }, nil
}
*/

type App struct {
	log        *zap.SugaredLogger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app.
func New(
	log *zap.SugaredLogger,
	authService authgrpc.Auth,
	port int,
) *App {

	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
