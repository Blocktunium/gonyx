package grpc

import (
	"encoding/json"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"log"
	"sync"
	"time"
)

// MARK: Variables

var (
	gRPCMaintenanceType = types.NewLogType("PROTOBUF_MANAGER_MAINTENANCE")
)

// Mark: manager

// manager struct
type manager struct {
	name      string
	lock      sync.Mutex
	servers   map[string]*ServerWrapper
	isStarted bool
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// Module init function
func init() {
	log.Println("gRPC Manager Package Initialized...")
}

// Manager Constructor - It initializes the db configuration params
func (m *manager) init() {
	m.name = "protobuf"

	m.lock.Lock()
	defer m.lock.Unlock()

	servers, err := config.GetManager().Get(m.name, "servers")
	if err != nil {
		return
	}

	serverArray := make([]string, len(servers.([]interface{})))
	for i, v := range servers.([]interface{}) {
		serverArray[i] = v.(string)
	}

	m.servers = make(map[string]*ServerWrapper)

	for _, item := range serverArray {
		conf, err := config.GetManager().Get(m.name, item)
		if err != nil {
			continue
		}

		jsonBody, err2 := json.Marshal(conf)
		if err2 != nil {
			continue
		}

		var obj ServerConfig
		err = json.Unmarshal(jsonBody, &obj)
		if err == nil {
			s, err := NewServer(m.name+"."+item, obj)
			if err == nil {
				m.servers[item] = s
			}
		}
	}

	// Config config server to reload
	wrapper, err := config.GetManager().GetConfigWrapper(m.name)
	if err == nil {
		wrapper.RegisterChangeCallback(func() interface{} {
			return nil
		})
	}
}

// MARK: Public Functions

// GetManager - This function returns singleton instance of gRPC Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}
		managerInstance.init()
	})
	return managerInstance
}

// StartServers - This function starts the gRPC servers
func (m *manager) StartServers() {
	l, _ := logger.GetManager().GetLogger()

	m.lock.Lock()
	defer m.lock.Unlock()

	ch := make(chan error, len(m.servers))

	for _, item := range m.servers {
		if item.IsInitialized() {
			err := item.Start(&ch)
			if err != nil {
				if l != nil {
					l.Log(types.NewLogObject(types.ERROR, "protobuf.manager.StartServer", gRPCMaintenanceType, time.Now(), "Cannot start server ...", err))
				}
			}
		}
	}

	go func(ch1 *chan error) {
		select {
		case err := <-*ch1:
			if l != nil {
				l.Log(types.NewLogObject(types.ERROR, "protobuf.manager.StartServer", gRPCMaintenanceType, time.Now(), "Cannot start server ...", err))
			}
		}
	}(&ch)

	if l != nil {
		l.Log(types.NewLogObject(types.INFO, "protobuf.manager.StartServer", gRPCMaintenanceType, time.Now(), "gRPC Engine Started ...", nil))
	}

	m.isStarted = true
}

// StopServers - This function stops gRPC server
func (m *manager) StopServers() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, item := range m.servers {
		item.Stop()
	}
	m.isStarted = false
}

// GetServerByName - get server instance by its name
func (m *manager) GetServerByName(name string) (*ServerWrapper, error) {
	if v, ok := m.servers[name]; ok {
		return v, nil
	}

	return nil, NewGrpcServerNotExistError(name)
}
