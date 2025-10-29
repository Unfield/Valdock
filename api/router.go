package api

import (
	"log"
	"os"
	"time"

	"github.com/Unfield/Valdock/api/handler"
	"github.com/Unfield/Valdock/api/middleware"
	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/logging"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

func NewAPIRouter(cfg *config.ValdockConfig) *gin.Engine {
	valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{cfg.KV.Url}})
	if err != nil {
		log.Fatal(err)
	}

	hdlr := handler.NewHandler(valkeyClient, cfg)

	logger := logging.Base.With(zap.String("service", "api"), zap.String("env", os.Getenv("ENV")))

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(middleware.RequestID())

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	v1group := router.Group("/api/v1")

	v1group.GET("/health", hdlr.HealthCheckHandler)

	// Instance Management

	v1group.GET("/instances", hdlr.GetInstancesHandler)
	v1group.POST("/instances", hdlr.CreateInstanceHandler)
	v1group.GET("/instances/:id", hdlr.GetInstanceHandler)
	v1group.DELETE("/instances/:id", hdlr.DeleteInstanceHandler)
	v1group.POST("/instances/:id/start")
	v1group.POST("/instances/:id/stop")
	v1group.POST("/instances/:id/restart")
	//v1group.GET("/instances/:id/stats")

	// ACL Management

	//v1group.GET("/instances/:id/acls")
	//v1group.POST("/instances/:id/acls")
	//v1group.PUT("/instances/:id/acls/:username")
	//v1group.DELETE("/instances/:id/acls/:username")

	return router
}
