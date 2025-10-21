package models

import "time"

type ACLUserModel struct {
	ID          string
	InstanceID  string
	Username    string
	Password    string
	Permissions string
	Keys        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
