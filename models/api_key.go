package models

import (
	"time"

	"github.com/Unfield/Valdock/permissions"
)

type APIKey struct {
	ID      string `json:"id"`
	Value   string `json:"key"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	// instances:read instances:write instances:delete acl:read acl:write acl:delete
	Permissions []permissions.Permission `json:"permissions"`
	ExpiresIn   time.Duration            `json:"expires_in"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	DeletedAt   *time.Time               `json:"deleted_at omitempty"`
}
