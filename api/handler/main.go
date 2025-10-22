package handler

import (
	"log"
	"os"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/jobs"
	"github.com/Unfield/Valdock/logging"
	"github.com/Unfield/Valdock/store"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

type Handler struct {
	cfg          *config.ValdockConfig
	managementKV valkey.Client
	store        *store.Store
	pa           *store.PortAllocator
	jobClient    *jobs.JobClient
	logger       *zap.Logger
}

func NewHandler(kv valkey.Client, cfg *config.ValdockConfig) *Handler {
	st := store.NewStore(kv)
	jobClient := jobs.NewJobClient(cfg.KV.Url)

	pa, err := store.NewPortAllocator(*st, cfg.PortAllocator.MinPort, cfg.PortAllocator.MaxPort)
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.Base.With(zap.String("service", "api"), zap.String("env", os.Getenv("ENV")))

	return &Handler{
		cfg:          cfg,
		managementKV: kv,
		store:        st,
		pa:           pa,
		jobClient:    jobClient,
		logger:       logger,
	}
}
