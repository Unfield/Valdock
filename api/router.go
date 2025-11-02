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
	router.GET("/api/v1/health", hdlr.HealthCheckHandler)

	router.Use(middleware.RequestID())

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	v1group := router.Group("/api/v1")

	authMiddleware := middleware.NewAuthMiddleware(valkeyClient, cfg)
	v1group.Use(authMiddleware.Serve())

	// Api Key Management
	v1group.GET("/api-keys", hdlr.GetApiKeysHandler)
	v1group.POST("/api-keys", hdlr.CreateApiKeyHandler)
	v1group.GET("/api-keys/:key", hdlr.GetApiKeyHandler)
	v1group.DELETE("/api-keys/:key", hdlr.DeleteApiKeyHandler)
	v1group.POST("/api-keys/:key/disable", hdlr.UpdateApiKeyStateHandler("disable"))
	v1group.POST("/api-keys/:key/enable", hdlr.UpdateApiKeyStateHandler("enable"))
	v1group.POST("/api-keys/:key/expire", hdlr.UpdateApiKeyStateHandler("expire"))

	// Instance Management
	v1group.GET("/instances", hdlr.GetInstancesHandler)
	v1group.POST("/instances", hdlr.CreateInstanceHandler)
	v1group.GET("/instances/:id", hdlr.GetInstanceHandler)
	v1group.DELETE("/instances/:id", hdlr.DeleteInstanceHandler)
	v1group.POST("/instances/:id/start", hdlr.StartInstanceHandler)
	v1group.POST("/instances/:id/stop", hdlr.StopInstanceHandler)
	v1group.POST("/instances/:id/restart", hdlr.RestartInstanceHandler)
	//v1group.GET("/instances/:id/stats")

	// ACL Management
	v1group.GET("/instances/:id/acls", hdlr.GetACLUsers)
	v1group.POST("/instances/:id/acls", hdlr.CreateNewACLUser)
	v1group.PUT("/instances/:id/acls/:username", hdlr.UpdateAclUserHandler)
	v1group.DELETE("/instances/:id/acls/:username", hdlr.DeleteAclUserHandler)

	return router
}
