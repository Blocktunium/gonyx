package commands

const (
	DefaultProjectDirectory = "."

	InitializeMessage         = `Zhycan > Create project skeleton ...`
	RootDirectoryIsCreated    = `Zhycan > Project root (%s) is created ...`
	RootDirectoryIsNotCreated = `Zhycan > Project root (%s) is not created correctly ...%v`
	GoModuleFileIsCreated     = `Zhycan > Go Module File "go.mod" is created ...`
	GoModuleFileIsNotCreated  = `Zhycan > Go Module File "go.mod" is not created ... %v`
	GoModuleIsCreated         = `Zhycan > Go Module "go.mod" is filled ...`
	GoModuleIsNotCreated      = `Zhycan > Go Module "go.mod" is not filled ... %v`
	MainGoFileIsCreated       = `Zhycan > Main program File "main.go" is created ...`
	MainGoFileIsNotCreated    = `Zhycan > Main program File "main.go" is not created ... %v`
	MainGoIsCreated           = `Zhycan > Main program "main.go" is filled ...`
	MainGoIsNotCreated        = `Zhycan > Main program "main.go" is not filled ... %v`
	UserExisted               = "Zhycan > User existed ..."
	UserNotExisted            = "Zhycan > User not existed ... %v"
	SubDirectoryIsNotCreated  = `Zhycan > Sub directory "%s" cannot be created ... %v`
	SubDirectoryIsCreated     = `Zhycan > Sub directory "%s" is created ...`
	AppControllerIsNotCreated = `Zhycan > App "controller.go" cannot be created ... %v`
	AppControllerIsCreated    = `Zhycan > App "controller.go" is created ...`
	AppModelIsNotCreated      = `Zhycan > App "model.go" cannot be created ... %v`
	AppModelIsCreated         = `Zhycan > App "model.go" is created ...`
	AppEngineIsNotCreated     = `Zhycan > App "app.go" cannot be created ... %v`
	AppEngineIsCreated        = `Zhycan > App "app.go" is created ...`

	RootCommandGoFileIsCreated    = `Zhycan > Root command File "commands/root.go" is created ...`
	RootCommandGoFileIsNotCreated = `Zhycan > Root command File "commands/root.go" is not created ... %v`
	RootCommandGoIsNotCreated     = `Zhycan > Root command "commands/root.go" is not filled ... %v`

	GitIgnoreFileIsCreated    = `Zhycan > Git Ignore File ".gitignore" is created ...`
	GitIgnoreFileIsNotCreated = `Zhycan > Git Ignore File ".gitignore" is not created ... %v`
	GitIgnoreIsNotCreated     = `Zhycan > Git Ignore ".gitignore" is not filled ... %v`

	GitInitExecutedError = `Zhycan > Cannot execute git init command ... %v`
	GitInitExecuted      = `Zhycan > Git repository is initialized ...`

	ConfigFileIsCreated       = `Zhycan > Config File "%s" is created ...`
	ConfigFileIsNotCreated    = `Zhycan > Config File "%s" is not created ... %v`
	ConfigDevFileIsCreated    = `Zhycan > Config File "%s" is created for "dev" mode ...`
	ConfigDevFileIsNotCreated = `Zhycan > Config File "%s" is not created for "dev" mode ... %v`

	GoModTidyExecutedError = `Zhycan > Cannot execute go mod tidy command ... %v`
	GoModTidyExecuted      = `Zhycan > "go mod tidy" command is executed ...`

	GreeterProtobufIsNotCreated = `Zhycan > App "greeter.proto" cannot be created ... %v`
	GreeterProtobufIsCreated    = `Zhycan > App "greeter.proto" is created ...`

	GreeterCompiledProtobufIsNotCreated = `Zhycan > App "greeter.pb.go" cannot be created ... %v`
	GreeterCompiledProtobufIsCreated    = `Zhycan > App "greeter.pb.go" is created ...`

	GreeterCompiledGrpcProtobufIsNotCreated = `Zhycan > App "greeter_grpc.pb.go" cannot be created ... %v`
	GreeterCompiledGrpcProtobufIsCreated    = `Zhycan > App "greeter_grpc.pb.go" is created ...`
)

const (
	goModTmpl = `module {{.ProjectName}}

go {{.Version}}

`

	mainTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "main.go" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/
package main

import (
    "fmt"
    "{{.ProjectName}}/commands"
    "github.com/abolfazlbeh/zhycan/pkg/config"
    "github.com/abolfazlbeh/zhycan/pkg/logger"
    "time"
)

/*
The main file of the project
*/
func main() {
    fmt.Println("{{.ProjectName}} is Started ...")

    // config module attributes
    baseConfigPath := "."           // the base path for the parameters. by default it's the current directory
    initialConfigMode := "dev"            // it can be override by environment value --> the value can be "dev" and "prod" and whatever you want
    configPrefix := "{{.ProjectName}}"    // this will be used in reading value from environment with this prefix

    err := config.InitializeManager(baseConfigPath, initialConfigMode, configPrefix)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Testing the logger module works properly
    logger.Log(logger.NewLogObject(
        logger.INFO, "main.go", logger.FuncMaintenanceType, time.Now().UTC(), "Logger Module Works Like A Charm ...", nil))

    // Execute the provided command
    commands.Execute()
}
`

	rootCommandTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "root.go" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/
package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/abolfazlbeh/zhycan/pkg/cli"
	"os"
	"{{.ProjectName}}/app"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{.ProjectName}}",
	Short: "A brief description of your application",
	Long: "A longer description that spans multiple lines and likely contains\nexamples and usage of using your application. For example:\nCobra is a CLI library for Go that empowers applications.\nThis application is a tool to generate the needed files\nto quickly create a Cobra application.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use == "runserver" {
			app1 := &app.App{}
			app1.Init()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Use)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test-corba-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

    // Attach Default Zhycan Cli Commands
    cli.AttachCommands(rootCmd)

	// MARK:Commands --- And New Commands Below ---
	// rootCmd.AddCommand(NewInitCmd())
}
`

	gitIgnoreTmpl = `# Create By Zhycan Framework
#
# Copyright © {{.Year}}
# Project: {{.ProjectName}}
# File: ".gitignore" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
# ---------------------------------------------------

### JetBrains template
# Covers JetBrains IDEs: IntelliJ, RubyMine, PhpStorm, AppCode, PyCharm, CLion, Android Studio, WebStorm and Rider
# Reference: https://intellij-support.jetbrains.com/hc/en-us/articles/206544839

# User-specific stuff
.idea/*

# Gradle and Maven with auto-import
# When using Gradle or Maven with auto-import, you should exclude module files,
# since they will be recreated, and may cause churn.  Uncomment if using
# auto-import.
# .idea/artifacts
# .idea/compiler.xml
# .idea/jarRepositories.xml
# .idea/modules.xml
# .idea/*.iml
# .idea/modules
# *.iml
# *.ipr

# CMake
cmake-build-*/

# Mongo Explorer plugin (remove the comment below to include it)
# .idea/**/mongoSettings.xml

# File-based project format
*.iws

# IntelliJ
out/

# mpeltonen/sbt-idea plugin
.idea_modules/

# JIRA plugin
atlassian-ide-plugin.xml

# Cursive Clojure plugin (remove the comment below to include it)
# .idea/replstate.xml

# Crashlytics plugin (for Android Studio and IntelliJ)
com_crashlytics_export_strings.xml
crashlytics.properties
crashlytics-build.properties
fabric.properties

# Editor-based Rest Client (remove the comment below to include it)
# .idea/httpRequests

# Android studio 3.1+ serialized cache file (remove the comment below to include it)
# .idea/caches/build_file_checksums.ser

### Go template
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

### VisualStudioCode template
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
*.code-workspace

# Local History for Visual Studio Code
.history/

### Xcode template
# Xcode
#
# gitignore contributors: remember to update Global/Xcode.gitignore, Objective-C.gitignore & Swift.gitignore

## User settings
xcuserdata/

## compatibility with Xcode 8 and earlier (ignoring not required starting Xcode 9)
*.xcscmblueprint
*.xccheckout

## compatibility with Xcode 3 and earlier (ignoring not required starting Xcode 4)
build/
DerivedData/
*.moved-aside
*.pbxuser
!default.pbxuser
*.mode1v3
!default.mode1v3
*.mode2v3
!default.mode2v3
*.perspectivev3
!default.perspectivev3

## Gcc Patch
/*.gcno
`

	baseConfigTmpl = `{
  "name": "{{.ProjectName}}",
  "config_must_watched": true,
  "config_remote_addr": "0.0.0.0:7777",
  "config_remote_infra": "grpc",
  "config_remote_duration": 300,
  "modules": [
    {"name":"logger", "type": "local"},
    {"name":"http", "type": "local"},
    {"name":"protobuf", "type": "local"}
  ]
}`
	loggerConfigTmpl = `{
  "type": "zap",
  "outputs": ["console", "file"],
  "channel_size": 1000,
  "options": ["caller", "stackTrace"],
  "console": {
    "level": "debug"
  },
  "file": {
    "level": "debug",
    "path": "/tmp"
  },
  "graylog": {
    "ip": "graylog.graylog",
    "port": 12201,
    "stdout": true
  },
  "syslog": {
    "ip": "172.25.205.37",
    "port": 514,
    "ctype": "tcp"
  }
}`
	httpConfigTmpl = `{
  "default": "s1",
  "servers": [
    {
      "name":                   "s1",
      "addr":                   ":3000",
      "versions":               ["v1"],
      "support_static":         false,
      "conf": {
        "read_timeout": -1,
        "write_timeout": -1,
        "request_methods": ["ALL"]
      },
      "middlewares": {
        "order": ["logger", "cors"],
        "logger": {
          "format": "> [${time}] ${status} - ${latency} ${method} ${path} ${queryParams}\n",
          "time_format": "15:04:05",
          "time_zone": "Local",
          "time_interval": 500,
          "output": "stdout"
        },
        "cors": {
            "allow_all_origins": true,
            "allow_origins": ["*"],
            "allow_methods": ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"],
            "allow_private_network": false,
            "allow_headers": [],
            "allow_credentials": false,
            "expose_headers": [],
            "max_age": 0,
            "allow_wildcard": true,
            "allow_browser_extension": false,
            "custom_schemes": [],
            "allow_websockets": false,
            "allow_files": false,
            "option_response_status_code": 204
        }
      }
    }
  ]
}`

	protobufConfigTmpl = `{
  "proto": 3,
  "servers": [
    "server1"
  ],
  "server1": {
    "host": "0.0.0.0",
    "port": 7777,
    "protocol": "tcp",
    "async": true,
    "reflection": true,
    "configs": {
      "maxReceiveMessageSize": 104857600,
      "maxSendMessageSize": 1048576000
    }
  }
}`

	appControllerTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/controller.go" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/

package app

import (
	"context"
	"fmt"
	"github.com/abolfazlbeh/zhycan/pkg/http"
	"github.com/gin-gonic/gin"
	"{{.ProjectName}}/app/proto/greeter"
)

// MARK: Controller

// SampleController - a sample controller to show the functionality
type SampleController struct{}

// GetName - return the name of the controller to be used as part of the route
func (ctrl *SampleController) GetName() string { return "Sample" }

// Routes - returning controller specific routes to be registered
func (ctrl *SampleController) Routes() []http.HttpRoute {
	return []http.HttpRoute{
		http.HttpRoute{
			Method:    http.MethodGet,
			Path:      "/hello",
			RouteName: "hello",
			F:         ctrl.GetHello,
		},
	}
}

// GetHello - just return the 'Hello World' string to user
func (ctrl *SampleController) GetHello(c *gin.Context) {
	c.String(200, "Hello World")
}

// MARK: gRPC Controller

// SampleProtoController - a sample protobuf controller to show the functionality
type SampleProtoController struct{
	greeter.UnimplementedGreeterServer
}

func (ctrl *SampleProtoController) GetName() string {
	return "Greeter"
}

func (ctrl *SampleProtoController) GetServerNames() []string {
	return []string{"server1"} // it must exist in the "protobuf" config
}

func (ctrl *SampleProtoController) SayHello(ctx context.Context, rq *greeter.HelloRequest) (*greeter.HelloResponse, error) {
	return &greeter.HelloResponse{
		Message: fmt.Sprintf("Hello, %s", rq.Name),
	}, nil
}
`

	appEngineTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/app.go" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/

package app

import (
	"github.com/abolfazlbeh/zhycan/pkg/engine"
	"google.golang.org/grpc"
	"{{.ProjectName}}/app/proto/greeter"
)

// MARK: App Engine

// App - application engine structure that must satisfy one of the engine interface such as 'engine.RestfulApp', ...
type App struct {}

// Init - initialize the app
func (app *App) Init() {
    engine.RegisterRestfulController(&SampleController{})

	greeterService := SampleProtoController{}
	engine.RegisterGrpcController(&greeterService, func(server *grpc.Server) {
		greeter.RegisterGreeterServer(server, &greeterService)
	})
}
`

	appModelTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/model.go" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/

package app

import (
	"errors"
	"github.com/abolfazlbeh/zhycan/pkg/db"
	"gorm.io/gorm"
)

// MARK: Models

// User - a sample model to show the functionality
type User struct {
    gorm.Model
    Name string
}

// CreateNewUser - create a new user record in database
func CreateNewUser(name string) (*User, int64, error) {
    u := User{Name: "test"}

    database, err := db.GetDb("default")
    if err != nil {
        return nil, 0, errors.New("UserCreateError")
    }

    result := database.Create(&u)
    if result.Error != nil {
        return nil, 0, result.Error
    }

    return &u, result.RowsAffected, nil
}

// GetAllUsers - get all user records from database
func GetAllUsers() (*[]User, int64, error) {
    database, err := db.GetDb("default")
    if err != nil {
        return nil, 0, errors.New("UserCreateError")
    }

    var users []User

    result := database.Find(&users)
    if result.Error != nil {
        return nil, 0, result.Error
    }

    return &users, result.RowsAffected, nil
}
`

	greeterProtobufTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/proto/greeter.proto" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/
syntax = "proto3";

package greeter;
option go_package = "{{.ProjectName}}/app/proto/greeter";

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloResponse) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloResponse {
  string message = 1;
}
`

	greeterProtobufPbTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/proto/greeter.proto" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: greeter.proto

package greeter

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The request message containing the user's name.
type HelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string U+0060protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"U+0060
}

func (x *HelloRequest) Reset() {
	*x = HelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greeter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}


func (x *HelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloRequest) ProtoMessage() {}

func (x *HelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloRequest.ProtoReflect.Descriptor instead.
func (*HelloRequest) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{0}
}

func (x *HelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The response message containing the greetings
type HelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string U+0060protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"U+0060
}

func (x *HelloResponse) Reset() {
	*x = HelloResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_greeter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloResponse) ProtoMessage() {}

func (x *HelloResponse) ProtoReflect() protoreflect.Message {
	mi := &file_greeter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloResponse.ProtoReflect.Descriptor instead.
func (*HelloResponse) Descriptor() ([]byte, []int) {
	return file_greeter_proto_rawDescGZIP(), []int{1}
}

func (x *HelloResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_greeter_proto protoreflect.FileDescriptor

var file_greeter_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x22, 0x22, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x29, 0x0a, 0x0d,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x46, 0x0a, 0x07, 0x47, 0x72, 0x65, 0x65, 0x74,
	0x65, 0x72, 0x12, 0x3b, 0x0a, 0x08, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x15,
	0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x2e,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x1a, 0x5a, 0x18, 0x7a, 0x68, 0x79, 0x63, 0x31, 0x31, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_greeter_proto_rawDescOnce sync.Once
	file_greeter_proto_rawDescData = file_greeter_proto_rawDesc
)

func file_greeter_proto_rawDescGZIP() []byte {
	file_greeter_proto_rawDescOnce.Do(func() {
		file_greeter_proto_rawDescData = protoimpl.X.CompressGZIP(file_greeter_proto_rawDescData)
	})
	return file_greeter_proto_rawDescData
}

var file_greeter_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_greeter_proto_goTypes = []any{
	(*HelloRequest)(nil),  // 0: greeter.HelloRequest
	(*HelloResponse)(nil), // 1: greeter.HelloResponse
}
var file_greeter_proto_depIdxs = []int32{
	0, // 0: greeter.Greeter.SayHello:input_type -> greeter.HelloRequest
	1, // 1: greeter.Greeter.SayHello:output_type -> greeter.HelloResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_greeter_proto_init() }
func file_greeter_proto_init() {
	if File_greeter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_greeter_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*HelloRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_greeter_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*HelloResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_greeter_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_greeter_proto_goTypes,
		DependencyIndexes: file_greeter_proto_depIdxs,
		MessageInfos:      file_greeter_proto_msgTypes,
	}.Build()
	File_greeter_proto = out.File
	file_greeter_proto_rawDesc = nil
	file_greeter_proto_goTypes = nil
	file_greeter_proto_depIdxs = nil
}
`

	greeterProtobufGrpcPbTmpl = `/*
Create By Zhycan Framework

Copyright © {{.Year}}
Project: {{.ProjectName}}
File: "app/proto/greeter.proto" --> {{ .Time.Format .TimeFormat }} by {{.CreatorUserName}}
------------------------------
*/

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.3
// source: greeter.proto

package greeter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Greeter_SayHello_FullMethodName = "/greeter.Greeter/SayHello"
)

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, Greeter_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServer is the server API for Greeter service.
// All implementations must embed UnimplementedGreeterServer
// for forward compatibility
type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloResponse, error)
	mustEmbedUnimplementedGreeterServer()
}

// UnimplementedGreeterServer must be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (UnimplementedGreeterServer) SayHello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedGreeterServer) mustEmbedUnimplementedGreeterServer() {}

// UnsafeGreeterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreeterServer will
// result in compilation errors.
type UnsafeGreeterServer interface {
	mustEmbedUnimplementedGreeterServer()
}

func RegisterGreeterServer(s grpc.ServiceRegistrar, srv GreeterServer) {
	s.RegisterService(&Greeter_ServiceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Greeter_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Greeter_ServiceDesc is the grpc.ServiceDesc for Greeter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Greeter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "greeter.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "greeter.proto",
}
`
)
