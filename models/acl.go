package models

import (
	"strings"
	"time"
)

type ACLUserModel struct {
	ID              string     `json:"id"`
	InstanceID      string     `json:"instance_id"`
	Username        string     `json:"username"`
	Enabled         bool       `json:"enabled"`
	NoPassword      bool       `json:"no_password"`
	PasswordHashes  []string   `json:"password_hashes"`
	KeyPatterns     []string   `json:"key_patterns"`
	ChannelPatterns []string   `json:"channel_patterns"`
	AllowedCommands []string   `json:"allowed_commands"`
	DeniedCommands  []string   `json:"denied_commands"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

func (u ACLUserModel) ToValkeyACL() string {
	parts := []string{"user", u.Username}

	if u.Enabled {
		parts = append(parts, "on")
	} else {
		parts = append(parts, "off")
	}

	for _, hash := range u.PasswordHashes {
		parts = append(parts, "#"+hash)
	}

	for _, p := range u.KeyPatterns {
		parts = append(parts, "~"+p)
	}

	for _, c := range u.ChannelPatterns {
		parts = append(parts, "&"+c)
	}

	for _, cmd := range u.AllowedCommands {
		parts = append(parts, "+"+cmd)
	}

	for _, cmd := range u.DeniedCommands {
		parts = append(parts, "-"+cmd)
	}

	return strings.Join(parts, " ")
}

func MakeACLFile(users []ACLUserModel) string {
	var res string
	for _, user := range users {
		res = res + user.ToValkeyACL()
	}
	return res
}
