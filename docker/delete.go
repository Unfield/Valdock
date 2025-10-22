package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

func (d *DockerService) DeleteValkeyContainer(ctx context.Context, containerID string) error {
	if containerID == "" {
		d.logger.Error("containerID missing")
		return fmt.Errorf("containerID missing")
	}

	if err := d.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{RemoveVolumes: true, Force: true}); err != nil {
		d.logger.Error("failed to delete container", zap.Error(err))
		return err
	}

	return nil
}
