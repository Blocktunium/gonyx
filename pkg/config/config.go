package config

import (
	"github.com/Blocktunium/gonyx/internal/config"
	"time"
)

// InitializeManager - Create a new config manager instance and wait to initialize it
func InitializeManager(configBasePath string, configInitialMode string, configEnvPrefix string) error {
	err := config.CreateManager(configBasePath, configInitialMode, configEnvPrefix)
	if err != nil {
		return err
	}

	for {
		flag := config.GetManager().IsInitialized()
		if flag {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// GetName - returns service instance name based on config
func GetName() string {
	if config.GetManager().IsInitialized() {
		return config.GetManager().GetName()
	}

	return ""
}

// GetOperationType - returns operation type which could be `dev`, `prod`
func GetOperationType() string {
	if config.GetManager().IsInitialized() {
		return config.GetManager().GetOperationType()
	}

	return ""
}

// GetHostName - returns hostname based on config
func GetHostName() string {
	if config.GetManager().IsInitialized() {
		return config.GetManager().GetHostName()
	}

	return ""
}

// Get - get value of the key in specific category
func Get(category string, name string) (interface{}, error) {
	return config.GetManager().Get(category, name)
}

// Set - set value in category by specified key
func Set(category string, name string, value interface{}) error {
	return config.GetManager().Set(category, name, value)
}

// IsInitialized - iterate over all config wrappers and see all initialised correctly
func IsInitialized() bool {
	return config.GetManager().IsInitialized()
}

// GetAllInitializedModuleList - get list of names that initialized truly
func GetAllInitializedModuleList() []string {
	return config.GetManager().GetAllInitializedModuleList()
}

// ManualLoadConfig - load manual config from the path and add to the current dict
func ManualLoadConfig(configBasePath string, configName string) error {
	return config.GetManager().ManualLoadConfig(configBasePath, configName)
}
