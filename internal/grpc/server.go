package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	// Сгенерированный код
	protogenerated "github.com/vitalikir156/SR_authservice/internal/grpc/protogenerated"
)
type serverAPI struct {
    protogenerated.UnimplementedAuthServer // Хитрая штука, о ней ниже
    auth Auth
}

type Auth interface {
    Login(
        ctx context.Context,
        login string,
        password string,
    ) (token string, err error)
    RegisterNewUser(
        ctx context.Context,
        login string,
        password string,
    ) (userID int64, err error)
}
func Register(gRPCServer *grpc.Server, auth Auth) {  
    protogenerated.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})  
}

func (s *serverAPI) Login(
    ctx context.Context,
    in *protogenerated.LoginRequest,
) (*protogenerated.LoginResponse, error) {
	if in.Login == "" {
        return nil, fmt.Errorf("email is required")
    }

    if in.Password == "" {
        return nil, fmt.Errorf("password is required")
    }

    if in.GetAppId() == 0 {
        return nil, fmt.Errorf("app_id is required")
    }

    token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
    if err != nil {
        // Ошибку auth.ErrInvalidCredentials мы создадим ниже
        if errors.Is(err, auth.ErrInvalidCredentials) {
            return nil, fmt.Errorf("invalid email or password")
        }

        return nil, fmt.Errorf("failed to login")
    }

    return &protogenerated.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
    ctx context.Context,
    in *protogenerated.RegisterRequest,
) (*protogenerated.RegisterResponse, error) {
	
    if in.Login == "" {
        return nil, fmt.Errorf("email is required")
    }

    if in.Password == "" {
        return nil, fmt.Errorf("password is required")
    }

    uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
    if err != nil {
        // Ошибку storage.ErrUserExists мы создадим ниже
        if errors.Is(err, storage.ErrUserExists) {
            return nil, fmt.Errorf("user already exists")
        }

        return nil, fmt.Errorf("failed to register user")
    }

    return &protogenerated.RegisterResponse{UserId: uid}, nil
}