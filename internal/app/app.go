package app

import (
	"FreeMusic/internal/config"
	"FreeMusic/internal/repository/mongodb"
	"FreeMusic/internal/server"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	handler "FreeMusic/internal/delivery/http/v1"
	"FreeMusic/internal/repository"
	"FreeMusic/internal/service"
)

func Run(configPath string) {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	config, err := config.InitConfig(configPath)
	if err != nil {
		logrus.Fatalf("error initializing configs: %v", err)
	}

	fileStorage, err := mongodb.NewMongoFileStorage(config)
	if err != nil {
		logrus.Fatalf("error initializing configs: %v", err)
	}

	repos := repository.NewRepository(fileStorage)

	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := server.NewServer(config)
	logrus.Print("FreeMusic Started")
	if err := srv.Run(handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %v", err.Error())
	}
	logrus.Print("FreeMusic Shutting Down")

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occured on server shutting down: %v", err.Error())
	}
}
