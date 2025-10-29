package models

type InstanceStatus string

//String representation of the container state. Can be one of “created”, “running”, “paused”, “restarting”, “removing”, “exited”, or “dead”

const (
	StatusCreating   InstanceStatus = "creating"
	StatusCreated    InstanceStatus = "created"
	StatusRunning    InstanceStatus = "running"
	StatusPaused     InstanceStatus = "paused"
	StatusRestarting InstanceStatus = "restarting"
	StatusRemoving   InstanceStatus = "removing"
	StatusExited     InstanceStatus = "exited"
	StatusDead       InstanceStatus = "dead"
	StatusDeleted    InstanceStatus = "deleted"
	StatusFailed     InstanceStatus = "failed"
)
