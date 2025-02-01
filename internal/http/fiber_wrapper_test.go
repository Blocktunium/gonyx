package http

import (
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/http/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestFiberWrapper_startServer(t *testing.T) {
	// create a new fiber wrapper and start it

	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	err = server.Start()
	if err != nil {
		t.Errorf("Starting HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}
}

func TestFiberWrapper_AddRoute(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, "TestGet", []string{}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) == 0 {
		t.Errorf("Expected 1 routes added to server, but got %v", len(routes))
		return
	}

	// check the name of the route
	r := server.app.GetRoute("TestGet")
	if r.Method != fiber.MethodGet {
		t.Errorf("Expected to get method %v, but got %v", fiber.MethodGet, r.Method)
		return
	}
	if r.Path != "/" {
		t.Errorf("Expected to get route '%v', but got '%v'", "/", r.Path)
		return
	}
}

func TestFiberWrapper_AddMultipleRoute(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	for i := range []int{1, 2, 3} {
		server.AddRoute(fiber.MethodGet, fmt.Sprintf("/%v", i), func(c *fiber.Ctx) error {
			return nil
		}, fmt.Sprintf("Test%v", i), []string{}, []string{})
	}

	routes := server.app.GetRoutes()
	if len(routes) == 0 {
		t.Errorf("Expected 1 routes added to server, but got %v", len(routes))
		return
	} else if len(routes) != 3 {
		t.Errorf("Expected 1 routes added to server, but got %v", len(routes))
		return
	}

	// check the name of the route
	r := server.app.GetRoute("Test2")
	if r.Method != fiber.MethodGet {
		t.Errorf("Expected to get method %v, but got %v", fiber.MethodGet, r.Method)
		return
	}
	if r.Path != "/2" {
		t.Errorf("Expected to get route '%v', but got '%v'", "/2", r.Path)
		return
	}
}

func TestFiberWrapper_AddBlockRouteMethod(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	server.AddRoute("TTT", "/", func(c *fiber.Ctx) error {
		return nil
	}, "TestGet", []string{}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) > 0 {
		t.Errorf("Expected the route not to be added, but added")
		return
	}
}

func TestFiberWrapper_GetRouteByName(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) == 0 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 1, len(routes))
		return
	}

	route, err := server.GetRouteByName(routeName)
	routes = server.app.GetRoutes()
	if err != nil {
		t.Errorf("Get route by name --> Expected no err:  but got %v", err)
		return
	}

	if route.Name != routeName {
		t.Errorf("Expected to get route by name: %v, but got %v", routeName, route)
	}
}

func TestFiberWrapper_AddRouteToOneSupportedVersion(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{"v1"}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) == 0 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 1, len(routes))
		return
	}

	expectedPath := "/v1"
	if routes[0].Path != expectedPath {
		t.Errorf("Expected that get route '%v', but got %v", expectedPath, routes[0].Path)
		return
	}
}

func TestFiberWrapper_AddRouteToTwoSupportedVersions(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{"v1", "v2"}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) != 2 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 2, len(routes))
		return
	}

	expectedPath := []string{"/v1/", "/v2/"}
	var actualPath []string
	for _, r := range routes {
		actualPath = append(actualPath, r.Path)
	}
	if !reflect.DeepEqual(actualPath, expectedPath) {
		t.Errorf("Expected that get route '%v', but got %v", expectedPath, actualPath)
		return
	}
}

func TestFiberWrapper_AddRouteToAllVersions(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{"all"}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) != 2 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 2, len(routes))
		return
	}

	expectedPath := []string{"/v1/", "/v2/"}
	var actualPath []string
	for _, r := range routes {
		actualPath = append(actualPath, r.Path)
	}
	if !reflect.DeepEqual(actualPath, expectedPath) {
		t.Errorf("Expected that get route '%v', but got %v", expectedPath, actualPath)
		return
	}
}

func TestFiberWrapper_AddRouteToNoVersions(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{""}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) != 1 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 1, len(routes))
		return
	}

	expectedPath := []string{"/"}
	var actualPath []string
	for _, r := range routes {
		actualPath = append(actualPath, r.Path)
	}
	if !reflect.DeepEqual(actualPath, expectedPath) {
		t.Errorf("Expected that get route '%v', but got %v", expectedPath, actualPath)
		return
	}
}

func TestFiberWrapper_AddGroup(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	groupName := "p"
	server.AddGroup(groupName, func(c *fiber.Ctx) error {
		return nil
	})

	expectedGroups := []string{"v1.p", "v2.p", "p"}
	var actualGroups []string
	for key, _ := range server.groups {
		actualGroups = append(actualGroups, key)
	}

	if !reflect.DeepEqual(expectedGroups, actualGroups) {
		t.Errorf("Expected that get group list '%v', but got %v", expectedGroups, actualGroups)
		return
	}
}

func TestFiberWrapper_AddGroupToGroup(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	groupName := "p2"
	server.AddGroup(groupName, func(c *fiber.Ctx) error {
		return nil
	})

	groupName = "p"
	server.AddGroup(groupName, func(c *fiber.Ctx) error {
		return nil
	}, "p2")

	expectedGroups := []string{"p2.p", "v1.p2", "v2.p2", "p2", "v1.p2.p", "v2.p2.p"}
	var actualGroups []string
	for key, _ := range server.groups {
		actualGroups = append(actualGroups, key)
	}

	if !reflect.DeepEqual(expectedGroups, actualGroups) {
		t.Errorf("Expected that get group list '%v', but got %v", expectedGroups, actualGroups)
		return
	}
}

func TestFiberWrapper_AddRouteToSpecificGroup(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000", Versions: []string{"v1", "v2"}}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	groupName := "p2"
	server.AddGroup(groupName, nil)

	routeName := "TestGet"
	server.AddRoute(fiber.MethodGet, "/ttt", func(c *fiber.Ctx) error {
		return nil
	}, routeName, []string{"v1"}, []string{groupName})

	routes := server.app.GetRoutes()
	if len(routes) != 1 {
		t.Errorf("Expected that %v route(s) exist, but got %v", 1, len(routes))
		return
	}
	expectedRoute := "/v1/p2/ttt"
	actualRoute := routes[0].Path
	if !reflect.DeepEqual(actualRoute, expectedRoute) {
		t.Errorf("Expected that the first route be %v, but got %v", expectedRoute, actualRoute)
		return
	}
}

func TestFiberWrapper_RequestMethodsNotAllowed(t *testing.T) {
	serverConfig := types.ServerConfig{
		ListenAddress: ":3000",
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
			RequestMethods:       []string{"GET"},
		},
	}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	err = server.AddRoute(fiber.MethodPost, "/p", func(c *fiber.Ctx) error {
		return nil
	}, "TestPOST", []string{}, []string{})
	if err == nil {
		t.Errorf("Expected to not add POST method, but it's added")
	}
}

func TestFiberWrapper_RequestMethodsFilter(t *testing.T) {
	serverConfig := types.ServerConfig{
		ListenAddress: ":3000",
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
			RequestMethods:       []string{"GET"},
		},
	}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	server.AddRoute(fiber.MethodPost, "/p", func(c *fiber.Ctx) error {
		return nil
	}, "TestPOST", []string{}, []string{})

	server.AddRoute(fiber.MethodGet, "/g", func(c *fiber.Ctx) error {
		return nil
	}, "TestGET", []string{}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) != 1 {
		t.Errorf("Expected 1 route added to server, but got %v", len(routes))
		return
	}

	// check the name of the route
	r := server.app.GetRoute("TestPOST")
	if r.Path == "/p" {
		t.Errorf("Expected to get route '%v', but got '%v'", "/", r.Path)
		return
	}
}

func TestFiberWrapper_StaticHandling(t *testing.T) {
	serverConfig := types.ServerConfig{
		ListenAddress: ":3000",
		SupportStatic: true,
		Static: struct {
			Prefix string       `json:"prefix"`
			Root   string       `json:"root"`
			Config fiber.Static `json:"config"`
		}{
			Prefix: "/public",
			Root:   "../../data/static/",
			Config: fiber.Static{Compress: false, ByteRange: false, Browse: false, Download: false, Index: "index.html", CacheDuration: 10 * time.Second, MaxAge: 0, ModifyResponse: nil, Next: nil},
		},
	}

	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}
	go server.Start()
	defer server.Stop()

	req := httptest.NewRequest(fiber.MethodGet, "http://127.0.0.1:3000/public/index.html", nil)
	resp, err := server.app.Test(req)
	utils.AssertEqual(t, nil, err, "Expected to not get error")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")

	body, err := io.ReadAll(resp.Body)
	utils.AssertEqual(t, nil, err, "Response Body Result")

	expectedText := `<!DOCTYPE html>
<html lang="en">
<body>
    Hi
</body>
</html>`
	utils.AssertEqual(t, expectedText, string(body), "Unmatched content")
}

func TestFiberWrapper_AddRouteWithMultipleHandler(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	server.AddRouteWithMultiHandlers(fiber.MethodGet, "/", []func(c *fiber.Ctx) error{
		func(c *fiber.Ctx) error { return nil },
		func(c *fiber.Ctx) error { return errors.New("error") },
	}, "TestGet", []string{}, []string{})

	routes := server.app.GetRoutes()
	if len(routes) == 0 {
		t.Errorf("Expected 1 routes added to server, but got %v", len(routes))
		return
	}

	// check the name of the route
	r := server.app.GetRoute("TestGet")
	if r.Method != fiber.MethodGet {
		t.Errorf("Expected to get method %v, but got %v", fiber.MethodGet, r.Method)
		return
	}
	if r.Path != "/" {
		t.Errorf("Expected to get route '%v', but got '%v'", "/", r.Path)
		return
	}
}

func TestFiberWrapper_CustomErrorHandler(t *testing.T) {
	serverConfig := types.ServerConfig{ListenAddress: ":3000"}
	server, err := NewServer("http", serverConfig)
	if err != nil {
		t.Errorf("Creating HTTP Server --> Expected: %v, but got %v", nil, err)
		return
	}

	server.AttachErrorHandler(func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError

		// Retrieve the custom status code if it's a *fiber.Error
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		utils.AssertEqual(t, fiber.StatusNotFound, code, "Status code")
		ctx.Status(code).SendString(e.Error())
		return nil
	})

	go server.Start()
	defer server.Stop()

	req := httptest.NewRequest(fiber.MethodGet, "http://127.0.0.1:3000/t", nil)
	resp, err := server.app.Test(req)
	utils.AssertEqual(t, nil, err, "Expected to not get error")
	utils.AssertEqual(t, 404, resp.StatusCode, "Status code")
}
