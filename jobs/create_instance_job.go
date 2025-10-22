package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Unfield/Valdock/models"
	"github.com/hibiken/asynq"
)

func NewCreateInstanceJob(instance *models.InstanceModel) (*asynq.Task, error) {
	payload, err := json.Marshal(instance)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	return asynq.NewTask(TypeCreateInstance, payload, asynq.MaxRetry(3)), nil
}

func (jh *JobHandler) HandlerCreateInstanceJob(ctx context.Context, t *asynq.Task) error {
	var p models.InstanceModel
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	containerID, err := jh.DockerService.CreateValkeyContainer(ctx, fmt.Sprintf("%s", p.Name), p.ID, p.DataPath, p.Port)
	if err != nil {
		retryCount, rcOK := asynq.GetRetryCount(ctx)
		maxRetryCount, mrOK := asynq.GetMaxRetry(ctx)
		if rcOK && mrOK && retryCount+1 >= maxRetryCount {
			followup, _ := NewUpdateStatusJob(p.ID, models.StatusFailed)
			jh.jobClient.client.Enqueue(followup)
			return fmt.Errorf("container creation failed: %w", err)
		}

		return fmt.Errorf("container creation failed: %w", err)
	}

	followup, _ := NewUpdateStatusJob(p.ID, models.StatusRunning)
	jh.jobClient.client.Enqueue(followup)

	followup2, _ := NewUpdateInstanceContainerIDJob(p.ID, containerID)
	jh.jobClient.client.Enqueue(followup2)

	return nil
}
