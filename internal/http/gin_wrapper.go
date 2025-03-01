package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/http/middlewares"
	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/Blocktunium/gonyx/internal/logger"
	logTypes "github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
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

}

func (s *GinServer) addGroup(keyName string, groupName string, router *gin.RouterGroup, f gin.HandlerFunc) {
	if f == nil {
		s.groups[keyName] = router.Group(groupName)
	} else {
		s.groups[keyName] = router.Group(groupName, f)
	}
}

//func (s *GinServer) addSwagger() {
//	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
//	s.baseRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
//}

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
