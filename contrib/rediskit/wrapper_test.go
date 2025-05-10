package rediskit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test the Set wrapper function
func TestSet(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)
	mockClient.On("Set", mock.Anything, "test-key", "test-value", time.Minute).Return(nil)

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the Set function
	ctx := context.Background()
	err := Set(ctx, "test-client", "test-key", "test-value", time.Minute)
	assert.NoError(t, err)

	// Test with an error from the client
	mockClient.On("Set", mock.Anything, "error-key", "error-value", time.Minute).Return(errors.New("set error"))
	err = Set(ctx, "test-client", "error-key", "error-value", time.Minute)
	assert.Error(t, err)
	assert.Equal(t, "set error", err.Error())

	// Test with a non-existent client
	err = Set(ctx, "non-existent", "test-key", "test-value", time.Minute)
	assert.Error(t, err)

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the Get wrapper function
func TestGet(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)
	mockClient.On("Get", mock.Anything, "test-key", mock.AnythingOfType("*string")).Return(nil).Run(
		func(args mock.Arguments) {
			val := args.Get(2).(*string)
			*val = "test-value"
		})

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the Get function
	ctx := context.Background()
	var value string
	err := Get(ctx, "test-client", "test-key", &value)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", value)

	// Test with an error from the client
	mockClient.On("Get", mock.Anything, "error-key", mock.Anything).Return(errors.New("get error"))
	err = Get(ctx, "test-client", "error-key", &value)
	assert.Error(t, err)
	assert.Equal(t, "get error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the SetStruct wrapper function
func TestSetStruct(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)

	type TestStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	testStruct := TestStruct{ID: 1, Name: "Test"}

	mockClient.On("SetStruct", mock.Anything, "test-key", testStruct, time.Minute).Return(nil)

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the SetStruct function
	ctx := context.Background()
	err := SetStruct(ctx, "test-client", "test-key", testStruct, time.Minute)
	assert.NoError(t, err)

	// Test with an error from the client
	mockClient.On("SetStruct", mock.Anything, "error-key", testStruct, time.Minute).Return(errors.New("set struct error"))
	err = SetStruct(ctx, "test-client", "error-key", testStruct, time.Minute)
	assert.Error(t, err)
	assert.Equal(t, "set struct error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the GetStruct wrapper function
func TestGetStruct(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)

	type TestStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	expectedStruct := TestStruct{ID: 1, Name: "Test"}

	// Setup the mock to set the value when Get is called
	mockClient.On("GetStruct", mock.Anything, "test-key", mock.AnythingOfType("*rediskit.TestStruct")).Return(nil).Run(
		func(args mock.Arguments) {
			val := args.Get(2).(*TestStruct)
			*val = expectedStruct
		})

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the GetStruct function
	ctx := context.Background()
	var result TestStruct
	err := GetStruct(ctx, "test-client", "test-key", &result)
	assert.NoError(t, err)
	assert.Equal(t, expectedStruct, result)

	// Test with an error from the client
	mockClient.On("GetStruct", mock.Anything, "error-key", mock.Anything).Return(errors.New("get struct error"))
	err = GetStruct(ctx, "test-client", "error-key", &result)
	assert.Error(t, err)
	assert.Equal(t, "get struct error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the HSet wrapper function
func TestHSet(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)
	mockClient.On("HSet", mock.Anything, "test-key", time.Minute, mock.Anything).Return(nil)

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the HSet function
	ctx := context.Background()
	err := HSet(ctx, "test-client", "test-key", time.Minute, "field1", "value1", "field2", "value2")
	assert.NoError(t, err)

	// Test with an error from the client
	mockClient.On("HSet", mock.Anything, "error-key", time.Minute, mock.Anything).Return(errors.New("hset error"))
	err = HSet(ctx, "test-client", "error-key", time.Minute, "field1", "value1")
	assert.Error(t, err)
	assert.Equal(t, "hset error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the HGet wrapper function
func TestHGet(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)
	mockClient.On("HGet", mock.Anything, "test-key", "test-field", mock.AnythingOfType("*string")).Return(nil).Run(
		func(args mock.Arguments) {
			val := args.Get(3).(*string)
			*val = "test-value"
		})

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the HGet function
	ctx := context.Background()
	var value string
	err := HGet(ctx, "test-client", "test-key", "test-field", &value)
	assert.NoError(t, err)
	assert.Equal(t, "test-value", value)

	// Test with an error from the client
	mockClient.On("HGet", mock.Anything, "error-key", "error-field", mock.Anything).Return(errors.New("hget error"))
	err = HGet(ctx, "test-client", "error-key", "error-field", &value)
	assert.Error(t, err)
	assert.Equal(t, "hget error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}

// Test the Ping wrapper function
func TestPing(t *testing.T) {
	// Create a mock Redis client
	mockClient := new(MockRedisClient)
	mockClient.On("IsInitialized").Return(true)
	mockClient.On("Ping", mock.Anything).Return(nil)

	// Create a manager with the mock client
	originalManager := managerInstance
	managerInstance = &Manager{
		isManagerInitialized: true,
		clients: map[string]IRedisClient{
			"test-client": mockClient,
		},
	}
	defer func() { managerInstance = originalManager }()

	// Test the Ping function
	ctx := context.Background()
	err := Ping(ctx, "test-client")
	assert.NoError(t, err)

	// Test with an error from the client
	mockClient.On("Ping", mock.Anything).Return(errors.New("ping error"))
	err = Ping(ctx, "test-client")
	assert.Error(t, err)
	assert.Equal(t, "ping error", err.Error())

	// Verify mock expectations
	mockClient.AssertExpectations(t)
}
