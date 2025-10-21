package handler

import (
	"log"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/store"
	"github.com/valkey-io/valkey-go"
)

type Handler struct {
	managementKV valkey.Client
	store        *store.Store
	pa           *store.PortAllocator
}

func NewHandler(kv valkey.Client, cfg *config.ValdockConfig) *Handler {
	st := store.NewStore(kv)

	pa, err := store.NewPortAllocator(*st, cfg.PortAllocator.MinPort, cfg.PortAllocator.MaxPort)
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		managementKV: kv,
		store:        st,
		pa:           pa,
	}
}
