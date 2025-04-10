package swagger

import (
	"fmt"
	"strings"
)

// Generator builds OpenAPI/Swagger documentation dynamically
type Generator struct {
	parser *CommentParser
}

// NewGenerator creates a new Swagger generator
func NewGenerator() *Generator {
	return &Generator{
		parser: NewCommentParser(),
	}
}

// GenerateAPI creates an OpenAPI specification from route information
func (g *Generator) GenerateAPI(
	routes interface{},
	appName, appVersion, host, port string,
) map[string]interface{} {
	// Create base OpenAPI structure
	swagger := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       appName,
			"description": "API Documentation automatically generated by Gonyx Framework",
			"version":     appVersion,
			"contact": map[string]interface{}{
				"name":  "API Support",
				"email": "support@gonyx.io",
			},
			"license": map[string]interface{}{
				"name": "MIT",
				"url":  "https://github.com/Blocktunium/gonyx/blob/main/LICENSE",
			},
		},
		"servers": []map[string]interface{}{
			{
				"url":         fmt.Sprintf("http://%s:%s", host, port),
				"description": "API Server",
			},
		},
		"paths": map[string]interface{}{},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{},
		},
	}

	// Process all routes
	g.processRoutes(routes, swagger)

	return swagger
}

// processRoutes analyzes the routes and adds them to the OpenAPI specification
func (g *Generator) processRoutes(routes interface{}, swagger map[string]interface{}) {
	// Get the paths section of the OpenAPI spec
	paths, _ := swagger["paths"].(map[string]interface{})

	// Track paths we've already processed (to avoid duplicates)
	processedPaths := make(map[string]bool)

	// Process routes based on the specific type
	// In this case, assuming gin.RoutesInfo which has a slice of RouteInfo objects
	// with Method, Path, and Handler fields
	if ginRoutes, ok := routes.([]struct {
		Method  string
		Path    string
		Handler string
	}); ok {
		for _, route := range ginRoutes {
			// Skip Swagger routes
			if strings.HasPrefix(route.Path, "/swagger") {
				continue
			}

			// Skip if already processed this path+method combo
			pathKey := route.Path + ":" + route.Method
			if processedPaths[pathKey] {
				continue
			}
			processedPaths[pathKey] = true

			// Process this route
			g.processGinRoute(route.Method, route.Path, route.Handler, paths)
		}
	}
}

// processGinRoute adds a single Gin route to the OpenAPI paths
func (g *Generator) processGinRoute(method, path, handler string, paths map[string]interface{}) {
	// Convert Gin path to OpenAPI path (e.g., /users/:id -> /users/{id})
	openAPIPath := path
	pathParams := []map[string]interface{}{}

	// Extract path parameters
	pathSegments := strings.Split(path, "/")
	for _, segment := range pathSegments {
		if strings.HasPrefix(segment, ":") {
			paramName := strings.TrimPrefix(segment, ":")
			pathParams = append(pathParams, map[string]interface{}{
				"name":        paramName,
				"in":          "path",
				"required":    true,
				"description": fmt.Sprintf("%s parameter", paramName),
				"schema": map[string]interface{}{
					"type": "string",
				},
			})

			// Replace :param with {param} for OpenAPI format
			openAPIPath = strings.Replace(
				openAPIPath,
				":"+paramName,
				"{"+paramName+"}",
				1,
			)
		}
	}

	// Initialize path object if it doesn't exist
	if _, exists := paths[openAPIPath]; !exists {
		paths[openAPIPath] = map[string]interface{}{}
	}

	// Get or create the path object
	pathObj := paths[openAPIPath].(map[string]interface{})

	// Build basic operation info
	operation := map[string]interface{}{
		"summary":     fmt.Sprintf("%s %s", method, path),
		"description": fmt.Sprintf("Endpoint for %s %s", method, path),
		"operationId": handler,
		"parameters":  pathParams,
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "Successful operation",
			},
			"400": map[string]interface{}{
				"description": "Bad request",
			},
			"500": map[string]interface{}{
				"description": "Internal server error",
			},
		},
	}

	// Try to extract metadata from handler's comments
	if metadata, err := g.parser.GetHandlerMetadata(handler); err == nil {
		// Override with extracted metadata
		for k, v := range metadata {
			// Special handling for parameters to merge them
			if k == "parameters" {
				// If we have path parameters, merge with the extracted ones
				if len(pathParams) > 0 {
					if extractedParams, ok := v.([]map[string]interface{}); ok {
						// Create a map to track param names we've already seen
						paramNames := make(map[string]bool)
						for _, p := range pathParams {
							if name, ok := p["name"].(string); ok {
								paramNames[name] = true
							}
						}

						// Add non-duplicate params from extracted ones
						for _, p := range extractedParams {
							if name, ok := p["name"].(string); ok {
								if !paramNames[name] {
									pathParams = append(pathParams, p)
								}
							}
						}

						// Update operation with merged parameters
						operation["parameters"] = pathParams
					}
				} else {
					operation[k] = v
				}
			} else {
				operation[k] = v
			}
		}
	}

	// Add operation to path
	httpMethod := strings.ToLower(method)
	pathObj[httpMethod] = operation
}
