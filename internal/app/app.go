package app

import (
	"context"
	"time"

	"github.com/vitalikir156/SR_authservice/internal/grpc/app/grpcsrv"

	"github.com/pkg/errors"
	"github.com/vitalikir156/SR_authservice/internal/config"
	"github.com/vitalikir156/SR_authservice/internal/repo"
	"go.uber.org/zap"
)
type App struct {
	GRPCServer *App
}

func New(
	log *zap.SugaredLogger,
	config config.AppConfig,
	tokenTTL time.Duration,
) *App {
	_, err := repo.NewRepository(context.Background(), config.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}
	if err != nil {
		panic(err)
	}

	//authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcsrv.New(log, _, config.GRPCConfig.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}