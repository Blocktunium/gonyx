package grpc

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

// MARK: Variables

var (
	ServerMaintenanceType = types.NewLogType("PROTOBUF_SERVER_MAINTENANCE")
)

// ServerWrapper struct
type ServerWrapper struct {
	name        string
	grpcServer  *grpc.Server
	listener    net.Listener
	initialized bool
	config      ServerConfig
	//authObj     *auth.Authentication
	//authEnable bool
}

// init - initialize the ServerWrapper with the configs
func (s *ServerWrapper) init(name string, serverConfig ServerConfig) error {
	s.name = name
	s.initialized = false

	s.config = serverConfig

	lis, err := net.Listen(s.config.Protocol, fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return err
	}

	//authEnable, err := config.GetManager().Get(parts[0], parts[1]+".auth_enable")
	//if err != nil {
	//	return err
	//}
	//s.authEnable = authEnable.(bool)

	//s.authObj = &auth.Authentication{}
	//s.authObj.Init()

	options := s.generateConfigs(s.config.Configs)

	s.listener = lis
	s.grpcServer = grpc.NewServer(options...)
	if s.config.Reflection {
		reflection.Register(s.grpcServer)
	}

	s.initialized = true
	return nil
}

// generateConfigs - generate grpc.ServerOption array from configs
func (s *ServerWrapper) generateConfigs(configs map[string]interface{}) []grpc.ServerOption {
	var options []grpc.ServerOption
	if v, ok := configs["maxreceivemessagesize"]; ok {
		options = append(options, grpc.MaxRecvMsgSize(int(v.(float64))))
	}

	if v, ok := configs["maxsendmessagesize"]; ok {
		options = append(options, grpc.MaxSendMsgSize(int(v.(float64))))
	}

	options = append(options, grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    15 * time.Minute,
		Timeout: 20 * time.Second,
	}))

	return options
}

// MARK: Public functions

// NewServer - create a new instance of Server and return it
func NewServer(name string, config ServerConfig) (*ServerWrapper, error) {
	server := &ServerWrapper{}
	err := server.init(name, config)
	if err != nil {
		return nil, NewCreateServerErr(err)
	}

	return server, nil
}

// Start - start the server with option of async capability
func (s *ServerWrapper) Start(ch *chan error) error {
	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.INFO, "protobuf.Server.Start", ServerMaintenanceType, time.Now(), "Starting the gRPC server ...", s.listener))
	}

	if s.config.Async {
		go func(ch1 *chan error) {
			err := s.grpcServer.Serve(s.listener)
			if err != nil && ch1 != nil {
				*ch <- NewGrpcServerStartError(err)
			}
		}(ch)
	} else {
		err := s.grpcServer.Serve(s.listener)
		if err != nil {
			return NewGrpcServerStartError(err)
		}
	}

	return nil
}

// Stop - stop the server
func (s *ServerWrapper) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
}

// IsInitialized - return whether the server is started or not
func (s *ServerWrapper) IsInitialized() bool {
	return s.initialized
}

func (s *ServerWrapper) RegisterController(desc *grpc.ServiceDesc, realClass interface{}) error {
	if realClass != nil {
		s.grpcServer.RegisterService(desc, realClass)
		return nil
	}

	return NewNilServiceRegistryError()
}

func (s *ServerWrapper) GetGrpcServer() *grpc.Server {
	return s.grpcServer
}
