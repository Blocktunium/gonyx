package middlewares

import (
	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/fufuok/favicon"
	"github.com/gin-gonic/gin"
)

func FaviconMiddleware(config types.FaviconMiddlewareConfig) gin.HandlerFunc {
	return favicon.New(favicon.Config{
		File:         config.File,
		CacheControl: config.CacheControl,
	})
}
