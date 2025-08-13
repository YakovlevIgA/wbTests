package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wbService/internal/api"
	"wbService/internal/config"
	"wbService/internal/logger"
	"wbService/internal/repo"
	"wbService/internal/service"
)

func main() {
	var cfg config.AppConfig
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
	}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logg, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	repository := repo.NewRepository()

	serviceInstance := service.NewService(repository, logg)

	app := api.NewRouters(&api.Routers{Service: serviceInstance}, logg)

	go func() {
		logg.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			logg.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logg.Info("Shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(); err != nil {
		logg.Error(errors.Wrap(err, "failed to shutdown fiber app"))
	}

	<-shutdownCtx.Done()
	logg.Info("Server stopped")
}
