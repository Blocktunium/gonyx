package rediskit

import (
	"errors"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"log"
	"sync"
	"time"
)

// Manager controls multiple Redis client instances
type Manager struct {
	name                 string
	lock                 sync.Mutex
	isManagerInitialized bool
	clients              map[string]IRedisClient
}

var (
	managerInstance *Manager
	once            sync.Once
)

// init function is called when the package is imported
func init() {
	log.Println("Initializing RedisKit Manager...")
}

// initialize sets up the manager configuration
func (m *Manager) initialize() {
	m.name = "rediskit"
	m.isManagerInitialized = false

	m.lock.Lock()
	defer m.lock.Unlock()

	prefix := config.GetManager().GetName()
	if prefix == "" {
		return
	}

	// Read Redis connections configuration
	connectionsObj, err := config.GetManager().Get(m.name, "connections")
	if err != nil {
		return
	}

	m.clients = make(map[string]IRedisClient)

	for _, item := range connectionsObj.([]interface{}) {
		instanceName := item.(string)

		redisType, err := config.GetManager().Get(m.name, fmt.Sprintf("%s.%s", instanceName, "type"))
		if err != nil {
			return
		}

		withPrefix, err := config.GetManager().Get(m.name, fmt.Sprintf("%s.%s", instanceName, "add_service_prefix"))
		if err != nil {
			return
		}
		withPrefixBool := withPrefix.(bool)

		logger, _ := logger.GetManager().GetLogger()

		// Register callback for config changes
		wrapper, err := config.GetManager().GetConfigWrapper(m.name)
		if err == nil {
			wrapper.RegisterChangeCallback(func() interface{} {
				err := m.Release()
				if err == nil {
					m.initialize()
				}
				return nil
			})
		}

		// Initialize client based on Redis type
		if redisType == "client" {
			client := &Client{}
			prefix := ""
			if withPrefixBool {
				prefix = config.GetManager().GetName()
			}

			err = client.Init(m.name, fmt.Sprintf("%s.%s", instanceName, redisType), prefix)
			if err != nil {
				if logger != nil {
					logger.Log(types.NewLogObject(types.ERROR, "RedisKit.Manager", redisMaintenanceType,
						time.Now(), "Redis Client initialization failed", err))
				}
				return
			}

			m.clients[instanceName] = client
		} else if redisType == "cluster" {
			// TODO: Add support for Redis cluster
			if logger != nil {
				logger.Log(types.NewLogObject(types.WARNING, "RedisKit.Manager", redisMaintenanceType,
					time.Now(), "Redis Cluster not implemented yet", nil))
			}
		}
	}

	m.isManagerInitialized = true
}

// restartOnConfigChange sets up a listener for configuration changes
func (m *Manager) restartOnConfigChange() {
	wrapper, err := config.GetManager().GetConfigWrapper(m.name)
	if err == nil {
		wrapper.RegisterChangeCallback(func() interface{} {
			if m.isManagerInitialized {
				m.initialize()
			}
			return nil
		})
	}
}

// GetManager returns the singleton instance of the Redis manager
func GetManager() *Manager {
	once.Do(func() {
		managerInstance = &Manager{}
		managerInstance.initialize()
		managerInstance.restartOnConfigChange()
	})
	return managerInstance
}

// Release closes all Redis client connections
func (m *Manager) Release() error {
	if m.clients != nil {
		for _, client := range m.clients {
			err := client.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetClient returns a Redis client by name
func (m *Manager) GetClient(name string) (IRedisClient, error) {
	if m.clients != nil {
		if client, ok := m.clients[name]; ok {
			if client.IsInitialized() {
				return client, nil
			}
		}
	}

	return nil, NewError(errors.New("redis client not found or not initialized"))
}

// IsInitialized checks if the manager has been initialized
func (m *Manager) IsInitialized() bool {
	return m.isManagerInitialized
}
