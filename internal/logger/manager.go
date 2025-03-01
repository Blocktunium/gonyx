package logger

// Imports needed list
import (
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"log"
	"sync"
)

// Mark: manager

// Manager object
type manager struct {
	name   string
	logger types.Logger
	lock   sync.Mutex
}

// MARK: Module variables
var managerInstance *manager = nil
var once sync.Once

// Module init function
func init() {
	log.Println("Logger Manager Package Initialized...")
}

// init - Manager Constructor - It initializes the logger configuration params
func (m *manager) init() {
	m.name = "logger"
	m.lock.Lock()
	defer m.lock.Unlock()

	t, err := config.GetManager().Get(m.name, "type")
	if err != nil {
		return
	}

	if t == "zap" {
		m.logger = &ZapWrapper{}
		m.logger.Constructor(m.name)
	} else if t == "logme" {
		m.logger = &LogMeWrapper{}
		m.logger.Constructor(m.name)
	}

	// Config config server to reload
	wrapper, err := config.GetManager().GetConfigWrapper(m.name)
	if err == nil {
		wrapper.RegisterChangeCallback(func() interface{} {
			return nil
		})
	}

	return
}

// MARK: Public Functions

// GetManager - This function returns singleton instance of Logger Manager
func GetManager() *manager {
	// once used for prevent race condition and manage critical section.
	once.Do(func() {
		managerInstance = &manager{}
		managerInstance.init()
	})
	return managerInstance
}

// GetLogger - This function returns logger instance
func (m *manager) GetLogger() (types.Logger, *Error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.logger != nil {
		return m.logger, nil
	}
	return nil, NewError(nil)
}
