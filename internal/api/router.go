package api

import (
	"log"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/internal/api/handler"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
)

func NewAPIRouter(cfg *config.ValdockConfig) *gin.Engine {
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.KV.Url}})
	if err != nil {
		log.Fatal(err)
	}

	hdlr := handler.NewHandler(valkeyClient)

	router := gin.Default()

	router.GET("/api/v1/health", hdlr.HealthCheckHandler)
	router.GET("/api/v1/test", hdlr.TestHandler)

	return router
}
