package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Unfield/Valdock/api"
	"github.com/Unfield/Valdock/config"
	dockerstatus "github.com/Unfield/Valdock/dockerStatus"
	"github.com/Unfield/Valdock/jobs"
	"github.com/Unfield/Valdock/logging"
	"github.com/Unfield/Valdock/version"
	"github.com/Unfield/cascade"
	"go.uber.org/zap"
)

func main() {
	logging.Init()
	defer logging.Base.Sync()

	serverLogger := logging.Base.With(zap.String("service", "server"), zap.String("env", os.Getenv("ENV")))

	cfg := config.ValdockConfig{}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8080

	cfg.PortAllocator.MinPort = 6380
	cfg.PortAllocator.MaxPort = 7380

	cfg.Docker.Instance.DefaultHostname = "127.0.0.1"

	cfgLoader := cascade.NewLoader(
		cascade.WithEnvPrefix("VALDOCK"),
		cascade.WithFlags(),
	)

	if err := cfgLoader.Load(&cfg); err != nil {
		serverLogger.Fatal("Failed to load config", zap.Error(err))
	}

	serverLogger.Info("Valdock starting", zap.String("version", version.FullVersion()), zap.String("host", cfg.Server.Host), zap.Int("port", cfg.Server.Port))

	apiRouter := api.NewAPIRouter(&cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverLogger.Info("Background services starting...")
	go dockerstatus.WatchStatusUpdates(ctx, &cfg)
	go jobs.ServeAsynqQueueServer(&cfg)

	go func() {
		err := apiRouter.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
		if err != nil {
			serverLogger.Fatal("Failed to run api router", zap.Error(err))
		}
	}()

	serverLogger.Info("Running initial status check...")
	dockerstatus.CheckInitialStatus(ctx, &cfg)
	serverLogger.Info("Initial status check completed")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	serverLogger.Info("Received termination signal, shutting down...")
	cancel()
}
