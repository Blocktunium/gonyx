package cache

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

// Mark: manager

// Manager object
type manager struct {
	name                 string
	lock                 sync.Mutex
	isManagerInitialized bool

	caches map[string]ICache
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// MARK: Module Initializer
func init() {
	log.Println("Initializing Cache Manager ...")
}

// init - Manager Constructor - It initializes the manager configuration params
func (m *manager) init() {
	m.name = "cache"
	m.isManagerInitialized = false

	m.lock.Lock()
	defer m.lock.Unlock()

	prefix := config.GetManager().GetName()
	if prefix == "" {
		return
	}

	// read configs
	connectionsObj, err := config.GetManager().Get(m.name, "connections")
	if err != nil {
		return
	}

	m.caches = make(map[string]ICache)

	for _, item := range connectionsObj.([]interface{}) {
		cacheInstanceName := item.(string)

		cacheType, err := config.GetManager().Get(m.name, fmt.Sprintf("%s.%s", cacheInstanceName, "type"))
		if err != nil {
			return
		}

		withPrefix, err := config.GetManager().Get(m.name, fmt.Sprintf("%s.%s", cacheInstanceName, "add_service_prefix"))
		if err != nil {
			return
		}
		withPrefixBool := withPrefix.(bool)

		logge, _ := logger.GetManager().GetLogger()

		wrapper, err := config.GetManager().GetConfigWrapper(m.name)
		if err == nil {
			wrapper.RegisterChangeCallback(func() interface{} {
				err := m.Release()
				if err == nil {
					m.init()
				}
				return nil
			})
		}

		if cacheType == "redis" {
			redisType, err := config.GetManager().Get(m.name, fmt.Sprintf("%s.%s", cacheInstanceName, "redis_type"))
			if err != nil {
				return
			}

			if redisType == "cluster" {
				// TODO: check for cluster config

			} else if redisType == "client" {
				tempCache := &RedisClientCache{}
				p := ""
				if withPrefixBool == true {
					p = prefix
				}
				err = tempCache.Init(m.name, fmt.Sprintf("%s.%s", cacheInstanceName, redisType), p)
				if err != nil {
					if logge != nil {
						logge.Log(types.NewLogObject(types.ERROR, "Cache.Manager", types.NilObject,
							time.Now(), "Redis Client Object is nil", err))
					}
					return
				}

				m.caches[cacheInstanceName] = tempCache
			}
		}
	}

	m.isManagerInitialized = true
}

// restartOnChangeConfig - subscribe a function for when the config is changed
func (m *manager) restartOnChangeConfig() {
	// Config config server to reload
	wrapper, err := config.GetManager().GetConfigWrapper(m.name)
	if err == nil {
		wrapper.RegisterChangeCallback(func() interface{} {
			if m.isManagerInitialized {
				m.init()
			}
			return nil
		})
	} else {
		// TODO: make some logs
	}
}

// MARK: Public Functions

// GetManager - This function returns singleton instance of Cache Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}
		managerInstance.init()
		managerInstance.restartOnChangeConfig()
	})
	return managerInstance
}

// Release receiver - releases the cache instance resource
func (m *manager) Release() error {
	if m.caches != nil {
		for _, cache := range m.caches {
			err := cache.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetCache - This function returns the instance of cache (interface of it)
func (m *manager) GetCache(cacheName string) (ICache, error) {
	if m.caches != nil {
		if val, ok := m.caches[cacheName]; ok {
			if val.IsInitialized() {
				return val, nil
			}
		}
	}

	return nil, NewError(errors.New("cannot get instance of cache"))
}
