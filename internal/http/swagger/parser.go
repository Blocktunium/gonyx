package swagger

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// CommentParser extracts Swagger documentation from Go source files
type CommentParser struct {
	// Cache of parsed files to avoid repeated I/O operations
	parsedFiles map[string]*ast.File
	// Map of fully qualified handler names to their metadata
	handlerDocs map[string]map[string]interface{}
}

// NewCommentParser creates a new CommentParser instance
func NewCommentParser() *CommentParser {
	return &CommentParser{
		parsedFiles: make(map[string]*ast.File),
		handlerDocs: make(map[string]map[string]interface{}),
	}
}

// GetHandlerMetadata extracts Swagger metadata from the given handler
func (p *CommentParser) GetHandlerMetadata(handlerName string) (map[string]interface{}, error) {
	// Check if we've already parsed this handler
	if metadata, exists := p.handlerDocs[handlerName]; exists {
		return metadata, nil
	}

	// Default metadata
	metadata := map[string]interface{}{
		"summary":     handlerName,
		"description": "Automatically generated endpoint",
	}

	// Parse the handler name to get package path and function name
	parts := strings.Split(handlerName, ".")
	if len(parts) < 2 {
		return metadata, fmt.Errorf("invalid handler name format: %s", handlerName)
	}

	// Last part is the function name
	funcName := parts[len(parts)-1]

	// Try to locate the function's source file
	sourceFile, err := p.findSourceFile(handlerName)
	if err != nil {
		return metadata, err
	}

	// Parse the source file if not already parsed
	file, ok := p.parsedFiles[sourceFile]
	if !ok {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
		if err != nil {
			return metadata, fmt.Errorf("failed to parse source file: %v", err)
		}
		p.parsedFiles[sourceFile] = parsedFile
		file = parsedFile
	}

	// Find the function in the AST
	var funcDecl *ast.FuncDecl
	for _, decl := range file.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == funcName {
			funcDecl = fd
			break
		}
	}

	if funcDecl == nil {
		return metadata, fmt.Errorf("function not found: %s in %s", funcName, sourceFile)
	}

	// Extract Swagger comments
	if funcDecl.Doc != nil {
		docText := funcDecl.Doc.Text()
		p.parseSwaggerComments(docText, metadata)
	}

	// Cache the result
	p.handlerDocs[handlerName] = metadata
	return metadata, nil
}

// parseSwaggerComments extracts Swagger annotations from comment text
func (p *CommentParser) parseSwaggerComments(commentText string, metadata map[string]interface{}) {
	// Extract @Summary
	summaryRegex := regexp.MustCompile(`@Summary\s+(.+)`)
	if matches := summaryRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		metadata["summary"] = strings.TrimSpace(matches[1])
	}

	// Extract @Description
	descRegex := regexp.MustCompile(`@Description\s+(.+)`)
	if matches := descRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		metadata["description"] = strings.TrimSpace(matches[1])
	}

	// Extract @Tags
	tagsRegex := regexp.MustCompile(`@Tags\s+(.+)`)
	if matches := tagsRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		tags := strings.Split(matches[1], ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		metadata["tags"] = tags
	}

	// Extract @Accept
	acceptRegex := regexp.MustCompile(`@Accept\s+(.+)`)
	if matches := acceptRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		metadata["consumes"] = []string{strings.TrimSpace(matches[1])}
	}

	// Extract @Produce
	produceRegex := regexp.MustCompile(`@Produce\s+(.+)`)
	if matches := produceRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		metadata["produces"] = []string{strings.TrimSpace(matches[1])}
	}

	// Extract @Param
	// Format: @Param name in type required description
	paramRegex := regexp.MustCompile(`@Param\s+(\w+)\s+(\w+)\s+(\w+)\s+(true|false)\s+"([^"]+)"`)
	paramMatches := paramRegex.FindAllStringSubmatch(commentText, -1)

	if len(paramMatches) > 0 {
		params := []map[string]interface{}{}
		for _, match := range paramMatches {
			if len(match) > 5 {
				param := map[string]interface{}{
					"name":        match[1],
					"in":          match[2], // path, query, header, body, form
					"required":    match[4] == "true",
					"description": match[5],
					"schema": map[string]interface{}{
						"type": match[3], // string, integer, number, boolean, array, object
					},
				}
				params = append(params, param)
			}
		}

		// Add to existing parameters or set as new
		if existingParams, ok := metadata["parameters"].([]map[string]interface{}); ok {
			metadata["parameters"] = append(existingParams, params...)
		} else {
			metadata["parameters"] = params
		}
	}

	// Extract @Success responses
	// Format: @Success code {type} model "description"
	successRegex := regexp.MustCompile(`@Success\s+(\d+)\s+{(\w+)}\s+(\w+)\s+"([^"]+)"`)
	successMatches := successRegex.FindAllStringSubmatch(commentText, -1)

	// Extract @Failure responses
	// Format: @Failure code {type} model "description"
	failureRegex := regexp.MustCompile(`@Failure\s+(\d+)\s+{(\w+)}\s+(\w+)\s+"([^"]+)"`)
	failureMatches := failureRegex.FindAllStringSubmatch(commentText, -1)

	responses := map[string]interface{}{}

	// Process success responses
	for _, match := range successMatches {
		if len(match) > 4 {
			code := match[1]
			respType := match[2] // object, array, string, etc.
			model := match[3]    // model name or schema type
			description := match[4]

			contentSchema := map[string]interface{}{}

			// If it's a reference to a model
			if respType == "object" && !isBuiltInType(model) {
				contentSchema = map[string]interface{}{
					"$ref": "#/components/schemas/" + model,
				}
			} else {
				contentSchema = map[string]interface{}{
					"type": respType,
				}

				// If it's an array of models
				if respType == "array" && !isBuiltInType(model) {
					contentSchema["items"] = map[string]interface{}{
						"$ref": "#/components/schemas/" + model,
					}
				} else if respType == "array" {
					contentSchema["items"] = map[string]interface{}{
						"type": model,
					}
				}
			}

			responses[code] = map[string]interface{}{
				"description": description,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": contentSchema,
					},
				},
			}
		}
	}

	// Process failure responses
	for _, match := range failureMatches {
		if len(match) > 4 {
			code := match[1]
			respType := match[2]
			model := match[3]
			description := match[4]

			contentSchema := map[string]interface{}{}

			if respType == "object" && !isBuiltInType(model) {
				contentSchema = map[string]interface{}{
					"$ref": "#/components/schemas/" + model,
				}
			} else {
				contentSchema = map[string]interface{}{
					"type": respType,
				}

				if respType == "array" && !isBuiltInType(model) {
					contentSchema["items"] = map[string]interface{}{
						"$ref": "#/components/schemas/" + model,
					}
				} else if respType == "array" {
					contentSchema["items"] = map[string]interface{}{
						"type": model,
					}
				}
			}

			responses[code] = map[string]interface{}{
				"description": description,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": contentSchema,
					},
				},
			}
		}
	}

	// Add default responses if none are provided
	if len(successMatches) == 0 && len(failureMatches) == 0 {
		responses["200"] = map[string]interface{}{
			"description": "Successful operation",
		}
		responses["400"] = map[string]interface{}{
			"description": "Bad request",
		}
		responses["500"] = map[string]interface{}{
			"description": "Internal server error",
		}
	}

	metadata["responses"] = responses

	// Extract @Security
	// Format: @Security SecuritySchemeName
	securityRegex := regexp.MustCompile(`@Security\s+(\w+)`)
	if matches := securityRegex.FindStringSubmatch(commentText); len(matches) > 1 {
		securityScheme := strings.TrimSpace(matches[1])
		// Map common security scheme names
		switch strings.ToLower(securityScheme) {
		case "apikeyauth":
			metadata["security"] = []map[string]interface{}{
				{"bearerAuth": []string{}},
			}
		case "bearerauth":
			metadata["security"] = []map[string]interface{}{
				{"bearerAuth": []string{}},
			}
		default:
			metadata["security"] = []map[string]interface{}{
				{securityScheme: []string{}},
			}
		}
	}

	// Extract @Router (for validation, though we get path from Gin routes)
	// Format: @Router /path [method]
	routerRegex := regexp.MustCompile(`@Router\s+([^\s]+)\s+\[([^\]]+)\]`)
	if matches := routerRegex.FindStringSubmatch(commentText); len(matches) > 2 {
		// This could be used for validation or override
		metadata["_router_path"] = strings.TrimSpace(matches[1])
		metadata["_router_method"] = strings.TrimSpace(matches[2])
	}
}

// findSourceFile attempts to locate the source file containing the given handler
func (p *CommentParser) findSourceFile(handlerName string) (string, error) {
	// Get the executable path
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("couldn't determine caller")
	}

	// Get the project root directory (assuming internal/http/swagger/parser.go structure)
	projectRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..")

	// Handle different types of handler names
	// Case 1: Full package path like "github.com/Blocktunium/gonyx/internal/http.HandlerFunc"
	fullPathRegex := regexp.MustCompile(`^(github\.com/\w+/\w+)/(.+?)\.(\w+)$`)
	if matches := fullPathRegex.FindStringSubmatch(handlerName); len(matches) > 3 {
		// repoPath := matches[1]
		packagePath := matches[2]
		funcName := matches[3]

		// Calculate the path to the package
		packageDir := filepath.Join(projectRoot, packagePath)
		return findFunctionInDir(packageDir, funcName)
	}

	// Case 2: Package identifier like "pkg/api.HandlerFunc"
	relativePathRegex := regexp.MustCompile(`^(\w+(?:/\w+)*?)\.(\w+)$`)
	if matches := relativePathRegex.FindStringSubmatch(handlerName); len(matches) > 2 {
		packagePath := matches[1]
		funcName := matches[2]

		// Look in different potential directories
		for _, baseDir := range []string{
			projectRoot,
			filepath.Join(projectRoot, "internal"),
			filepath.Join(projectRoot, "pkg"),
		} {
			packageDir := filepath.Join(baseDir, packagePath)
			if filePath, err := findFunctionInDir(packageDir, funcName); err == nil {
				return filePath, nil
			}
		}
	}

	// Case 3: Just the func name (most likely in the same package)
	// This might be the case with the handler string from Gin routes
	if strings.Count(handlerName, ".") == 0 {
		// Try in the current HTTP package directory
		httpDir := filepath.Join(filepath.Dir(thisFile), "..")
		if filePath, err := findFunctionInDir(httpDir, handlerName); err == nil {
			return filePath, nil
		}

		// Try in common API directories
		for _, dir := range []string{
			filepath.Join(projectRoot, "internal", "api"),
			filepath.Join(projectRoot, "pkg", "api"),
			filepath.Join(projectRoot, "api"),
		} {
			if filePath, err := findFunctionInDir(dir, handlerName); err == nil {
				return filePath, nil
			}
		}
	}

	return "", fmt.Errorf("could not locate source file for handler: %s", handlerName)
}

// findFunctionInDir searches for a function in all Go files in the given directory
func findFunctionInDir(dir string, funcName string) (string, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", dir)
	}

	// Look at all Go files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return "", err
	}

	// Parse each file and look for the function
	for _, file := range files {
		fset := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue // Skip files with parse errors
		}

		// Check if the function exists in this file
		for _, decl := range parsedFile.Decls {
			if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == funcName {
				return file, nil
			}
		}
	}

	return "", fmt.Errorf("function %s not found in directory %s", funcName, dir)
}

// isBuiltInType checks if a type is a built-in Go type
func isBuiltInType(typeName string) bool {
	builtInTypes := map[string]bool{
		"string":  true,
		"int":     true,
		"int64":   true,
		"float64": true,
		"bool":    true,
		"byte":    true,
		"rune":    true,
		"error":   true,
		"any":     true,
	}
	return builtInTypes[typeName]
}
