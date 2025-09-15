package swagger

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
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

// HandlerInfo holds information about a handler function and its annotations
type HandlerInfo struct {
	FuncName    string
	FilePath    string
	Annotations map[string]string
}

// extractHandlerInfo extracts handler function information from gin route
func (sg *SwaggoGenerator) extractHandlerInfo(handlerName string) *HandlerInfo {
	// Try to get function pointer from handler name
	// This is a simplified approach - in real scenarios you might need more sophisticated parsing
	parts := strings.Split(handlerName, ".")
	if len(parts) < 2 {
		return nil
	}

	return &HandlerInfo{
		FuncName:    parts[len(parts)-1],
		FilePath:    "",
		Annotations: make(map[string]string),
	}
}

// parseSwagAnnotations parses swag annotations from source file
func (sg *SwaggoGenerator) parseSwagAnnotations(filePath, funcName string) map[string]string {
	annotations := make(map[string]string)

	// Parse the source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Failed to parse file %s: %v", filePath, err)
		return annotations
	}

	// Find the function and its comments
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == funcName {
			if fn.Doc != nil {
				for _, comment := range fn.Doc.List {
					text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
					if strings.HasPrefix(text, "@") {
						parts := strings.SplitN(text, " ", 2)
						if len(parts) == 2 {
							key := strings.TrimPrefix(parts[0], "@")
							annotations[key] = strings.TrimSpace(parts[1])
						}
					}
				}
			}
			return false // Found the function, stop searching
		}
		return true
	})

	return annotations
}

// getHandlerFunction attempts to get the actual function pointer from gin handler
func (sg *SwaggoGenerator) getHandlerFunction(handlerName string) (string, string) {
	// Parse handler name to extract function information
	parts := strings.Split(handlerName, ".")
	if len(parts) < 2 {
		return "", ""
	}

	funcName := parts[len(parts)-1]

	// Try to resolve file path using runtime information
	// This approach attempts to find the function by name pattern matching
	filePath := sg.findHandlerFile(funcName, handlerName)

	return funcName, filePath
}

// findHandlerFile searches for the source file containing the handler function
func (sg *SwaggoGenerator) findHandlerFile(funcName, fullHandlerName string) string {
	// Extract package path from handler name
	parts := strings.Split(fullHandlerName, ".")
	if len(parts) < 2 {
		return ""
	}

	// Get package name (everything except the last part which is function name)
	packageParts := parts[:len(parts)-1]
	packageName := strings.Join(packageParts, ".")

	// Common patterns for Go source files
	searchPaths := []string{
		// Try different common patterns
		fmt.Sprintf("**/%s.go", strings.ToLower(funcName)),
		fmt.Sprintf("**/*handler*.go"),
		fmt.Sprintf("**/*controller*.go"),
		fmt.Sprintf("**/*api*.go"),
	}

	// Search in project directory structure
	for _, pattern := range searchPaths {
		if filePath := sg.searchFileByPattern(pattern, funcName); filePath != "" {
			return filePath
		}
	}

	// Fallback: try to find by package structure
	if packagePath := sg.resolvePackageToPath(packageName); packagePath != "" {
		// Look for Go files in the package directory
		files := sg.findGoFiles(packagePath)
		for _, file := range files {
			if sg.containsFunction(file, funcName) {
				return file
			}
		}
	}

	return ""
}

// searchFileByPattern searches for files matching pattern that contain the function
func (sg *SwaggoGenerator) searchFileByPattern(pattern, funcName string) string {
	// This would implement file pattern matching
	// For now, return empty - would need filepath.Walk or similar
	return ""
}

// resolvePackageToPath attempts to resolve package name to file system path
func (sg *SwaggoGenerator) resolvePackageToPath(packageName string) string {
	// Try to resolve package name to actual path
	// This is a simplified implementation

	// Remove common module prefixes
	cleanPackage := packageName
	if strings.Contains(packageName, "/") {
		parts := strings.Split(packageName, "/")
		// Find the last meaningful part
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "" && !strings.HasPrefix(parts[i], "func") {
				cleanPackage = parts[i]
				break
			}
		}
	}

	// Common package to directory mappings
	commonMappings := map[string]string{
		"handlers":    "handlers",
		"controllers": "controllers",
		"api":         "api",
		"routes":      "routes",
		"endpoints":   "endpoints",
	}

	if dir, exists := commonMappings[cleanPackage]; exists {
		return filepath.Join(".", dir)
	}

	return ""
}

// findGoFiles finds all .go files in a directory
func (sg *SwaggoGenerator) findGoFiles(dir string) []string {
	var files []string

	// This would use filepath.Walk to find .go files
	// Simplified implementation for now
	return files
}

// containsFunction checks if a file contains the specified function
func (sg *SwaggoGenerator) containsFunction(filePath, funcName string) bool {
	// Parse the file and check if it contains the function
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return false
	}

	// Search for the function
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == funcName {
			found = true
			return false
		}
		return true
	})

	return found
}

// GenerateFromGinRoutes generates swagger spec from Gin routes with annotation parsing
func (sg *SwaggoGenerator) GenerateFromGinRoutes(routes gin.RoutesInfo) ([]byte, error) {
	// Start with basic swagger spec
	swaggerSpec := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"title":       sg.appName,
			"description": "API Documentation generated from Gonyx routes with annotations",
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

	// Add paths from routes with annotation parsing
	paths := make(map[string]interface{})
	for _, route := range routes {
		if route.Path == "/swagger.json" || route.Path == "/swagger/*any" {
			continue // Skip swagger routes
		}

		// Extract handler information
		funcName, filePath := sg.getHandlerFunction(route.Handler)

		// Parse annotations if we have file path
		var annotations map[string]string
		if filePath != "" {
			annotations = sg.parseSwagAnnotations(filePath, funcName)
		} else {
			annotations = make(map[string]string)
		}

		// Create path item with annotation data
		method := strings.ToLower(route.Method)
		pathItem := make(map[string]interface{})

		// Build operation from annotations or defaults
		operation := sg.buildOperationFromAnnotations(route, annotations)
		pathItem[method] = operation

		// Handle multiple operations for the same path
		if existingPath, exists := paths[route.Path]; exists {
			if pathMap, ok := existingPath.(map[string]interface{}); ok {
				pathMap[method] = operation
				continue
			}
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

// buildOperationFromAnnotations builds swagger operation from parsed annotations
func (sg *SwaggoGenerator) buildOperationFromAnnotations(route gin.RouteInfo, annotations map[string]string) map[string]interface{} {
	operation := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "Success",
			},
		},
	}

	// Use annotations if available, otherwise use defaults
	if summary, exists := annotations["Summary"]; exists {
		operation["summary"] = summary
	} else {
		operation["summary"] = fmt.Sprintf("%s %s", route.Method, route.Path)
	}

	if description, exists := annotations["Description"]; exists {
		operation["description"] = description
	} else {
		operation["description"] = fmt.Sprintf("Handler: %s", route.Handler)
	}

	if tags, exists := annotations["Tags"]; exists {
		operation["tags"] = strings.Split(tags, ",")
	}

	if accept, exists := annotations["Accept"]; exists {
		operation["consumes"] = []string{accept}
	}

	if produce, exists := annotations["Produce"]; exists {
		operation["produces"] = []string{produce}
	}

	// Parse parameters from annotations
	if params := sg.parseParameters(annotations); len(params) > 0 {
		operation["parameters"] = params
	}

	// Parse additional responses
	if responses := sg.parseResponses(annotations); len(responses) > 0 {
		if existingResponses, ok := operation["responses"].(map[string]interface{}); ok {
			for code, response := range responses {
				existingResponses[code] = response
			}
		}
	}

	return operation
}

// parseParameters parses parameter annotations
func (sg *SwaggoGenerator) parseParameters(annotations map[string]string) []map[string]interface{} {
	var parameters []map[string]interface{}

	for key, value := range annotations {
		if strings.HasPrefix(key, "Param") {
			// Parse parameter format: "name in type required description"
			parts := strings.Fields(value)
			if len(parts) >= 4 {
				param := map[string]interface{}{
					"name":        parts[0],
					"in":          parts[1],
					"type":        parts[2],
					"required":    parts[3] == "true",
					"description": strings.Join(parts[4:], " "),
				}
				parameters = append(parameters, param)
			}
		}
	}

	return parameters
}

// parseResponses parses response annotations
func (sg *SwaggoGenerator) parseResponses(annotations map[string]string) map[string]interface{} {
	responses := make(map[string]interface{})

	for key, value := range annotations {
		if strings.HasPrefix(key, "Success") || strings.HasPrefix(key, "Failure") {
			// Parse response format: "code {type} description"
			parts := strings.Fields(value)
			if len(parts) >= 2 {
				code := parts[0]
				description := strings.Join(parts[1:], " ")

				// Remove type information for now (would need schema parsing)
				if strings.Contains(description, "{") && strings.Contains(description, "}") {
					start := strings.Index(description, "}")
					if start != -1 && start+1 < len(description) {
						description = strings.TrimSpace(description[start+1:])
					}
				}

				responses[code] = map[string]interface{}{
					"description": description,
				}
			}
		}
	}

	return responses
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
