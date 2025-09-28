package rediskit

import (
	"errors"
	"sync"
	"testing"

	"github.com/Blocktunium/gonyx/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock implementation of the config manager
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) Get(name string, key string) (interface{}, error) {
	args := m.Called(name, key)
	return args.Get(0), args.Error(1)
}

func (m *MockConfig) GetName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfig) GetConfigWrapper(name string) (config.IConfigWrapper, error) {
	args := m.Called(name)
	return args.Get(0).(config.IConfigWrapper), args.Error(1)
}

// MockConfigWrapper is a mock implementation of the config wrapper
type MockConfigWrapper struct {
	mock.Mock
}

func (m *MockConfigWrapper) RegisterChangeCallback(callback func() interface{}) {
	m.Called(callback)
}

//// MockRedisClient is a mock implementation of the Redis client
//type MockRedisClient struct {
//	mock.Mock
//}
//
//func (m *MockRedisClient) Init(name, configPrefix, keyPrefix string) error {
//	args := m.Called(name, configPrefix, keyPrefix)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) Ping(ctx context.Context) error {
//	args := m.Called(ctx)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) IsInitialized() bool {
//	args := m.Called()
//	return args.Bool(0)
//}
//
//func (m *MockRedisClient) Close() error {
//	args := m.Called()
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) Get(ctx context.Context, key string, val any) error {
//	args := m.Called(ctx, key, val)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
//	args := m.Called(ctx, key, val, expiration)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) SetStruct(ctx context.Context, key string, val any, expiration time.Duration) error {
//	args := m.Called(ctx, key, val, expiration)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) GetStruct(ctx context.Context, key string, val any) error {
//	args := m.Called(ctx, key, val)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) HSet(ctx context.Context, key string, expiration time.Duration, val ...any) error {
//	args := m.Called(ctx, key, expiration, val)
//	return args.Error(0)
//}
//
//func (m *MockRedisClient) HGet(ctx context.Context, key, field string, val any) error {
//	args := m.Called(ctx, key, field, val)
//	return args.Error(0)
//}

// Replace the real config manager with our mock for testing
func setupMockConfig() (*MockConfig, *MockConfigWrapper) {
	mockConfig := new(MockConfig)
	mockWrapper := new(MockConfigWrapper)

	// Save the original manager instance
	originalManager := config.GetManager

	// Replace with mock
	config.GetManager = func() config.IManager {
		return mockConfig
	}

	// Return a clean-up function
	return mockConfig, mockWrapper
}

// Restore the original config manager
func teardownMockConfig(originalManager func() config.IManager) {
	config.GetManager = originalManager
}

// Test creating a new Manager instance
func TestGetManager(t *testing.T) {
	// Setup
	originalManagerFn := config.GetManager
	mockConfig, mockWrapper := setupMockConfig()
	defer teardownMockConfig(originalManagerFn)

	// Reset the singleton for testing
	managerInstance = nil
	once = sync.Once{}

	// Configure mocks
	mockConfig.On("GetName").Return("test-service")
	mockConfig.On("Get", "rediskit", "connections").Return([]interface{}{"main"}, nil)
	mockConfig.On("Get", "rediskit", "main.type").Return("client", nil)
	mockConfig.On("Get", "rediskit", "main.add_service_prefix").Return(false, nil)
	mockConfig.On("GetConfigWrapper", "rediskit").Return(mockWrapper, nil)
	mockWrapper.On("RegisterChangeCallback", mock.AnythingOfType("func() interface {}")).Return()
	mockConfig.On("Get", "rediskit", "main.client.address").Return("localhost:6379", nil)
	mockConfig.On("Get", "rediskit", "main.client.password").Return("", nil)
	mockConfig.On("Get", "rediskit", "main.client.db").Return(float64(0), nil)
	mockConfig.On("Get", "rediskit", "main.client.max_retries").Return(float64(3), nil)
	mockConfig.On("Get", "rediskit", "main.client.min_retry_backoff").Return(float64(100), nil)
	mockConfig.On("Get", "rediskit", "main.client.max_retry_backoff").Return(float64(500), nil)
	mockConfig.On("Get", "rediskit", "main.client.dial_timeout").Return(float64(5000), nil)
	mockConfig.On("Get", "rediskit", "main.client.read_timeout").Return(float64(3000), nil)
	mockConfig.On("Get", "rediskit", "main.client.write_timeout").Return(float64(3000), nil)
	mockConfig.On("Get", "rediskit", "main.client.on_connect_log").Return(true, nil)
	mockConfig.On("Get", "rediskit", "main.client.enable_lock").Return(false, nil)

	// Test
	manager := GetManager()

	// Assertions
	assert.NotNil(t, manager)
	assert.True(t, manager.isManagerInitialized)

	// Call again to test singleton
	manager2 := GetManager()
	assert.Same(t, manager, manager2)

	// Verify mocks
	mockConfig.AssertExpectations(t)
	mockWrapper.AssertExpectations(t)
}

// Test getting a client from the manager
func TestManager_GetClient(t *testing.T) {
	// Create a manager with a mock client
	manager := &Manager{
		name:                 "rediskit",
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test": &MockRedisClient{},
		},
	}

	// Setup mock client
	mockClient := manager.clients["test"].(*MockRedisClient)
	mockClient.On("IsInitialized").Return(true)

	// Test getting an existing client
	client, err := manager.GetClient("test")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test getting a non-existent client
	client, err = manager.GetClient("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, client)

	// Test getting an uninitialized client
	mockClient.On("IsInitialized").Return(false)
	client, err = manager.GetClient("test")
	assert.Error(t, err)
	assert.Nil(t, client)

	// Verify mock
	mockClient.AssertExpectations(t)
}

// Test releasing all clients
func TestManager_Release(t *testing.T) {
	// Create client mocks
	mockClient1 := new(MockRedisClient)
	mockClient2 := new(MockRedisClient)

	// Setup manager with mock clients
	manager := &Manager{
		name:                 "rediskit",
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"client1": mockClient1,
			"client2": mockClient2,
		},
	}

	// Setup expectations
	mockClient1.On("Close").Return(nil)
	mockClient2.On("Close").Return(nil)

	// Test successful release
	err := manager.Release()
	assert.NoError(t, err)

	// Test with error
	mockClient1.On("Close").Return(nil)
	mockClient2.On("Close").Return(errors.New("close error"))

	manager.clients = map[string]IRedisClient{
		"client1": mockClient1,
		"client2": mockClient2,
	}

	// Test release with error
	err = manager.Release()
	assert.Error(t, err)
	assert.Equal(t, "close error", err.Error())

	// Verify mocks
	mockClient1.AssertExpectations(t)
	mockClient2.AssertExpectations(t)
}

// Test manager initialization with various scenarios
func TestManager_initialize(t *testing.T) {
	// Setup
	originalManagerFn := config.GetManager
	mockConfig, mockWrapper := setupMockConfig()
	defer teardownMockConfig(originalManagerFn)

	// Test case: empty config name
	mockConfig.On("GetName").Return("").Once()

	manager := &Manager{}
	manager.initialize()
	assert.False(t, manager.isManagerInitialized)

	// Test case: missing connections
	mockConfig.On("GetName").Return("test").Once()
	mockConfig.On("Get", "rediskit", "connections").Return(nil, errors.New("config error")).Once()

	manager = &Manager{}
	manager.initialize()
	assert.False(t, manager.isManagerInitialized)

	// Test case: invalid Redis type
	mockConfig.On("GetName").Return("test").Once()
	mockConfig.On("Get", "rediskit", "connections").Return([]interface{}{"main"}, nil).Once()
	mockConfig.On("Get", "rediskit", "main.type").Return(nil, errors.New("missing type")).Once()

	manager = &Manager{}
	manager.initialize()
	assert.False(t, manager.isManagerInitialized)

	// Test case: missing add_service_prefix
	mockConfig.On("GetName").Return("test").Once()
	mockConfig.On("Get", "rediskit", "connections").Return([]interface{}{"main"}, nil).Once()
	mockConfig.On("Get", "rediskit", "main.type").Return("client", nil).Once()
	mockConfig.On("Get", "rediskit", "main.add_service_prefix").Return(nil, errors.New("missing prefix")).Once()

	manager = &Manager{}
	manager.initialize()
	assert.False(t, manager.isManagerInitialized)

	// Verify mocks
	mockConfig.AssertExpectations(t)
	mockWrapper.AssertExpectations(t)
}

// Test the IsInitialized method
func TestManager_IsInitialized(t *testing.T) {
	manager := &Manager{isManagerInitialized: true}
	assert.True(t, manager.IsInitialized())

	manager = &Manager{isManagerInitialized: false}
	assert.False(t, manager.IsInitialized())
}
