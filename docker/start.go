package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

func (d *DockerService) StartDockerContainer(ctx context.Context, id string) {
	if err := d.cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		d.logger.Error("failed to start container", zap.String("container_id", id), zap.Error(err))
		return
	}

	// maybe we should notify the status job to set status to running but in this case it will be picked up by the docker-status service
}
