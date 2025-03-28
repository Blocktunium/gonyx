package mongokit

import (
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/Blocktunium/gonyx/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"sync"
)

// Mark: manager

// manager object
type manager struct {
	name             string
	lock             sync.Mutex
	mongoDbInstances map[string]*MongoWrapper
	supportedDBs     []string

	isManagerInitialized bool
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// Module init function
func init() {
	log.Println("DB Manager Package Initialized...")
}

// init - Manager Constructor - It initializes the manager configuration params
func (m *manager) init() {
	m.name = "db"
	m.isManagerInitialized = false

	m.lock.Lock()
	defer m.lock.Unlock()

	m.supportedDBs = []string{"mongodb"}

	// read configs
	connectionsObj, err := config.GetManager().Get(m.name, "connections")
	if err != nil {
		return
	}

	m.mongoDbInstances = make(map[string]*MongoWrapper)

	for _, item := range connectionsObj.([]interface{}) {
		dbInstanceName := item.(string)

		dbTypeKey := fmt.Sprintf("%s.%s", dbInstanceName, "type")
		dbTypeInf, err := config.GetManager().Get(m.name, dbTypeKey)
		if err != nil {
			continue
		}

		//  create a new instance based on type
		dbType := strings.ToLower(dbTypeInf.(string))
		if utils.ArrayContains(&m.supportedDBs, dbType) {
			switch dbType {
			case "mongodb":
				obj, err := NewMongoWrapper(fmt.Sprintf("db/%s", dbInstanceName))
				if err != nil {
					// TODO: log error here
					continue
				}
				m.mongoDbInstances[dbInstanceName] = obj
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

// GetManager - This function returns singleton instance of Db Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}
		managerInstance.init()
		managerInstance.restartOnChangeConfig()
	})
	return managerInstance
}

// GetMongoDb - Get *mongo.Client instance from the underlying interfaces
func (m *manager) GetMongoDb(instanceName string) (*mongo.Database, error) {
	if m.isManagerInitialized {
		if v, ok := m.mongoDbInstances[instanceName]; ok {
			return v.GetDb()
		}
	}
	return nil, NewNotExistServiceNameErr(instanceName)
}

func (m *manager) RegisterLogger(l types.Logger) {
	//for _, item := range m.mongoDbInstances {
	//}
}
