package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

func (d *DockerService) StopDockerContainer(ctx context.Context, id string) {
	timeout := 10

	if err := d.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		d.logger.Error("failed to stop container", zap.String("container_id", id), zap.Error(err))
		return
	}

	// maybe we should notify the status job to set status to exited but in this case it will be picked up by the docker-status service
}
