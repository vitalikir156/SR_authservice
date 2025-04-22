package authgrpc

import (
	"context"

	protogenerated "github.com/vitalikir156/SR_authservice/internal/grpc/protogenerated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Auth interface {
	Login(
		ctx context.Context,
		username string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		username string,
		password string,
	) (userID int64, err error)
//	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	protogenerated.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	protogenerated.RegisterAuthServiceServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	in *protogenerated.LoginRequest,
) (*protogenerated.LoginResponse, error) {
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	


	return &protogenerated.LoginResponse{Token: "token1"}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *protogenerated.RegisterRequest,
) (*protogenerated.RegisterResponse, error) {
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	

	return &protogenerated.RegisterResponse{Message: "lollloloo"}, nil
}

