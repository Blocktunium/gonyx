package gormkit

import (
	"github.com/Blocktunium/gonyx/internal/config"
	"gorm.io/gorm"
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

func TestManager_CheckInitialization(t *testing.T) {
	makeReadyConfigManager()

	m := manager{}
	m.init()

	if len(m.sqliteDbInstances) != 1 {
		t.Errorf("Expected manager have %v instance of sqlite, but got %v", 1, len(m.sqliteDbInstances))
		return
	}

	if v, ok := m.sqliteDbInstances["server1"]; ok {
		expected := reflect.TypeOf(&SqlWrapper[Sqlite]{})
		got := reflect.ValueOf(v).Type()
		if got != expected {
			t.Errorf("Expected manager one sqlite instance: %v, but got: %v ", expected, got)
			return
		}
	} else {
		t.Errorf("Expected manager have %v instance of sqlite for %v, but got nothing", 1, "server1")
		return
	}
}

func TestManager_TestGetDbFunc(t *testing.T) {
	makeReadyConfigManager()

	m := manager{}
	m.init()

	if len(m.sqliteDbInstances) != 1 {
		t.Errorf("Expected manager have %v instance of sqlite, but got %v", 1, len(m.sqliteDbInstances))
		return
	}

	db, err := m.GetDb("server1")
	if err != nil {
		t.Errorf("Get Db Instance --> Expected error: %v, but got %v", nil, err)
		return
	}

	expected := reflect.TypeOf(&gorm.DB{})
	got := reflect.ValueOf(db).Type()

	if got != expected {
		t.Errorf("Expected manager GetDb function return: %v, but got: %v ", expected, got)
		return
	}
}

func makeReadyConfigManager() {
	path := "../.."
	initialMode := "test"
	prefix := "Gonyx"

	_ = config.CreateManager(path, initialMode, prefix)
}
