package dockerstatus

import (
	"context"
	"os"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/jobs"
	"github.com/Unfield/Valdock/logging"
	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/store"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/hibiken/asynq"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

func CheckInitialStatus(ctx context.Context, cfg *config.ValdockConfig) {
	logger := logging.Base.With(
		zap.String("service", "docker-status-watcher"),
		zap.String("env", os.Getenv("ENV")),
	)

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer func() {
		_ = cli.Close()
	}()

	if err != nil {
		logger.Error("failed to connect to docker", zap.Error(err))
		return
	}

	valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.KV.Url}})
	if err != nil {
		logger.Error("failed to connect to valkey server", zap.Error(err))
		return
	}
	defer valkeyClient.Close()

	store := store.NewStore(valkeyClient)

	jobClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.KV.Url})

	eventFilter := filters.NewArgs()
	eventFilter.Add("type", "container")

	managedContainers, err := store.ListInstances()
	if len(managedContainers) == 0 {
		logger.Warn("no containers to check")
		return
	}

	for _, c := range managedContainers {
		cInfo, err := cli.ContainerInspect(ctx, c.ContainerID)
		if err != nil {
			logger.Warn("failed to check container stats", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.Error(err))
			continue
		}
		logger.Warn("checking status of container", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID))
		switch cInfo.State.Status {
		case container.StateRunning:
			if c.Status != "running" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusRunning)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusRunning)
			}
			continue
		case container.StateRestarting:
			if c.Status != "restarting" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusRestarting)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusRestarting)
			}
			continue
		case container.StateRemoving:
			if c.Status != "removing" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusRemoving)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusRemoving)
			}
			continue
		case container.StatePaused:
			if c.Status != "paused" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusPaused)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusPaused)
			}
			continue
		case container.StateExited:
			if c.Status != "exited" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusExited)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusExited)
			}
			continue
		case container.StateDead:
			if c.Status != "dead" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusDead)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusDead)
			}
			continue
		case container.StateCreated:
			if c.Status != "created" {
				logger.Warn("switching status", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.String("status", string(models.StatusCreated)))
				enqueStatusUpdate(logger, jobClient, c, models.StatusCreated)
			}
			continue
		}
	}
}

func enqueStatusUpdate(logger *zap.Logger, jobClient *asynq.Client, c models.InstanceModel, status models.InstanceStatus) {
	updateStatus, err := jobs.NewUpdateStatusJob(c.ID, status)
	if err != nil {
		logger.Warn("failed to create status update job", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.Error(err))
		return
	}
	_, err = jobClient.Enqueue(updateStatus)
	if err != nil {
		logger.Warn("failed to enque status update job", zap.String("container_id", c.ContainerID), zap.String("instance_id", c.ID), zap.Error(err))
		return
	}
}

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
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusExited)
				jobClient.Enqueue(updateStatus)
			case "die":
				updateStatus, _ := jobs.NewUpdateStatusJob(event.Actor.Attributes["name"], models.StatusDead)
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
