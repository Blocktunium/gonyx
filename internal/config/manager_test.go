package config

import (
	"reflect"
	"testing"
	"time"
)

func TestCreateManager(t *testing.T) {
	err := createManager()
	if err != nil {
		t.Errorf("There is an error in creating manager: %v", err)
	}
}

func TestAllModuleIsInitialized(t *testing.T) {
	err := createManager()
	if err == nil {
		time.Sleep(2 * time.Second)

		expectedFlag := true
		actualFlag := GetManager().IsInitialized()
		if expectedFlag != actualFlag {
			t.Errorf("Manager initialization status --> expected: %v, got %v", expectedFlag, actualFlag)
		} else {
			expectedInitializedModules := []string{"logger"}
			actualInitializedModules := GetManager().GetAllInitializedModuleList()
			if !reflect.DeepEqual(actualInitializedModules, expectedInitializedModules) {
				t.Errorf("Config Initialized Modules --> expected: %v, got %v", expectedInitializedModules, actualInitializedModules)
			}
		}
	}
}

func TestGetWrapperInstance(t *testing.T) {
	err := createManager()
	if err == nil {
		time.Sleep(3 * time.Second)

		expectedWrapper := &ViperWrapper{
			ConfigPath:          []string{"../../configs/test/"},
			ConfigName:          "logger",
			ConfigEnvPrefix:     "",
			ConfigResourcePlace: "",
		}
		_ = expectedWrapper.Load()

		actualWrapper, err := GetManager().GetConfigWrapper("logger")
		if err != nil {
			t.Errorf("Get Warpper Object --> Got Error --> Expected: %v, but got, %v", nil, err)
		}

		if !reflect.DeepEqual(expectedWrapper, actualWrapper) {
			t.Errorf("Get Warpper Object --> Expected: %v, but got %v", expectedWrapper, actualWrapper)
		}
	}
}

func TestGetDifferentConfigs(t *testing.T) {
	expectedValues := []struct {
		key         string
		expectedVal interface{}
	}{
		{key: "type", expectedVal: "zap"},
		{key: "typ2", expectedVal: KeyNotExistErr{
			Key: "typ2", Category: "logger",
		}},
	}

	err := createManager()
	if err != nil {
		t.Errorf("Creating Config Manager Caused Error --> Expected: %v, but got %v", nil, err)
	}

	actualVal, _ := GetManager().Get("logger", expectedValues[0].key)
	if expectedValues[0].expectedVal != actualVal {
		t.Errorf("Get the key for the `type` --> Expected: %v, got %v", expectedValues[0].expectedVal, actualVal)
	}

	_, err = GetManager().Get("logger", expectedValues[1].key)
	if reflect.DeepEqual(expectedValues[1].expectedVal, err) {
		t.Errorf("Get the key for the `type` --> Expected: %v, got %v", expectedValues[1].expectedVal, err)
	}
}

func TestManager_ManualLoadConfig(t *testing.T) {
	err := createManager()
	if err != nil {
		t.Errorf("Creating Config Manager Caused Error --> Expected: %v, but got %v", nil, err)
	}

	// add logger to the manager
	err = GetManager().ManualLoadConfig("/tmp/", "zh")
	if err != nil {
		t.Errorf("Manual loading config file --> Expected: %v, but got %v", nil, err)
	}

	expectedValue := "founded"
	val, err := GetManager().Get("zh", "test")
	if err != nil {
		t.Errorf("Get config value from manual file the key is `test`, got error --> Expected: %v, but got %v", nil, err)
	}

	if !reflect.DeepEqual(expectedValue, val) {
		t.Errorf("Expected the value of key `test` to be `%v`, but got `%v`", expectedValue, val)
	}
}

func createManager() error {
	path := "../.."
	initialMode := "test"
	prefix := "Gonyx"

	return CreateManager(path, initialMode, prefix)
}
