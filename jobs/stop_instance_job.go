package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

func NewStopInstanceJob(container_id string) (*asynq.Task, error) {
	payload, err := json.Marshal(container_id)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	return asynq.NewTask(TypeStopInstance, payload, asynq.MaxRetry(3)), nil
}

func (jh *JobHandler) HandlerStopInstanceJob(ctx context.Context, t *asynq.Task) error {
	var containerID string
	if err := json.Unmarshal(t.Payload(), &containerID); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	jh.DockerService.StopDockerContainer(ctx, containerID)

	return nil
}
