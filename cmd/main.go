package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/vitalikir156/SR_authservice/internal/config"
	customLogger "github.com/vitalikir156/SR_authservice/internal/logger"
	"github.com/vitalikir156/SR_authservice/internal/repo"
)

func main() {
	loadenv := pflag.BoolP("loadenv", "e", false, "load .env file")
	pflag.Parse()

	if *loadenv {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to load env file"))
		}
	}

	// Загружаем конфигурацию из переменных окружения
	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	// Инициализация логгера
	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	_, err = repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to initialize repository"))
	}




	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
