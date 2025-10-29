package models

import (
	"time"
)

type InstanceModel struct {
	ID              string     `json:"id"`
	ContainerID     string     `json:"container_id"`
	DataPath        string     `json:"data_path"`
	Name            string     `json:"name"`
	Port            int        `json:"port"`
	Status          string     `json:"status"`
	PrimaryHostname string     `json:"primary_hostname"`
	ConfigTemplate  string     `json:"config_template"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
