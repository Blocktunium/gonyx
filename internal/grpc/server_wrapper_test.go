package grpc

import (
	"github.com/Blocktunium/gonyx/internal/config"
	"testing"
)

func TestServerWrapper_Init(t *testing.T) {
	makeReadyConfigManager()

	serverConfig := ServerConfig{
		Host:     "127.0.0.1",
		Port:     7777,
		Protocol: "tcp",
		Async:    true,
		Configs:  map[string]interface{}{},
	}

	server, err := NewServer("protobuf", serverConfig)
	if err != nil {
		t.Errorf("Creating gRPC Server --> Expected: %v, but got %v", nil, err)
		return
	}

	err = server.Start(nil)
	if err != nil {
		t.Errorf("Starting gRPC Server --> Expected: %v, but got %v", nil, err)
		return
	}
}

func makeReadyConfigManager() {
	path := "../.."
	initialMode := "test"
	prefix := "Gonyx"

	_ = config.CreateManager(path, initialMode, prefix)
}
