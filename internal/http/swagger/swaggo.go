package swagger

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

// SwaggoGenerator handles programmatic swagger generation using swaggo/swag
type SwaggoGenerator struct {
	appName    string
	appVersion string
	host       string
	port       string
}

// NewSwaggoGenerator creates a new SwaggoGenerator instance
func NewSwaggoGenerator(projectRoot string) *SwaggoGenerator {
	return &SwaggoGenerator{
		appName:    "Gonyx API",
		appVersion: "1.0.0",
		host:       "localhost",
		port:       "3000",
	}
}

// GenerateSwaggerDocs generates swagger documentation using the registered swag info
func (sg *SwaggoGenerator) GenerateSwaggerDocs() ([]byte, error) {
	// Try to get swagger info from swag registry
	doc := swag.GetSwagger("swagger")
	if doc == nil {
		log.Println("Warning: No swagger info found in registry, creating basic spec")
		return sg.generateBasicSwaggerSpec()
	}

	// Get the swagger JSON from the registered doc
	swaggerJSON := doc.ReadDoc()

	// Parse and validate the JSON
	var swaggerSpec map[string]interface{}
	if err := json.Unmarshal([]byte(swaggerJSON), &swaggerSpec); err != nil {
		log.Printf("Warning: Failed to parse registered swagger JSON: %v", err)
		return sg.generateBasicSwaggerSpec()
	}

	// Update with current host and port
	sg.updateSwaggerSpec(swaggerSpec)

	// Convert back to JSON bytes
	jsonBytes, err := json.MarshalIndent(swaggerSpec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal swagger JSON: %w", err)
	}

	return jsonBytes, nil
}

// generateBasicSwaggerSpec creates a basic swagger specification
func (sg *SwaggoGenerator) generateBasicSwaggerSpec() ([]byte, error) {
	swaggerSpec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       sg.appName,
			"description": "API Documentation for Gonyx Framework",
			"version":     sg.appVersion,
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@gonyx.io",
			},
			"license": map[string]interface{}{
				"name": "MIT",
				"url":  "https://github.com/Blocktunium/gonyx/blob/main/LICENSE",
			},
		},
		"host":        fmt.Sprintf("%s:%s", sg.host, sg.port),
		"basePath":    "/",
		"schemes":     []string{"http"},
		"paths":       map[string]interface{}{},
		"definitions": map[string]interface{}{},
	}

	jsonBytes, err := json.MarshalIndent(swaggerSpec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal basic swagger JSON: %w", err)
	}

	return jsonBytes, nil
}

// GenerateFromGinRoutes generates swagger spec from Gin routes
func (sg *SwaggoGenerator) GenerateFromGinRoutes(routes gin.RoutesInfo) ([]byte, error) {
	// Start with basic swagger spec
	swaggerSpec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       sg.appName,
			"description": "API Documentation generated from Gonyx routes",
			"version":     sg.appVersion,
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@gonyx.io",
			},
			"license": map[string]interface{}{
				"name": "MIT",
				"url":  "https://github.com/Blocktunium/gonyx/blob/main/LICENSE",
			},
		},
		"host":        fmt.Sprintf("%s:%s", sg.host, sg.port),
		"basePath":    "/",
		"schemes":     []string{"http"},
		"paths":       map[string]interface{}{},
		"definitions": map[string]interface{}{},
	}

	// Add paths from routes
	paths := make(map[string]interface{})
	for _, route := range routes {
		if route.Path == "/swagger.json" || route.Path == "/swagger/*any" {
			continue // Skip swagger routes
		}

		pathItem := make(map[string]interface{})
		method := strings.ToLower(route.Method)

		pathItem[method] = map[string]interface{}{
			"summary":     fmt.Sprintf("%s %s", route.Method, route.Path),
			"description": fmt.Sprintf("Handler: %s", route.Handler),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Success",
				},
				"400": map[string]interface{}{
					"description": "Bad Request",
				},
				"500": map[string]interface{}{
					"description": "Internal Server Error",
				},
			},
		}

		paths[route.Path] = pathItem
	}

	swaggerSpec["paths"] = paths

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(swaggerSpec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal swagger JSON: %w", err)
	}

	return jsonBytes, nil
}

// updateSwaggerSpec updates host and port in existing swagger spec
func (sg *SwaggoGenerator) updateSwaggerSpec(swaggerSpec map[string]interface{}) {
	swaggerSpec["host"] = fmt.Sprintf("%s:%s", sg.host, sg.port)

	// Ensure schemes are set
	if _, exists := swaggerSpec["schemes"]; !exists {
		swaggerSpec["schemes"] = []string{"http"}
	}

	// Set base path if not exists
	if _, exists := swaggerSpec["basePath"]; !exists {
		swaggerSpec["basePath"] = "/"
	}
}

// SetHostAndPort updates the host and port for the generator
func (sg *SwaggoGenerator) SetHostAndPort(host, port string) {
	sg.host = host
	sg.port = port
}
