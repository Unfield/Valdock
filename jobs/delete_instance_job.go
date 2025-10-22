package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/hibiken/asynq"
)

type DeleteInstancePayload struct {
	ID string `json:"id"`
}

func NewDeleteInstanceJob(id string) (*asynq.Task, error) {
	payload, err := json.Marshal(DeleteInstancePayload{
		ID: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delete instance payload: %w", err)
	}
	return asynq.NewTask(TypeDeleteInstance, payload), nil
}

func (jh *JobHandler) HandlerDeleteInstanceJob(ctx context.Context, t *asynq.Task) error {
	var p DeleteInstancePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	var instance models.InstanceModel
	if err := jh.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.INSTANCES, p.ID), &instance); err != nil {
		return fmt.Errorf("container deletion failed: %w", err)
	}

	err := jh.DockerService.DeleteValkeyContainer(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("container deletion failed: %w", err)
	}

	followup, _ := NewUpdateStatusJob(p.ID, models.StatusDeleted)
	jh.jobClient.client.Enqueue(followup)

	return nil
}
