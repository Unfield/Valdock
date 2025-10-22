package jobs

import (
	"os"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/logging"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func ServeAsynqQueueServer(cfg *config.ValdockConfig) {
	logger := logging.GetBase().With(zap.String("service", "asynq-server"), zap.String("env", os.Getenv("ENV")))

	asynqLogger := logging.NewAsynqLogger(logger)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.KV.Url},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			Logger: asynqLogger,
		},
	)

	jobHandler := NewJobHandler(cfg)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeCreateInstance, jobHandler.HandlerCreateInstanceJob)
	mux.HandleFunc(TypeUpdateInstanceStatus, jobHandler.HandleUpdateInstanceStatus)
	mux.HandleFunc(TypeUpdateInstanceContainerID, jobHandler.HandleUpdateInstanceContainerID)
	mux.HandleFunc(TypeDeleteInstance, jobHandler.HandlerDeleteInstanceJob)

	if err := srv.Run(mux); err != nil {
		logger.Fatal("asynq server failed", zap.Error(err))
	}
}
