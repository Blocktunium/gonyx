package http

import (
	"encoding/json"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"strings"
	"time"
)

// Mark: Definitions

// Server struct
type FiberServer struct {
	name                  string
	config                types.ServerConfig
	app                   *fiber.App
	versionGroups         map[string]fiber.Router
	groups                map[string]fiber.Router
	supportedMiddlewares  []string
	defaultRequestMethods []string

	predefinedGroups []struct {
		name       string
		f          func(c *fiber.Ctx) error
		groupNames []string
	}

	predefinedRoutes []struct {
		method    string
		path      string
		f         []func(c *fiber.Ctx) error
		routeName string
		versions  []string
		groups    []string
	}
}

// init - Server Constructor - It initializes the server
func (s *FiberServer) init(name string, serverConfig types.ServerConfig, rawConfig map[string]interface{}) error {
	s.name = name
	s.config = serverConfig

	// Get application name from the config manager
	appName := config.GetManager().GetName()
	requestMethods := fiber.DefaultMethods
	if s.config.Config.RequestMethods != nil {
		if s.config.Config.RequestMethods[0] != "ALL" {
			requestMethods = s.config.Config.RequestMethods
		}
	}
	s.defaultRequestMethods = requestMethods

	s.app = fiber.New(fiber.Config{
		Prefork:        false,
		ServerHeader:   "Gonyx",
		AppName:        appName,
		RequestMethods: requestMethods,
	})

	if s.config.SupportStatic == true {
		s.setupStatic()
	}

	s.groups = make(map[string]fiber.Router)
	s.supportedMiddlewares = []string{
		"logger",
		"favicon",
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

func (s *FiberServer) createVersionGroups(versions []string) {
	s.versionGroups = make(map[string]fiber.Router)
	for _, item := range versions {
		s.versionGroups[item] = s.app.Group(item)
	}
}

func (s *FiberServer) attachMiddlewares(orders []string, rawConfig map[string]interface{}) {
	for _, item := range orders {
		if utils.ArrayContains(&s.supportedMiddlewares, item) {
			switch item {
			case "logger":
				//key := fmt.Sprintf("middlewares.%s", item)
				// read config
				if loggerCfg, ok := rawConfig[item].(map[string]interface{}); ok {
					jsonBody, err2 := json.Marshal(loggerCfg)
					if err2 == nil {
						var obj types.LoggerMiddlewareConfig
						err := json.Unmarshal(jsonBody, &obj)
						if err == nil {
							// Everything is ok and let's go define logger config
							loggerMiddlewareCfg := logger.Config{
								Next:         nil,
								Done:         nil,
								Format:       obj.Format,
								TimeFormat:   obj.TimeFormat,
								TimeZone:     obj.TimeZone,
								TimeInterval: time.Duration(obj.TimeInterval) * time.Millisecond,
							}
							if obj.Output == "stdout" {
								loggerMiddlewareCfg.Output = os.Stdout
							}
							s.app.Use(logger.New(loggerMiddlewareCfg))
							break
						}
					}
				} else {
					s.app.Use(logger.New())
					break
				}
			case "favicon":
				//key := fmt.Sprintf("middlewares.%s", item)
				// read config
				if loggerCfg, ok := rawConfig[item].(map[string]interface{}); ok {
					jsonBody, err2 := json.Marshal(loggerCfg)
					if err2 == nil {
						var obj types.FaviconMiddlewareConfig
						err := json.Unmarshal(jsonBody, &obj)
						if err != nil {
							faviconMiddlewareCfg := favicon.Config{
								File:         obj.File,
								URL:          obj.URL,
								CacheControl: obj.CacheControl,
							}
							s.app.Use(favicon.New(faviconMiddlewareCfg))
							break
						}
					}
				} else {
					s.app.Use(favicon.New())
					break
				}
			}
		}
	}

	fmt.Println(s.app)
}

func (s *FiberServer) addGroup(keyName string, groupName string, router fiber.Router, f func(c *fiber.Ctx) error) {
	if f == nil {
		s.groups[keyName] = router.Group(groupName)
	} else {
		s.groups[keyName] = router.Group(groupName, f)
	}
}

func (s *FiberServer) setupStatic() {
	s.app.Static(s.config.Static.Prefix, s.config.Static.Root, s.config.Static.Config)
}

// MARK: Public functions

// NewServer - create a new instance of Server and return it
func NewServer(name string, config types.ServerConfig, rawConfig map[string]interface{}) (*FiberServer, error) {
	server := &FiberServer{}
	err := server.init(name, config, rawConfig)
	if err != nil {
		return nil, NewCreateServerErr(err)
	}
	return server, nil
}

func (s *FiberServer) UpdateConfigs(config types.ServerConfig, rawConfig map[string]interface{}) error {
	err := s.init(s.name, config, rawConfig)
	if err != nil {
		return NewUpdateServerConfigErr(err)
	}

	return nil
}

// Start - start the server and listen to provided address
func (s *FiberServer) Start() error {
	err := s.app.Listen(s.config.ListenAddress)
	if err != nil {
		return NewStartServerErr(s.config.ListenAddress, err)
	}
	return nil
}

// Stop - stop the server
func (s *FiberServer) Stop() error {
	err := s.app.Shutdown()
	if err != nil {
		return NewShutdownServerErr(err)
	}
	return nil
}

// AttachErrorHandler - attach a custom error handler to the server
func (s *FiberServer) AttachErrorHandler(f func(ctx *fiber.Ctx, err error) error) {
	oldConfig := s.app.Config()
	oldConfig.ErrorHandler = f
	s.app = fiber.New(oldConfig)
}

// AddRoute - add a route to the server
func (s *FiberServer) AddRoute(method string, path string, f func(c *fiber.Ctx) error, routeName string, versions []string, groups []string) error {
	s.predefinedRoutes = append(s.predefinedRoutes, struct {
		method    string
		path      string
		f         []func(c *fiber.Ctx) error
		routeName string
		versions  []string
		groups    []string
	}{method: method, path: path, f: []func(c *fiber.Ctx) error{f}, routeName: routeName, versions: versions, groups: groups})

	// check that whether is acceptable to add this route method
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
									router.Add(method, path, f)
									if strings.TrimSpace(routeName) != "" {
										//router.Name(routeName)
										s.app.Name(routeName)
									}
								}
							}
							break
						} else if v == "" {
							if router, ok := s.groups[g]; ok {
								router.Add(method, path, f)
								if strings.TrimSpace(routeName) != "" {
									//router.Name(routeName)
									s.app.Name(routeName)
								}
							}
							break
						} else {
							newKey := fmt.Sprintf("%s.%s", v, g)
							if router, ok := s.groups[newKey]; ok {
								router.Add(method, path, f)
								if strings.TrimSpace(routeName) != "" {
									s.app.Name(routeName)
									//router.Name(routeName)
								}
							}
						}
					}
				} else {
					if savedGroup, ok := s.groups[g]; ok {
						savedGroup.Add(method, path, f)
						if strings.TrimSpace(routeName) != "" {
							//savedGroup.Name(routeName)
							s.app.Name(routeName)

						}
					}
				}
			}
		} else {
			if versionsExist {
				for _, v := range versions {
					if router, ok := s.versionGroups[v]; ok {
						router.Add(method, path, f)
						if strings.TrimSpace(routeName) != "" {
							//router.Name(routeName)
							s.app.Name(routeName)
						}
					} else {
						if v == "all" {
							for _, router1 := range s.versionGroups {
								router1.Add(method, path, f)
								if strings.TrimSpace(routeName) != "" {
									//router1.Name(routeName)
									s.app.Name(routeName)
								}
							}
							break
						} else if v == "" {
							s.app.Add(method, path, f)
							if strings.TrimSpace(routeName) != "" {
								s.app.Name(routeName)
							}
							break
						}
					}
				}
			} else {
				s.app.Add(method, path, f)
				if strings.TrimSpace(routeName) != "" {
					s.app.Name(routeName)
				}
			}
		}
		return nil
	}

	return NewNotSupportedHttpMethodErr(method)
}

// AddGroup - add a group to the server
func (s *FiberServer) AddGroup(groupName string, f func(c *fiber.Ctx) error, groups ...string) error {
	s.predefinedGroups = append(s.predefinedGroups, struct {
		name       string
		f          func(c *fiber.Ctx) error
		groupNames []string
	}{name: groupName, f: f, groupNames: groups})
	if len(groups) > 0 {
		for _, g := range groups {
			for key := range s.versionGroups {
				gKey := fmt.Sprintf("%s.%s", key, g)
				if r, ok := s.groups[gKey]; ok {
					newKey := fmt.Sprintf("%s.%s.%s", key, g, groupName)
					s.addGroup(newKey, groupName, r, f)
				} else {
					return NewGroupRouteNotExistErr(gKey)
				}
			}

			newKey := fmt.Sprintf("%s.%s", g, groupName)
			s.addGroup(newKey, groupName, s.app, f)
		}
	} else {
		for key, item := range s.versionGroups {
			newKey := fmt.Sprintf("%s.%s", key, groupName)
			s.addGroup(newKey, groupName, item, f)
		}

		s.addGroup(groupName, groupName, s.app, f)
	}

	return nil
}

// AddRouteWithMultiHandlers - add a route to the server
func (s *FiberServer) AddRouteWithMultiHandlers(method string, path string, f []func(c *fiber.Ctx) error, routeName string, versions []string, groups []string) error {
	s.predefinedRoutes = append(s.predefinedRoutes, struct {
		method    string
		path      string
		f         []func(c *fiber.Ctx) error
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
									router.Add(method, path, f...)
									if strings.TrimSpace(routeName) != "" {
										//router.Name(routeName)
										s.app.Name(routeName)
									}
								}
							}
							break
						} else if v == "" {
							if router, ok := s.groups[g]; ok {
								router.Add(method, path, f...)
								if strings.TrimSpace(routeName) != "" {
									//router.Name(routeName)
									s.app.Name(routeName)
								}
							}
							break
						} else {
							newKey := fmt.Sprintf("%s.%s", v, g)
							if router, ok := s.groups[newKey]; ok {
								router.Add(method, path, f...)
								if strings.TrimSpace(routeName) != "" {
									s.app.Name(routeName)
									//router.Name(routeName)
								}
							}
						}
					}
				} else {
					if savedGroup, ok := s.groups[g]; ok {
						savedGroup.Add(method, path, f...)
						if strings.TrimSpace(routeName) != "" {
							//savedGroup.Name(routeName)
							s.app.Name(routeName)

						}
					}
				}
			}
		} else {
			if len(versions) > 0 {
				for _, v := range versions {
					if router, ok := s.versionGroups[v]; ok {
						router.Add(method, path, f...)
						if strings.TrimSpace(routeName) != "" {
							//router.Name(routeName)
							s.app.Name(routeName)
						}
					} else {
						if v == "all" {
							for _, router1 := range s.versionGroups {
								router1.Add(method, path, f...)
								if strings.TrimSpace(routeName) != "" {
									//router1.Name(routeName)
									s.app.Name(routeName)
								}
							}
							break
						} else if v == "" {
							s.app.Add(method, path, f...)
							if strings.TrimSpace(routeName) != "" {
								s.app.Name(routeName)
							}
							break
						}
					}
				}
			} else {
				s.app.Add(method, path, f...)
				if strings.TrimSpace(routeName) != "" {
					s.app.Name(routeName)
				}
			}
		}
		return nil
	}

	return NewNotSupportedHttpMethodErr(method)
}

// GetRouteByName _ get route by its name
func (s *FiberServer) GetRouteByName(name string) (*fiber.Route, error) {
	route := s.app.GetRoute(name)
	if route.Name != name {
		return nil, NewGetRouteByNameErr(name)
	}
	return &route, nil
}

// GetAllRoutes - Get all Routes
func (s *FiberServer) GetAllRoutes() []fiber.Route {
	return s.app.GetRoutes(true)
}
