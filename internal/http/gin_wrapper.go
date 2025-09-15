package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/http/middlewares"
	"github.com/Blocktunium/gonyx/internal/http/swagger"
	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/Blocktunium/gonyx/internal/logger"
	logTypes "github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// MARK: Variables

var (
	HttpServerMaintenanceType = logTypes.NewLogType("HTTP_SERVER_MAINTENANCE")
)

// Mark: Definitions

// GinServer struct
type GinServer struct {
	name                  string
	config                types.GinServerConfig
	app                   *http.Server
	baseRouter            *gin.Engine
	versionGroups         map[string]*gin.RouterGroup
	groups                map[string]*gin.RouterGroup
	supportedMiddlewares  []string
	defaultRequestMethods []string
	swaggerJSON           []byte // Store generated swagger JSON

	predefinedGroups []struct {
		name       string
		f          gin.HandlerFunc
		groupNames []string
	}

	predefinedRoutes []struct {
		method    string
		path      string
		f         []func(c *gin.Context)
		routeName string
		versions  []string
		groups    []string
	}
}

// init - Server Constructor - It initializes the server
func (s *GinServer) init(name string, serverConfig types.GinServerConfig, rawConfig map[string]interface{}) error {
	log.Println("New Http Server have been created...")
	s.name = name
	s.config = serverConfig

	if s.config.Config.RequestMethods != nil {
		if s.config.Config.RequestMethods[0] == "ALL" {
			s.defaultRequestMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
		} else {
			s.defaultRequestMethods = s.config.Config.RequestMethods
		}
	}

	// Set Gin mode based on environment
	ginMode := gin.ReleaseMode
	if env, err := config.GetManager().Get("app", "env"); err == nil {
		if env == "dev" {
			ginMode = gin.DebugMode
		}
	}
	gin.SetMode(ginMode)

	s.baseRouter = gin.New()

	s.groups = make(map[string]*gin.RouterGroup)
	s.supportedMiddlewares = []string{
		"logger",
		"favicon",
		"cors",
	}

	// get middleware objects and pass it to the attachMiddlewares function
	if v, ok := rawConfig["middlewares"].(map[string]interface{}); ok {
		s.attachMiddlewares(serverConfig.Middlewares.Order, v)
	}

	s.createVersionGroups(serverConfig.Versions)

	// if predefined before and just restarting
	if len(s.predefinedGroups) > 0 {
		for _, item := range s.predefinedGroups {
			s.AddGroup(item.name, item.f, item.groupNames...)
		}
	}

	if len(s.predefinedRoutes) > 0 {
		for _, item := range s.predefinedRoutes {
			if len(item.f) > 1 {
				s.AddRouteWithMultiHandlers(item.method, item.path, item.f, item.routeName, item.versions, item.groups)
			} else {
				s.AddRoute(item.method, item.path, item.f[0], item.routeName, item.versions, item.groups)
			}
		}
	}

	// Add Swagger documentation if enabled
	if s.config.Swagger.Enabled {
		// Generate Swagger JSON programmatically using swaggo before adding routes
		if err := s.generateSwaggerJSON(); err != nil {
			log.Printf("Warning: Failed to generate swagger JSON: %v", err)
		}
		s.addSwagger()
	}

	return nil
}

func (s *GinServer) createVersionGroups(versions []string) {
	s.versionGroups = make(map[string]*gin.RouterGroup)
	for _, item := range versions {
		s.versionGroups[item] = s.baseRouter.Group(item)
	}
}

func (s *GinServer) attachMiddlewares(orders []string, rawConfig map[string]interface{}) {
	for _, item := range orders {
		if utils.ArrayContains(&s.supportedMiddlewares, item) {
			switch item {
			case "logger":
				{
					// check which logger must be used
					loggerType, err := config.GetManager().Get("logger", "type")
					if err == nil {
						if loggerType == "zap" {
							s.baseRouter.Use(middlewares.ZapLogger())
							s.baseRouter.Use(middlewares.ZapRecoveryLogger())
						} else if loggerType == "logme" {
							s.baseRouter.Use(middlewares.LogMeLogger())
							s.baseRouter.Use(middlewares.LogMeRecoveryLogger())
						}
					}
				}
			case "favicon":
				if loggerCfg, ok := rawConfig[item].(map[string]interface{}); ok {
					jsonBody, err2 := json.Marshal(loggerCfg)
					if err2 == nil {
						var obj types.FaviconMiddlewareConfig
						err := json.Unmarshal(jsonBody, &obj)
						if err != nil {
							s.baseRouter.Use(middlewares.FaviconMiddleware(obj))
							break
						}
					}
				}
			case "cors":
				s.baseRouter.Use(func(c *gin.Context) {
					c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
					c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
					c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
					c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(s.defaultRequestMethods, ","))

					if c.Request.Method == "OPTIONS" {
						c.AbortWithStatus(204)
						return
					}

					c.Next()
				})

			}
		}
	}
}

func (s *GinServer) setupStatic() {
	// Add static file serving functionality here if needed
}

func (s *GinServer) addGroup(keyName string, groupName string, router *gin.RouterGroup, f gin.HandlerFunc) {
	if f == nil {
		s.groups[keyName] = router.Group(groupName)
	} else {
		s.groups[keyName] = router.Group(groupName, f)
	}
}

// getProjectRoot returns the project root directory dynamically
func getProjectRoot() string {
	// Try to get from environment variable first
	if projectRoot := os.Getenv("GONYX_PROJECT_ROOT"); projectRoot != "" {
		return projectRoot
	}

	// Try to get from config
	if projectRoot, err := config.GetManager().Get("base", "project_root"); err == nil {
		if rootPath, ok := projectRoot.(string); ok && rootPath != "" {
			return rootPath
		}
	}

	// Fallback: try to determine from current working directory
	if wd, err := os.Getwd(); err == nil {
		// Look for go.mod file to determine project root
		dir := wd
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break // reached root directory
			}
			dir = parent
		}
		// If go.mod not found, use current working directory
		return wd
	}

	// Final fallback: current directory
	return "."
}

// generateSwaggerJSON generates swagger JSON using swaggo library programmatically
func (s *GinServer) generateSwaggerJSON() error {
	// Get project root directory dynamically
	projectRoot := getProjectRoot()
	log.Printf("Using project root: %s", projectRoot)

	// Create swaggo generator
	swaggoGen := swagger.NewSwaggoGenerator(projectRoot)

	// Parse host and port for the generator
	host := "localhost"
	port := "3000"
	if s.config.ListenAddress != "" {
		parts := strings.Split(s.config.ListenAddress, ":")
		if len(parts) > 1 {
			if strings.HasPrefix(s.config.ListenAddress, ":") {
				host = "127.0.0.1"
			} else {
				host = parts[0]
			}
			port = parts[1]
		} else {
			if strings.HasPrefix(s.config.ListenAddress, ":") {
				host = "127.0.0.1"
				port = parts[0]
			} else {
				host = parts[0]
				port = "80"
			}
		}
	}

	// Set host and port in generator
	swaggoGen.SetHostAndPort(host, port)

	// First try to generate from swag registry (if annotations exist)
	routes := s.baseRouter.Routes()
	swaggerJSONBytes, err := swaggoGen.GenerateFromGinRoutes(routes)
	if err != nil {
		return fmt.Errorf("failed to generate swagger JSON from routes: %w", err)
	}

	// Store the generated JSON
	s.swaggerJSON = swaggerJSONBytes

	return nil
}

// addSwagger adds Swagger documentation endpoints to the server
func (s *GinServer) addSwagger() {
	// Parse host and port from the listen address
	host := "localhost"
	port := "3000"

	// Get host and port from listen address
	if s.config.ListenAddress != "" {
		parts := strings.Split(s.config.ListenAddress, ":")
		if len(parts) > 1 {
			// If address is like "localhost:8080" or "127.0.0.1:8080"
			if strings.HasPrefix(s.config.ListenAddress, ":") {
				host = "127.0.0.1"
			} else {
				host = parts[0]
			}
			port = parts[1]
		} else {
			if strings.HasPrefix(s.config.ListenAddress, ":") {
				// If address is just a port like ":8080"
				host = "127.0.0.1"
				port = parts[0]
			} else {
				// If address is just a port like ":8080"
				host = parts[0]
				port = "80"
			}
		}
	}

	// Create endpoint to serve the pre-generated swagger JSON
	s.baseRouter.GET("/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		if s.swaggerJSON != nil {
			// Update host and port in the swagger JSON dynamically
			var swaggerSpec map[string]interface{}
			if err := json.Unmarshal(s.swaggerJSON, &swaggerSpec); err == nil {
				swaggerSpec["host"] = fmt.Sprintf("%s:%s", host, port)
				if updatedJSON, err := json.Marshal(swaggerSpec); err == nil {
					c.Data(http.StatusOK, "application/json", updatedJSON)
					return
				}
			}
			// Fallback to original JSON if update fails
			c.Data(http.StatusOK, "application/json", s.swaggerJSON)
		} else {
			c.Data(http.StatusOK, "application/json", []byte("{}"))
		}
	})

	// Register Swagger UI handler with custom JSON URL
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://%s:%s/swagger.json", host, port))
	s.baseRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
}

// MARK: Public functions

// NewGinServer - create a new instance of Server and return it
func NewGinServer(name string, config types.GinServerConfig, rawConfig map[string]interface{}) (*GinServer, error) {
	server := &GinServer{}
	err := server.init(name, config, rawConfig)
	if err != nil {
		return nil, NewCreateServerErr(err)
	}
	return server, nil
}

func (s *GinServer) UpdateConfigs(config types.GinServerConfig, rawConfig map[string]interface{}) error {
	err := s.init(s.name, config, rawConfig)
	if err != nil {
		return NewUpdateServerConfigErr(err)
	}

	return nil
}

// Start - start the server and listen to provided address
func (s *GinServer) Start() error {
	s.app = &http.Server{
		Addr:         s.config.ListenAddress,
		Handler:      s.baseRouter,
		ReadTimeout:  s.config.Config.ReadTimeout,
		WriteTimeout: s.config.Config.WriteTimeout,
	}

	errCh := make(chan error)
	go func(ch chan error) {
		if err := s.app.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			close(ch)
		} else {
			ch <- err
			close(ch)
		}
	}(errCh)

	time.AfterFunc(3*time.Second, func() {
		errCh <- nil
	})
	err := <-errCh

	if err == nil {
		l, _ := logger.GetManager().GetLogger()
		if l != nil {
			l.Log(logTypes.NewLogObject(logTypes.INFO, "http.Server.Start", HttpServerMaintenanceType, time.Now(), "Starting the Http server ...", s.config.ListenAddress))
		}
	} else {
		l, _ := logger.GetManager().GetLogger()
		if l != nil {
			l.Log(logTypes.NewLogObject(logTypes.ERROR, "http.Server.Start", HttpServerMaintenanceType, time.Now(), "Starting the Http server failed ...", err))
		}
	}

	return err
}

// Stop - stop the server
func (s *GinServer) Stop() error {
	err := s.app.Shutdown(context.Background())
	if err != nil {
		return NewShutdownServerErr(err)
	}
	return nil
}

// AttachErrorHandler - attach a custom error handler to the server
func (s *GinServer) AttachErrorHandler(f func(ctx *gin.Context, err any)) {
	s.baseRouter.Use(gin.CustomRecovery(f))
}

func (s *GinServer) AddGroup(groupName string, f gin.HandlerFunc, groups ...string) error {
	s.predefinedGroups = append(s.predefinedGroups, struct {
		name       string
		f          gin.HandlerFunc
		groupNames []string
	}{name: groupName, groupNames: groups, f: f})

	if len(groups) > 0 {

	} else {
		for key, item := range s.versionGroups {
			newKey := fmt.Sprintf("%s.%s", key, groupName)
			s.addGroup(newKey, groupName, item, f)
		}

		s.addGroup(groupName, groupName, &s.baseRouter.RouterGroup, f)
	}

	return nil
}

func (s *GinServer) AddRoute(method string, path string, f func(c *gin.Context), routeName string, versions []string, groups []string) error {
	s.predefinedRoutes = append(s.predefinedRoutes, struct {
		method    string
		path      string
		f         []func(c *gin.Context)
		routeName string
		versions  []string
		groups    []string
	}{method: method, path: path, f: []func(c *gin.Context){f}, routeName: routeName, versions: versions, groups: groups})

	if utils.ArrayContains(&s.defaultRequestMethods, method) {
		groupsExist := false
		if groups != nil {
			if len(groups) > 0 {
				groupsExist = true
			}
		}

		versionsExist := false
		if versions != nil {
			if len(versions) > 0 {
				versionsExist = true
			}
		}

		if groupsExist {
			for _, g := range groups {
				if versionsExist {
					for _, v := range versions {
						if v == "all" {
							for k := range s.versionGroups {
								newKey := fmt.Sprintf("%s.%s", k, g)
								if router, ok := s.groups[newKey]; ok {
									router.Handle(method, path, f)
								}
							}
							break
						} else if v == "" {
							if router, ok := s.groups[g]; ok {
								router.Handle(method, path, f)
							}
							break
						} else {
							newKey := fmt.Sprintf("%s.%s", v, g)
							if router, ok := s.groups[newKey]; ok {
								router.Handle(method, path, f)
							}
						}
					}
				} else {
					if savedGroup, ok := s.groups[g]; ok {
						savedGroup.Handle(method, path, f)
					}
				}
			}
		} else {
			if versionsExist {
				for _, v := range versions {
					if router, ok := s.versionGroups[v]; ok {
						router.Handle(method, path, f)
					} else {
						if v == "all" {
							for _, router1 := range s.versionGroups {
								router1.Handle(method, path, f)
							}
							break
						} else if v == "" {
							s.baseRouter.Handle(method, path, f)
							break
						}
					}
				}
			} else {
				s.baseRouter.Handle(method, path, f)
			}
		}
		return nil
	}

	return NewNotSupportedHttpMethodErr(method)
}

// swaggerGenerator is a singleton instance of the Swagger generator
var swaggerGenerator *swagger.Generator

// getSwaggerGenerator initializes or returns the Swagger generator
func getSwaggerGenerator() *swagger.Generator {
	if swaggerGenerator == nil {
		swaggerGenerator = swagger.NewGenerator()
	}
	return swaggerGenerator
}

// AddRouteWithMultiHandlers - add a route to the server
func (s *GinServer) AddRouteWithMultiHandlers(method string, path string, f []func(c *gin.Context), routeName string, versions []string, groups []string) error {
	s.predefinedRoutes = append(s.predefinedRoutes, struct {
		method    string
		path      string
		f         []func(c *gin.Context)
		routeName string
		versions  []string
		groups    []string
	}{method: method, path: path, f: f, routeName: routeName, versions: versions, groups: groups})

	// check that whether is acceptable to add this route method
	if utils.ArrayContains(&s.defaultRequestMethods, method) {
		if len(groups) > 0 {
			for _, g := range groups {
				if len(versions) > 0 {
					for _, v := range versions {
						if v == "all" {
							for k := range s.versionGroups {
								newKey := fmt.Sprintf("%s.%s", k, g)
								if router, ok := s.groups[newKey]; ok {
									var p []gin.HandlerFunc
									for _, item := range f {
										p = append(p, item)
									}
									router.Handle(method, path, p...)
								}
							}
							break
						} else if v == "" {
							if router, ok := s.groups[g]; ok {
								var p []gin.HandlerFunc
								for _, item := range f {
									p = append(p, item)
								}
								router.Handle(method, path, p...)
							}
							break
						} else {
							newKey := fmt.Sprintf("%s.%s", v, g)
							if router, ok := s.groups[newKey]; ok {
								var p []gin.HandlerFunc
								for _, item := range f {
									p = append(p, item)
								}

								router.Handle(method, path, p...)
							}
						}
					}
				} else {
					if savedGroup, ok := s.groups[g]; ok {
						var p []gin.HandlerFunc
						for _, item := range f {
							p = append(p, item)
						}

						savedGroup.Handle(method, path, p...)
					}
				}
			}
		} else {
			if len(versions) > 0 {
				for _, v := range versions {
					if router, ok := s.versionGroups[v]; ok {
						var p []gin.HandlerFunc
						for _, item := range f {
							p = append(p, item)
						}
						router.Handle(method, path, p...)
					} else {
						if v == "all" {
							for _, router1 := range s.versionGroups {
								var p []gin.HandlerFunc
								for _, item := range f {
									p = append(p, item)
								}
								router1.Handle(method, path, p...)
							}
							break
						} else if v == "" {
							var p []gin.HandlerFunc
							for _, item := range f {
								p = append(p, item)
							}

							s.baseRouter.Handle(method, path, p...)
						}
					}
				}
			} else {
				var p []gin.HandlerFunc
				for _, item := range f {
					p = append(p, item)
				}
				s.baseRouter.Handle(method, path, p...)
			}
		}
		return nil
	}

	return NewNotSupportedHttpMethodErr(method)
}

// GetAllRoutes - Get all Routes
func (s *GinServer) GetAllRoutes() gin.RoutesInfo {
	return s.baseRouter.Routes()
}
