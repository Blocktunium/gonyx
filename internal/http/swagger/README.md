# Dynamic Swagger Documentation for Gonyx

This package provides dynamic Swagger/OpenAPI documentation for the Gonyx framework's HTTP module. It automatically generates API documentation from your route definitions and enhances it by parsing standard Swagger annotations in your handler function comments.

## Features

- Dynamically generates OpenAPI 3.0 documentation from your registered routes
- No pre-generation step required - documentation is created at runtime
- Parses standard Swagger annotations from your handler functions' comments
- Automatically extracts path parameters from route definitions
- Supports all standard Swagger annotations

## Usage

### 1. Enable Swagger in your HTTP configuration

```json
{
  "http": {
    "s1": {
      "addr": "localhost:8080",
      "name": "s1",
      "versions": ["v1"],
      "swagger": {
        "enabled": true
      }
      ...
    }
  }
}
```

### 2. Add Swagger annotations to your handler functions

```go
// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "User found"
// @Success 404 {object} ErrorResponse "User not found"
// @Router /v1/users/{id} [get]
func GetUserByID(c *gin.Context) {
    // Implementation...
}
```

### 3. Access your Swagger documentation

Once your server is running with Swagger enabled, you can access:

- Swagger UI: `http://localhost:8080/swagger/index.html`
- Swagger JSON: `http://localhost:8080/swagger/doc.json`

## Supported Annotations

The parser supports the following standard Swagger annotations:

| Annotation | Example | Description |
|------------|---------|-------------|
| @Summary | `@Summary Get user by ID` | Brief summary of the endpoint |
| @Description | `@Description Detailed description...` | Longer description of what the endpoint does |
| @Tags | `@Tags users,admin` | Tags for grouping endpoints |
| @Accept | `@Accept json` | Content types accepted by the endpoint |
| @Produce | `@Produce json` | Content types produced by the endpoint |
| @Param | `@Param id path string true "User ID"` | Parameter definition |
| @Success | `@Success 200 {object} UserResponse "Success"` | Success response definition |
| @Failure | `@Failure 400 {object} ErrorResponse "Bad request"` | Failure response definition |
| @Router | `@Router /users/{id} [get]` | Endpoint path and method |

## How It Works

1. When an HTTP request is made to `/swagger/doc.json`, the Gonyx framework dynamically generates the Swagger documentation.
2. The generator examines all registered routes in your application.
3. For each route, it extracts path parameters and builds basic documentation.
4. It then attempts to locate and parse your handler functions' source code.
5. If your handler functions contain Swagger annotations, these are extracted and added to the documentation.
6. The complete OpenAPI specification is returned as JSON.

## Example

See the `example.go` file in this package for a complete example of how to document your API handlers.
