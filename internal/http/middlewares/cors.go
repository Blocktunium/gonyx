package middlewares

import (
	"time"

	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware creates a CORS middleware based on the provided configuration
// If config is nil, it returns the default CORS middleware
func CorsMiddleware(config *types.CorsMiddlewareConfig) gin.HandlerFunc {
	if config == nil {
		return cors.Default()
	}

	corsConfig := cors.Config{
		AllowAllOrigins:           config.AllowAllOrigins,
		AllowOrigins:              config.AllowOrigins,
		AllowMethods:              config.AllowMethods,
		AllowPrivateNetwork:       config.AllowPrivateNetwork,
		AllowHeaders:              config.AllowHeaders,
		AllowCredentials:          config.AllowCredentials,
		ExposeHeaders:             config.ExposeHeaders,
		MaxAge:                    time.Duration(config.MaxAge) * time.Second,
		AllowWildcard:             config.AllowWildcard,
		AllowBrowserExtensions:    config.AllowBrowserExtensions,
		CustomSchemas:             config.CustomSchemas,
		AllowWebSockets:           config.AllowWebSockets,
		AllowFiles:                config.AllowFiles,
		OptionsResponseStatusCode: int(config.OptionsResponseStatusCode),
	}

	return cors.New(corsConfig)
}
