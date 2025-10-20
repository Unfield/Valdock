package handler

import "github.com/valkey-io/valkey-go"

type Handler struct {
	managementKV valkey.Client
}

func NewHandler(kv valkey.Client) *Handler {
	return &Handler{
		managementKV: kv,
	}
}
