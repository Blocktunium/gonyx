package mongokit

import (
	"github.com/Blocktunium/gonyx/internal/config"
	"reflect"
	"testing"
)

func TestManager_Init(t *testing.T) {
	makeReadyConfigManager()

	m := manager{}
	m.init()

	if m.name != "db" {
		t.Errorf("Expected manager name to be 'db', got '%s'", m.name)
	}
}

func TestManager_CheckInitializationMongo(t *testing.T) {
	makeReadyConfigManager()

	m := manager{}
	m.init()

	if len(m.mongoDbInstances) != 1 {
		t.Errorf("Expected manager have %v instance of mongodb, but got %v", 1, len(m.mongoDbInstances))
		return
	}

	if v, ok := m.mongoDbInstances["server4"]; ok {
		expected := reflect.TypeOf(&MongoWrapper{})
		got := reflect.ValueOf(v).Type()
		if got != expected {
			t.Errorf("Expected manager one mongodb instance: %v, but got: %v ", expected, got)
			return
		}
	} else {
		t.Errorf("Expected manager have %v instance of mongodb for %v, but got nothing", 1, "server1")
		return
	}
}

func makeReadyConfigManager() {
	path := "../.."
	initialMode := "test"
	prefix := "Gonyx"

	_ = config.CreateManager(path, initialMode, prefix)
}
