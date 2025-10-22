package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

func (d *DockerService) CreateValkeyContainer(ctx context.Context, name, instanceID, dataPath string, port int) (string, error) {
	if dataPath == "" {
		d.logger.Error("dataPath invalid")
		return "", fmt.Errorf("dataPath invalid")
	}

	img := "valkey/valkey:latest"
	_, err := d.cli.ImagePull(ctx, img, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	portMapping := nat.Port(fmt.Sprintf("%d/tcp", 6379))
	hostMapping := nat.PortBinding{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", port)}

	config := &container.Config{
		Image: img,
	}

	hostConfig := &container.HostConfig{
		Binds:        []string{fmt.Sprintf("%s:/data", dataPath)},
		PortBindings: nat.PortMap{portMapping: []nat.PortBinding{hostMapping}},
	}

	networkName := d.cfg.Docker.Instance.Net
	networking := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {},
		},
	}

	containerResponse, err := d.cli.ContainerCreate(ctx, config, hostConfig, networking, nil, instanceID)
	if err != nil {
		return "", fmt.Errorf("container create failed: %w", err)
	}

	d.logger.Info("container created", zap.String("container_id", containerResponse.ID), zap.String("instance_id", instanceID))

	if err := d.cli.ContainerStart(ctx, containerResponse.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("start failed: %w", err)
	}

	d.logger.Info("container started", zap.String("container_id", containerResponse.ID), zap.String("instance_name", name), zap.String("instance_id", instanceID))
	return containerResponse.ID, nil
}
