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

// CorsMiddlewareConfig defines the configuration for CORS middleware
// Based on gin-contrib/cors configuration reference
type CorsMiddlewareConfig struct {
	// AllowAllOrigins sets whether all origins are allowed
	AllowAllOrigins bool `json:"allow_all_origins"`

	// AllowOrigins is a list of origins that may access the resource
	AllowOrigins []string `json:"allow_origins"`

	// AllowMethods is a list of methods the client is allowed to use
	AllowMethods []string `json:"allow_methods"`

	// AllowPrivateNetwork indicates if the resource allows requests from private network
	AllowPrivateNetwork bool `json:"allow_private_network"`

	// AllowHeaders is a list of request headers the client is allowed to use
	AllowHeaders []string `json:"allow_headers"`

	// AllowCredentials indicates if the request can include user credentials
	AllowCredentials bool `json:"allow_credentials"`

	// ExposeHeaders indicates which headers are safe to expose
	ExposeHeaders []string `json:"expose_headers"`

	// MaxAge indicates how long (in seconds) the results can be cached
	MaxAge time.Duration `json:"max_age"`

	// AllowWildcard indicates if wildcards are allowed in origins
	AllowWildcard bool `json:"allow_wildcard"`

	// AllowBrowserExtensions indicates if browser extensions are allowed
	AllowBrowserExtensions bool `json:"allow_browser_extensions"`

	// CustomSchemas is a list of custom schemas like tauri://
	CustomSchemas []string `json:"custom_schemas"`

	// AllowWebSockets indicates if WebSocket origins are allowed
	AllowWebSockets bool `json:"allow_websockets"`

	// AllowFiles indicates if file:// origins are allowed
	AllowFiles bool `json:"allow_files"`

	// OptionsResponseStatusCode sets the status code for OPTIONS responses
	OptionsResponseStatusCode int `json:"options_response_status_code"`
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
