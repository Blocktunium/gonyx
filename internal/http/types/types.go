package types

import (
	"time"
)

// SwaggerConfig - defines the config for Swagger documentation.
type SwaggerConfig struct {
	Enabled bool `json:"enabled"`
}

// LoggerMiddlewareConfig - defines the config for middleware.
type LoggerMiddlewareConfig struct {
	Format       string `json:"format"`
	TimeFormat   string `json:"time_format"`
	TimeZone     string `json:"time_zone"`
	TimeInterval int    `json:"time_interval"`
	Output       string `json:"output"`
}

type FaviconMiddlewareConfig struct {
	File         string `json:"file"`
	URL          string `json:"url"`
	CacheControl string `json:"cache_control"`
}

type GinServerConfig struct {
	ListenAddress string   `json:"addr"`
	Name          string   `json:"name"`
	Versions      []string `json:"versions"`
	SupportStatic bool     `json:"support_static"`
	Config        struct {
		ReadTimeout    time.Duration `json:"read_timeout"`
		WriteTimeout   time.Duration `json:"write_timeout"`
		RequestMethods []string      `json:"request_methods"`
	} `json:"conf"`
	Middlewares struct {
		Order []string `json:"order"`
	} `json:"middlewares"`
	Static struct {
		Prefix string `json:"prefix"`
		Root   string `json:"root"`
	} `json:"static"`
	Swagger SwaggerConfig `json:"swagger"`
}
