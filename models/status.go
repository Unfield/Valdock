package models

type InstanceStatus string

const (
	StatusCreating   InstanceStatus = "creating"
	StatusRunning    InstanceStatus = "running"
	StatusRestarting InstanceStatus = "restarting"
	StatusRecreating InstanceStatus = "recreating"
	StatusStopped    InstanceStatus = "stopped"
	StatusDied       InstanceStatus = "died"
	StatusFailed     InstanceStatus = "failed"
	StatusDeleting   InstanceStatus = "deleting"
	StatusDeleted    InstanceStatus = "deleted"
	StatusPaused     InstanceStatus = "paused"
)
