package app

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vitalikir156/SR_authservice/internal/config"
	grpcsrv "github.com/vitalikir156/SR_authservice/internal/grpc/app"
	"github.com/vitalikir156/SR_authservice/internal/repo"
	auth "github.com/vitalikir156/SR_authservice/internal/service"
	"go.uber.org/zap"
)
type App struct {
	GRPCServer *grpcsrv.App
}

func New(
	log *zap.SugaredLogger,
	config config.AppConfig,
	tokenTTL time.Duration,
) *App {
	storage, err := repo.NewRepository(context.Background(), config.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, config.GRPCConfig.Timeout, config.GRPCConfig.Secret)

	grpcApp := grpcsrv.New(log, authService, config.GRPCConfig.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}