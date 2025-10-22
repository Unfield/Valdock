package docker

import (
	"fmt"
	"os"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/logging"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type DockerService struct {
	cli    *client.Client
	cfg    *config.ValdockConfig
	logger *zap.Logger
}

func NewDockerService(cfg *config.ValdockConfig) (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize docker service: %w", err)
	}

	logger := logging.GetBase().With(zap.String("service", "docker-service"), zap.String("env", os.Getenv("ENV")))

	return &DockerService{
		cli:    cli,
		cfg:    cfg,
		logger: logger,
	}, nil
}
