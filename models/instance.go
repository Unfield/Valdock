package models

import (
	"time"
)

type InstanceModel struct {
	ID             string
	ContainerID    string
	DataPath       string
	Name           string
	Port           int
	Status         string
	ConfigTemplate string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
