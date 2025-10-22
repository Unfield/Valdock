package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unfield/Valdock/models"
	"github.com/hibiken/asynq"
)

type UpdateInstanceStatusPayload struct {
	ID     string                `json:"id"`
	Status models.InstanceStatus `json:"status"`
}

func NewUpdateStatusJob(id string, status models.InstanceStatus) (*asynq.Task, error) {
	payload, err := json.Marshal(UpdateInstanceStatusPayload{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update status payload: %w", err)
	}
	return asynq.NewTask(TypeUpdateInstanceStatus, payload), nil
}

func (jh *JobHandler) HandleUpdateInstanceStatus(ctx context.Context, t *asynq.Task) error {
	var payload UpdateInstanceStatusPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload for update status: %w", err)
	}

	if err := jh.store.UpdateInstanceStatus(payload.ID, payload.Status); err != nil {
		return fmt.Errorf("failed to update status for instance %s: %w", payload.ID, err)
	}

	return nil
}
