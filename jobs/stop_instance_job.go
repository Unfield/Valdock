package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type StopInstanceContainerIDPayload struct {
	InstanceID  string `json:"instance_id"`
	ContainerID string `json:"containerID"`
}

func NewStopInstanceJob(containerIDPayload StopInstanceContainerIDPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(containerIDPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	return asynq.NewTask(TypeStopInstance, payload, asynq.MaxRetry(3)), nil
}

func (jh *JobHandler) HandlerStopInstanceJob(ctx context.Context, t *asynq.Task) error {
	var containerIDPayload StopInstanceContainerIDPayload
	if err := json.Unmarshal(t.Payload(), &containerIDPayload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	jh.DockerService.StopDockerContainer(ctx, containerIDPayload.ContainerID)

	return nil
}
