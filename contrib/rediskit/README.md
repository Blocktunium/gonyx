# RedisKit

RedisKit is a contributed package for the Gonyx framework that provides a clean, easy-to-use Redis client implementation.

## Features

- Singleton Manager to control multiple Redis instances
- Simple Redis client initialization from configuration
- Support for key prefixing to avoid key collisions
- Methods for setting and getting simple values and structs
- Support for Redis hash operations
- Clear error handling with type-specific error messages

## Usage

### Using the Redis Manager

```go
import (
    "github.com/Blocktunium/gonyx/contrib/rediskit"
    "context"
    "time"
)

// Get the RedisKit manager
manager := rediskit.GetManager()

// Get a configured Redis client by name
client, err := manager.GetClient("main")
if err != nil {
    // Handle error
}

// When your application is shutting down
err = manager.Release()
if err != nil {
    // Handle error
}
```

### Using an individual client directly

```go
import (
    "github.com/Blocktunium/gonyx/contrib/rediskit"
    "context"
    "time"
)

// Create a new client
client := rediskit.NewClient()

// Initialize the client with configuration
err := client.Init("configName", "redis", "app-prefix")
if err != nil {
    // Handle error
}

// Check if the client is initialized
if client.IsInitialized() {
    // Use the client
}

// Don't forget to close the client when done
defer client.Close()
```

### Basic operations

```go
// Set a string value with expiration
ctx := context.Background()
err := client.Set(ctx, "my-key", "my-value", 1*time.Hour)
if err != nil {
    // Handle error
}

// Get a string value
var value string
err = client.Get(ctx, "my-key", &value)
if err != nil {
    // Handle error
}

// Store a struct
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

user := User{ID: 1, Name: "John Doe"}
err = client.SetStruct(ctx, "user:1", user, 24*time.Hour)
if err != nil {
    // Handle error
}

// Retrieve a struct
var retrievedUser User
err = client.GetStruct(ctx, "user:1", &retrievedUser)
if err != nil {
    // Handle error
}

// Working with hashes
err = client.HSet(ctx, "user:hash:1", 24*time.Hour, "name", "John", "age", 30)
if err != nil {
    // Handle error
}

var name string
err = client.HGet(ctx, "user:hash:1", "name", &name)
if err != nil {
    // Handle error
}
```

## Configuration

RedisKit requires the following configuration structure when using the Manager:

```json
{
  "connections": ["localhost", "computeServer"],
  "localhost": {
    "type": "redis",
    "redis_type": "client",
    "add_service_prefix": false,
    "client": {
      "enable_lock": true,
      "address": "127.0.0.1:6579",
      "password": "",
      "db": 0,
      "max_retries": 0,
      "min_retry_backoff": 8,
      "max_retry_backoff": 512,
      "dial_timeout": 5000,
      "read_timeout": 3000,
      "write_timeout": 3000,
      "pool_size_per_cpu": 10,
      "min_idle_conn": 1,
      "max_conn_age": -1,
      "pool_timeout": 4000,
      "idle_timeout": 5000,
      "idle_check_frequency": 1000,
      "on_connect_log": true
    }
  },
  "computeServer": {
    "type": "redis",
    "redis_type": "client",
    "add_service_prefix": false,
    "client": {
      "enable_lock": true,
      "address": "127.0.0.1:6379",
      "password": "",
      "db": 0,
      "max_retries": 0,
      "min_retry_backoff": 8,
      "max_retry_backoff": 512,
      "dial_timeout": 5000,
      "read_timeout": 3000,
      "write_timeout": 3000,
      "pool_size_per_cpu": 10,
      "min_idle_conn": 1,
      "max_conn_age": -1,
      "pool_timeout": 4000,
      "idle_timeout": 5000,
      "idle_check_frequency": 1000,
      "on_connect_log": true
    }
  }
}
```

All time values are in milliseconds.
