package jobs

import (
	"fmt"

	"github.com/Unfield/Valdock/models"
	"github.com/hibiken/asynq"
)

type JobClient struct {
	client *asynq.Client
}

func NewJobClient(addr string) *JobClient {
	return &JobClient{
		client: asynq.NewClient(asynq.RedisClientOpt{Addr: addr}),
	}
}

func (jc *JobClient) EnqueueCreateInstance(instance *models.InstanceModel) error {
	task, err := NewCreateInstanceJob(instance)
	if err != nil {
		return fmt.Errorf("failed to enqueue create instance job: %w", err)
	}

	_, err = jc.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (jc *JobClient) EnqueueDeleteInstance(instanceID string) error {
	task, err := NewDeleteInstanceJob(instanceID)
	if err != nil {
		return fmt.Errorf("failed to enqueue delete instance job: %w", err)
	}

	_, err = jc.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (jc *JobClient) EnqueueStartInstance(containerID, instanceID string) error {
	task, err := NewStartInstanceJob(StartInstanceContainerIDPayload{InstanceID: instanceID, ContainerID: containerID})
	if err != nil {
		return fmt.Errorf("failed to enqueue start instance job: %w", err)
	}

	_, err = jc.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (jc *JobClient) EnqueueStopInstance(containerID, instanceID string) error {
	task, err := NewStopInstanceJob(StopInstanceContainerIDPayload{InstanceID: instanceID, ContainerID: containerID})
	if err != nil {
		return fmt.Errorf("failed to enqueue stop instance job: %w", err)
	}

	_, err = jc.client.Enqueue(task, asynq.Queue("default"))
	return err
}

func (jc *JobClient) EnqueueRestartInstance(containerID, instanceID string) error {
	task, err := NewRestartInstanceJob(RestartInstanceContainerIDPayload{InstanceID: instanceID, ContainerID: containerID})
	if err != nil {
		return fmt.Errorf("failed to enqueue restart instance job: %w", err)
	}

	_, err = jc.client.Enqueue(task, asynq.Queue("default"))
	return err
}
