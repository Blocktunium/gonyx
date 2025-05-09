package types

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"testing"
)

func TestServerConfig_UnmarshalJson(t *testing.T) {
	// Test input data
	input := []byte(`{
      "addr":                   ":3000",
      "name":                   "s1",
	  "versions":               ["v1", "v2"],
	  "support_static":         true,
      "conf": {
        "server_header": "",
        "strict_routing": false,
        "case_sensitive": false,
        "unescape_path": false,
        "etag": false,
        "body_limit": 4194304,
        "concurrency": 262144,
        "read_timeout": -1,
        "write_timeout": -1,
        "idle_timeout": -1,
        "read_buffer_size": 4096,
        "write_buffer_size": 4096,
        "compressed_file_suffix": ".gz",
        "get_only": false,
        "disable_keepalive": false,
        "network": "tcp",
        "enable_print_routes": true,
        "attach_error_handler": true,
		"request_methods": ["ALL"]
      },
      "middlewares": {
		"order": ["logger", "favicon"],
		"favicon": {
            "file": "./favicon.ico",
            "url": "/favicon.ico",
            "cache_control": "public, max-age=31536000"
        }
	  },
      "static": {
        "prefix": "/",
        "root": "./public",
        "config": {
          "compress": false,
          "byte_range": false,
          "browse": false,
          "download": false,
          "index": "index.html",
          "cache_duration": 10,
          "max_age": 0
        }
      }
    }`)

	// Test expected result
	v := make([]string, 2)
	v[0] = "v1"
	v[1] = "v2"
	expected := ServerConfig{
		ListenAddress: ":3000",
		Name:          "s1",
		Versions:      v,
		SupportStatic: true,
		Config: struct {
			ServerHeader         string   `json:"server_header"`
			StrictRouting        bool     `json:"strict_routing"`
			CaseSensitive        bool     `json:"case_sensitive"`
			UnescapePath         bool     `json:"unescape_path"`
			Etag                 bool     `json:"etag"`
			BodyLimit            int      `json:"body_limit"`
			Concurrency          int      `json:"concurrency"`
			ReadTimeout          int      `json:"read_timeout"`
			WriteTimeout         int      `json:"write_timeout"`
			IdleTimeout          int      `json:"idle_timeout"`
			ReadBufferSize       int      `json:"read_buffer_size"`
			WriteBufferSize      int      `json:"write_buffer_size"`
			CompressedFileSuffix string   `json:"compressed_file_suffix"`
			GetOnly              bool     `json:"get_only"`
			DisableKeepalive     bool     `json:"disable_keepalive"`
			Network              string   `json:"network"`
			EnablePrintRoutes    bool     `json:"enable_print_routes"`
			AttachErrorHandler   bool     `json:"attach_error_handler"`
			RequestMethods       []string `json:"request_methods"`
		}{
			ServerHeader:         "",
			StrictRouting:        false,
			CaseSensitive:        false,
			UnescapePath:         false,
			Etag:                 false,
			BodyLimit:            4194304,
			Concurrency:          262144,
			ReadTimeout:          -1,
			WriteTimeout:         -1,
			IdleTimeout:          -1,
			ReadBufferSize:       4096,
			WriteBufferSize:      4096,
			CompressedFileSuffix: ".gz",
			GetOnly:              false,
			DisableKeepalive:     false,
			Network:              "tcp",
			EnablePrintRoutes:    true,
			AttachErrorHandler:   true,
			RequestMethods:       []string{"ALL"},
		},
		Middlewares: struct {
			Order []string `json:"order"`
		}{Order: []string{"logger", "favicon"}},
		Static: struct {
			Prefix string       `json:"prefix"`
			Root   string       `json:"root"`
			Config fiber.Static `json:"config"`
		}{Prefix: "/", Root: "./public", Config: fiber.Static{
			Compress:      false,
			ByteRange:     false,
			Browse:        false,
			Download:      false,
			Index:         "index.html",
			CacheDuration: 10,
			MaxAge:        0,
		}},
	}

	// Unmarshal the input data into a ServerConfig instance
	var config ServerConfig
	err := json.Unmarshal(input, &config)
	if err != nil {
		t.Errorf("Error unmarshaling JSON: %v", err)
	}

	// Check if the result is what we expected
	if !reflect.DeepEqual(config, expected) {
		t.Errorf("Expected %+v, got %+v", expected, config)
	}
}
