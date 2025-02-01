package http

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/http"
	"github.com/gin-gonic/gin"
)

// Http Methods
const (
	MethodGet     = "GET"     // RFC 7231, 4.3.1
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    = "POST"    // RFC 7231, 4.3.3
	MethodPut     = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   = "PATCH"   // RFC 5789
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8
	methodUse     = "USE"
)

// HttpRoute - Structure of the route
type HttpRoute struct {
	Method     string
	Path       string
	RouteName  string
	Versions   []string
	GroupNames []string
	F          func(c *gin.Context)
	Servers    []string
}

// HttpGroup - Structure of the group
type HttpGroup struct {
	GroupName string
	F         func(c *gin.Context)
	Groups    []string
	Servers   []string
}

// AddHttpRouteByObj - add route by HttpRoute obj
func AddHttpRouteByObj(httpRoute HttpRoute) error {
	return http.GetManager().AddRoute(httpRoute.Method,
		httpRoute.Path,
		httpRoute.F,
		httpRoute.RouteName,
		httpRoute.Versions,
		httpRoute.GroupNames,
		httpRoute.Servers...)
}

// AddHttpRoute - Add route by parameters
func AddHttpRoute(method string, path string, f func(c *gin.Context), routeName string, versions []string, groupNames []string, serverName ...string) error {
	return http.GetManager().AddRoute(method,
		path,
		f,
		routeName,
		versions,
		groupNames,
		serverName...)
}

// AddBulkHttpRoutes - add bulk http routes to the server
func AddBulkHttpRoutes(httpRoutes []HttpRoute) error {
	for _, httpRoute := range httpRoutes {
		err := http.GetManager().AddRoute(httpRoute.Method,
			httpRoute.Path,
			httpRoute.F,
			httpRoute.RouteName,
			httpRoute.Versions,
			httpRoute.GroupNames,
			httpRoute.Servers...)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRouteByName - Get route by providing the route name from specific server
//func GetRouteByName(routeName string, serverName ...string) (*fiber.Route, error) {
//	return http.GetManager().GetRouteByName(routeName, serverName...)
//}

// AddHttpGroupByObj - add group by HttpGroup obj
func AddHttpGroupByObj(group HttpGroup) error {
	return http.GetManager().AddGroup(
		group.GroupName,
		group.F,
		group.Groups,
		group.Servers...,
	)
}

// AddHttpGroup - Add group by parameters
func AddHttpGroup(groupName string, f func(c *gin.Context), groups []string, serverName ...string) error {
	return http.GetManager().AddGroup(groupName,
		f, groups, serverName...)
}

// AddBulkHttpGroups - add bulk http groups to the server
func AddBulkHttpGroups(httpGroups []HttpGroup) error {
	for _, group := range httpGroups {
		err := http.GetManager().AddGroup(
			group.GroupName,
			group.F,
			group.Groups,
			group.Servers...,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// AttachHttpErrorHandler - Attach http error handler to the manager
func AttachHttpErrorHandler(f func(ctx *gin.Context, err any), serverNames ...string) error {
	return http.GetManager().AttachErrorHandler(f, serverNames...)
}

// PrintAllRoutes - Print all routes on the screen
func PrintAllRoutes() {
	routes := http.GetManager().GetAllRoutes()
	for _, item := range routes {
		fmt.Println(item)
	}
}
