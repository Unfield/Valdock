package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"

	configtemplates "github.com/Unfield/Valdock/configTemplates"
	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

func (d *DockerService) CreateValkeyContainer(
	ctx context.Context,
	name, instanceID, dataPath string,
	port int,
) (string, error) {
	if dataPath == "" {
		d.logger.Error("dataPath invalid")
		return "", fmt.Errorf("dataPath invalid")
	}

	var config models.ConfigModel
	if err := d.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.CONFIG, instanceID), &config); err != nil {
		d.logger.Warn("no config found! creating new default config...")
		config = *configtemplates.NewDefaultConfig(
			instanceID,
			configtemplates.DefaultConfigOptions{MaxMemoryMB: 256, MaxClients: 10000},
		)
		if err := d.store.SetJSON(fmt.Sprintf("%s:%s", namespaces.CONFIG, instanceID), config); err != nil {
			d.logger.Error("failed to write config (no existing config found)", zap.Error(err))
			return "", fmt.Errorf("failed to write config (no existing config found)")
		}
	}

	aclUsers, err := d.store.ListACLs()
	if err != nil || len(aclUsers) < 1 {
		d.logger.Warn("no acl users found!")
	}

	configStr := config.ToConf()
	aclUsersStr := models.MakeACLFile(aclUsers)

	img := "valkey/valkey:latest"
	_, err = d.cli.ImagePull(ctx, img, image.PullOptions{})
	if err != nil {
		d.logger.Error("failed to pull image", zap.Error(err))
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	portMapping := nat.Port(fmt.Sprintf("%d/tcp", 6379))
	hostMapping := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: fmt.Sprintf("%d", port),
	}

	cConfig := &container.Config{
		Image: img,
		Cmd:   []string{"valkey-server", "/etc/valkey/valkey.conf"},
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

	containerResponse, err := d.cli.ContainerCreate(ctx, cConfig, hostConfig, networking, nil, instanceID)
	if err != nil {
		d.logger.Error("container creation failed", zap.Error(err))
		return "", fmt.Errorf("container creation failed: %w", err)
	}

	d.logger.Info("container created",
		zap.String("container_id", containerResponse.ID),
		zap.String("instance_id", instanceID),
	)

	confTar, err := makeTarFromString("/etc/valkey/valkey.conf", configStr)
	if err != nil {
		d.logger.Error("failed to make config tar", zap.Error(err))
		return "", fmt.Errorf("failed to make config tar: %w", err)
	}

	aclTar, err := makeTarFromString("/data/users.acl", aclUsersStr)
	if err != nil {
		d.logger.Error("failed to make acl users tar", zap.Error(err))
		return "", fmt.Errorf("failed to make acl users tar: %w", err)
	}

	if err := d.cli.CopyToContainer(
		ctx,
		containerResponse.ID,
		"/",
		confTar,
		container.CopyToContainerOptions{AllowOverwriteDirWithFile: true},
	); err != nil {
		return "", fmt.Errorf("failed to copy config: %w", err)
	}

	if err := d.cli.CopyToContainer(
		ctx,
		containerResponse.ID,
		"/",
		aclTar,
		container.CopyToContainerOptions{AllowOverwriteDirWithFile: true},
	); err != nil {
		return "", fmt.Errorf("failed to copy acl users: %w", err)
	}

	if err := d.cli.ContainerStart(ctx, containerResponse.ID, container.StartOptions{}); err != nil {
		d.logger.Error("failed to start container", zap.Error(err))
		return "", fmt.Errorf("start failed: %w", err)
	}

	d.logger.Info("container started",
		zap.String("container_id", containerResponse.ID),
		zap.String("instance_name", name),
		zap.String("instance_id", instanceID),
	)
	return containerResponse.ID, nil
}

func makeTarFromString(path, content string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: path,
		Mode: 0644,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(content)); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}
