package rediskit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock implementation of the Redis client
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Init(name, configPrefix, keyPrefix string) error {
	args := m.Called(name, configPrefix, keyPrefix)
	return args.Error(0)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedisClient) IsInitialized() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) Get(ctx context.Context, key string, val any) error {
	args := m.Called(ctx, key, val)
	return args.Error(0)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	args := m.Called(ctx, key, val, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) SetStruct(ctx context.Context, key string, val any, expiration time.Duration) error {
	args := m.Called(ctx, key, val, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) GetStruct(ctx context.Context, key string, val any) error {
	args := m.Called(ctx, key, val)
	return args.Error(0)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, expiration time.Duration, val ...any) error {
	args := m.Called(ctx, key, expiration, val)
	return args.Error(0)
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string, val any) error {
	args := m.Called(ctx, key, field, val)
	return args.Error(0)
}

// Basic tests for the Client struct methods
func TestNewClient(t *testing.T) {
	client := &Client{}
	assert.NotNil(t, client)
	assert.IsType(t, &Client{}, client)
}

func TestClient_generateKey(t *testing.T) {
	client := &Client{prefix: "test"}
	key := client.generateKey("example")
	assert.Equal(t, "test$example", key)

	client = &Client{prefix: ""}
	key = client.generateKey("example")
	assert.Equal(t, "example", key)
}
