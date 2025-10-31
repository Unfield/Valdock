package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Unfield/Valdock/config"
	"github.com/Unfield/Valdock/logging"
	"github.com/Unfield/Valdock/models"
	"github.com/Unfield/Valdock/namespaces"
	"github.com/Unfield/Valdock/permissions"
	"github.com/Unfield/Valdock/response"
	"github.com/Unfield/Valdock/store"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	store  *store.Store
	logger *zap.Logger
	cfg    *config.ValdockConfig
}

func NewAuthMiddleware(valkeyClient valkey.Client, cfg *config.ValdockConfig) *AuthMiddleware {
	return &AuthMiddleware{
		store:  store.NewStore(valkeyClient),
		logger: logging.Base.With(zap.String("service", "auth-middleware"), zap.String("env", os.Getenv("ENV"))),
		cfg:    cfg,
	}
}

func (am *AuthMiddleware) Serve() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			response.SendError(c, http.StatusUnauthorized, response.Unauthorized, "unauthorized")
			am.logger.Warn("invalid authorization header provided", zap.String("client_ip", c.ClientIP()), zap.String("user_agent", c.Request.UserAgent()))
			c.Abort()
			return
		}
		authToken := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))

		if authToken == "" {
			response.SendError(c, http.StatusUnauthorized, response.Unauthorized, "unauthorized")
			am.logger.Warn("authorization token missing", zap.String("client_ip", c.ClientIP()), zap.String("user_agent", c.Request.UserAgent()))
			c.Abort()
			return
		}

		if am.cfg.API.MasterKey != "" && authToken == am.cfg.API.MasterKey {
			am.logger.Warn("master key was used", zap.String("client_ip", c.ClientIP()), zap.String("user_agent", c.Request.UserAgent()))
			c.Set("permissions", permissions.ParsePermissionString(string(permissions.RootAdmin)))
			c.Next()
			return
		}

		var apiKey models.APIKey
		if err := am.store.GetJSON(fmt.Sprintf("%s:%s", namespaces.API_KEYS, authToken), &apiKey); err != nil {
			response.SendError(c, http.StatusUnauthorized, response.Unauthorized, "unauthorized")
			am.logger.Warn("api key not found in db", zap.String("client_ip", c.ClientIP()), zap.String("user_agent", c.Request.UserAgent()))
			c.Abort()
			return
		}

		if (apiKey.DeletedAt != nil && apiKey.DeletedAt.Before(time.Now())) || !apiKey.Enabled || time.Now().After(apiKey.CreatedAt.Add(apiKey.ExpiresIn)) {
			response.SendError(c, http.StatusUnauthorized, response.Unauthorized, "unauthorized")
			am.logger.Warn("api key disabled, revoked or expired",
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.Bool("expired", time.Now().After(apiKey.CreatedAt.Add(apiKey.ExpiresIn))),
				zap.Bool("disabled", !apiKey.Enabled),
				zap.Bool("deleted", apiKey.DeletedAt != nil),
			)
			c.Abort()
			return
		}

		c.Set("permissions", apiKey.Permissions)
		c.Next()
	}
}
