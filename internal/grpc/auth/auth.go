package authgrpc

import (
	"context"
	"errors"
	"fmt"

	interrors "github.com/vitalikir156/SR_authservice/internal/errors"
	protogenerated "github.com/vitalikir156/SR_authservice/internal/grpc/protogenerated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
type Auth interface {
	Login(ctx context.Context, username string,	password string) (token string, err error)
	RegisterNewUser(ctx context.Context, username string, password string) (userID int, err error)
	CheckToken(ctx context.Context, login string, token string) (bool, error)
	UpdatePassword(ctx context.Context, login string, oldpass string, newpass string) (error)
	DeleteUser(ctx context.Context, login string) (error)
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
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

token, err:=s.auth.Login(ctx, in.Username, in.Password)
if err != nil {
	if errors.Is(err, interrors.ErrBadCred){
		return nil, status.Error(codes.Unauthenticated, "login or password is wrong")
	}
	fmt.Println(err)
	return nil, status.Error(codes.Internal, "failed to login")
}
	


	return &protogenerated.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	in *protogenerated.RegisterRequest,
) (*protogenerated.RegisterResponse, error) {
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err:=s.auth.RegisterNewUser(ctx, in.Username, in.Password)
	if err != nil {
		if errors.Is(err, interrors.ErrBusyLogin){
			return nil, status.Error(codes.AlreadyExists, "login is busy")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}
	out:=fmt.Sprint(id)

	return &protogenerated.RegisterResponse{Message: out}, nil
}

func (s *serverAPI) CheckToken(ctx context.Context, in *protogenerated.TokenRequest) (*protogenerated.TokenResponse, error){
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	valid, err:=s.auth.CheckToken(ctx, in.Username, in.Token)
	if err != nil {
		if errors.Is(err, interrors.ErrBadCred){
			return nil, status.Error(codes.Unauthenticated, "bad creds")
		}
		return nil, status.Error(codes.InvalidArgument, "bad token")
	}
	return &protogenerated.TokenResponse{Valid: valid}, nil
}

func (s *serverAPI) UpdatePassword(ctx context.Context, in *protogenerated.PassUpdRequest)(*protogenerated.PassUpdResponse, error){
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.Oldpass == "" {
		return nil, status.Error(codes.InvalidArgument, "empty old password")
	}

	if in.Newpass == "" {
		return nil, status.Error(codes.InvalidArgument, "empty new password")
	}
	valid, err:=s.auth.CheckToken(ctx, in.Username, in.Token)
	if err != nil {
		if errors.Is(err, interrors.ErrBadCred){
			return nil, status.Error(codes.Unauthenticated, "bad creds")
		}
		fmt.Println(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !valid {return nil, status.Error(codes.Unauthenticated, "bad token")}
	err=s.auth.UpdatePassword(ctx, in.Username, in.Oldpass, in.Newpass)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update password")
	}
	return &protogenerated.PassUpdResponse{Success: true}, nil
}

func (s *serverAPI) SelfDeleteUser(ctx context.Context, in *protogenerated.SelfDelRequest)(*protogenerated.SelfDelResponse, error){
	if in.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	ok, err:=s.auth.CheckToken(ctx, in.Username, in.Token)
	if err != nil {
		if errors.Is(err, interrors.ErrBadCred){
			return nil, status.Error(codes.Unauthenticated, "bad creds")
		}
		fmt.Println(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !ok {return nil, status.Error(codes.Unauthenticated, "bad token")}
	err=s.auth.DeleteUser(ctx, in.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &protogenerated.SelfDelResponse{Success: true}, nil
}