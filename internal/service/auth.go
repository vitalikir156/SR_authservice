package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	interrors "github.com/vitalikir156/SR_authservice/internal/errors"
	"github.com/vitalikir156/SR_authservice/internal/jwt"
	"github.com/vitalikir156/SR_authservice/internal/types"
	"go.uber.org/zap"
)

type Auth struct {
	log       *zap.SugaredLogger
	usrRepo   UserRepo
	tokenRepo TokenRepo
	tokenTTL  time.Duration
	tokenSecret string
}

type UserRepo interface {
	CreateUser(ctx context.Context, user string, pass string) (uid int, err error)
	GetUserOverLogin(ctx context.Context, login string) (types.User, error)
	UpdateUserPassword(ctx context.Context, login string, oldpass string, newpass string) error
	DeleteUser(ctx context.Context, login string) error
}

type TokenRepo interface {
	CreateToken(ctx context.Context, token types.Token) (int, error)
	UpdateToken(ctx context.Context, token types.Token) (error)
	GetToken(ctx context.Context, userid int) (types.Token, error)
	DelToken(ctx context.Context, tokenid int) error
}

func New(
	log *zap.SugaredLogger,
	userRepo UserRepo,
	tokenRepo TokenRepo,
	tokenTTL time.Duration,
	tokenSecret string,
) *Auth {
	return &Auth{
		usrRepo:   userRepo,
		tokenRepo: tokenRepo,
		log:       log,
		tokenTTL:  tokenTTL,
		tokenSecret: tokenSecret,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	user, err := a.usrRepo.GetUserOverLogin(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", interrors.ErrBadCred
		}
		return "", err
	}
	if user.Password != password {
		return "", interrors.ErrBadCred
	}
	token, err := a.tokenRepo.GetToken(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			token.Token, err = jwt.NewToken(user, a.tokenTTL, a.tokenSecret)
			if err != nil {
				return "", err
			}
			token.UID = user.ID
			_, err = a.tokenRepo.CreateToken(ctx, token)
			if err != nil {
				return "", err
			}
		}
		if err != nil {
			return "", err
		}
	}
	_, err= jwt.ValidateAccessToken(token.Token, a.tokenSecret)
	if err != nil{
		token.Token, err = jwt.NewToken(user, a.tokenTTL, a.tokenSecret)
			if err != nil {
				return "", err
			}
			token.UID = user.ID
			err = a.tokenRepo.UpdateToken(ctx, token)
			if err != nil {
				return "", err
			}
	}
	return token.Token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, login string, pass string) (int, error) {
	id, err := a.usrRepo.CreateUser(ctx, login, pass)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *Auth) CheckToken(ctx context.Context, login string, token string) (bool, error) {
	user, err := a.usrRepo.GetUserOverLogin(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, interrors.ErrBadCred
		}
		return false, err
	}
	claims, err:= jwt.ValidateAccessToken(token,  a.tokenSecret)
	if err != nil {
		return false, err
	}
if claims.UserID!=user.ID || claims.Login!=login {
	return false, err
}
	return true, nil
}

func (a *Auth) UpdatePassword(ctx context.Context, login string, oldpass string, newpass string) (error) {

return  a.usrRepo.UpdateUserPassword(ctx, login, oldpass, newpass)
}

func (a *Auth) DeleteUser(ctx context.Context, login string) (error) {

return  a.usrRepo.DeleteUser(ctx, login)
}