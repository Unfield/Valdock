package jobs

import (
	"log"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/docker"
	"github.com/Unfield/Valdock/store"
	"github.com/valkey-io/valkey-go"
)

const (
	TypeCreateInstance            = "instance:create"
	TypeStartInstance             = "instance:start"
	TypeStopInstance              = "instance:stop"
	TypeRestartInstance           = "instance:restart"
	TypeDeleteInstance            = "instance:delete"
	TypeUpdateInstanceStatus      = "instance:updateStatus"
	TypeUpdateInstanceContainerID = "instance:updateContainerID"
)

type JobHandler struct {
	DockerService *docker.DockerService
	store         *store.Store
	jobClient     *JobClient
}

func NewJobHandler(cfg *config.ValdockConfig) *JobHandler {
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.KV.Url}})
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStore(valkeyClient)

	dService, err := docker.NewDockerService(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &JobHandler{
		DockerService: dService,
		store:         store,
		jobClient:     NewJobClient(cfg.KV.Url),
	}
}
