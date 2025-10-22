package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type UpdateInstanceContainerIDPayload struct {
	InstanceID  string `json:"instance_id"`
	ContainerID string `json:"containerID"`
}

func NewUpdateInstanceContainerIDJob(instanceID, containerID string) (*asynq.Task, error) {
	payload, err := json.Marshal(UpdateInstanceContainerIDPayload{
		InstanceID:  instanceID,
		ContainerID: containerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update containerID payload: %w", err)
	}
	return asynq.NewTask(TypeUpdateInstanceContainerID, payload), nil
}

func (jh *JobHandler) HandleUpdateInstanceContainerID(ctx context.Context, t *asynq.Task) error {
	var payload UpdateInstanceContainerIDPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload for update containerID: %w", err)
	}

	if err := jh.store.UpdateInstanceContainerID(payload.InstanceID, payload.ContainerID); err != nil {
		return fmt.Errorf("failed to update containerID for instance %s: %w", payload.InstanceID, err)
	}

	return nil
}
