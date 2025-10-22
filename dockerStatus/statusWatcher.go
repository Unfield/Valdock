package dockerstatus

import (
	"context"
	"os"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/jobs"
	"github.com/Unfield/Valdock/logging"
	"github.com/Unfield/Valdock/models"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

func WatchStatusUpdates(ctx context.Context, cfg *config.ValdockConfig) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}

	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")

	eventsChan, errChan := cli.Events(ctx, events.ListOptions{Filters: eventFilter})

	logger := logging.Base.With(
		zap.String("service", "docker-status-watcher"),
		zap.String("env", os.Getenv("ENV")),
	)

	jobClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.KV.Url})

	for {
		select {
		case <-ctx.Done():
			logger.Info("Stopping Docker status watcher")
			return

		case event := <-eventsChan:
			switch event.Action {
			case "start":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusRunning)
				jobClient.Enqueue(updateStatus)
			case "stop":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusStopped)
				jobClient.Enqueue(updateStatus)
			case "die":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusStopped)
				jobClient.Enqueue(updateStatus)
			case "pause":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusPaused)
				jobClient.Enqueue(updateStatus)
			case "unpause":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusRunning)
				jobClient.Enqueue(updateStatus)
			}

		case err := <-errChan:
			if err != nil {
				logger.Error("Docker event stream error", zap.Error(err))
				return
			}
		}
	}
}
