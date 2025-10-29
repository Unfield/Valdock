package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type RestartInstanceContainerIDPayload struct {
	InstanceID  string `json:"instance_id"`
	ContainerID string `json:"containerID"`
}

func NewRestartInstanceJob(containerIDPayload RestartInstanceContainerIDPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(containerIDPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	return asynq.NewTask(TypeRestartInstance, payload, asynq.MaxRetry(3)), nil
}

func (jh *JobHandler) HandlerRestartInstanceJob(ctx context.Context, t *asynq.Task) error {
	var containerIDPayload RestartInstanceContainerIDPayload
	if err := json.Unmarshal(t.Payload(), &containerIDPayload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	jh.DockerService.RestartDockerContainer(ctx, containerIDPayload.ContainerID)

	return nil
}
