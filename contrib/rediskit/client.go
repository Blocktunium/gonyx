package rediskit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Blocktunium/gonyx/internal/config"
	"github.com/Blocktunium/gonyx/internal/logger"
	"github.com/Blocktunium/gonyx/internal/logger/types"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	redisMaintenanceType = types.NewLogType("REDIS_MAINTENANCE")
)

// Client represents a Redis client with various operations
type Client struct {
	name        string
	prefix      string
	initialized bool
	client      *redis.Client
	wg          sync.WaitGroup
	lock        sync.Mutex
	lockEnable  bool
}

// Init initializes the Redis client with the provided configuration
func (c *Client) Init(name string, configPrefix string, keyPrefix string) error {
	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Init Start", nil))
	}

	c.wg.Add(1)
	defer c.wg.Done()

	c.name = name
	c.prefix = keyPrefix
	c.initialized = false

	addr, err := config.GetManager().Get(name, configPrefix+".address")
	if err != nil {
		return err
	}

	password, err := config.GetManager().Get(name, configPrefix+".password")
	if err != nil {
		return err
	}

	db, err := config.GetManager().Get(name, configPrefix+".db")
	if err != nil {
		return err
	}

	maxRetries, err := config.GetManager().Get(name, configPrefix+".max_retries")
	if err != nil {
		return err
	}

	minRetryBackOff, err := config.GetManager().Get(name, configPrefix+".min_retry_backoff")
	if err != nil {
		return err
	}

	maxRetryBackOff, err := config.GetManager().Get(name, configPrefix+".max_retry_backoff")
	if err != nil {
		return err
	}

	dialTimeout, err := config.GetManager().Get(name, configPrefix+".dial_timeout")
	if err != nil {
		return err
	}

	readTimeout, err := config.GetManager().Get(name, configPrefix+".read_timeout")
	if err != nil {
		return err
	}

	writeTimeout, err := config.GetManager().Get(name, configPrefix+".write_timeout")
	if err != nil {
		return err
	}

	onConnectLog, err := config.GetManager().Get(name, configPrefix+".on_connect_log")
	if err != nil {
		return err
	}

	lockEnable, err := config.GetManager().Get(name, configPrefix+".enable_lock")
	if err != nil {
		return err
	}
	c.lockEnable = lockEnable.(bool)

	options := &redis.Options{
		Addr:            addr.(string),
		Password:        password.(string),
		DB:              int(db.(float64)),
		MaxRetries:      int(maxRetries.(float64)),
		MinRetryBackoff: time.Duration(minRetryBackOff.(float64)) * time.Millisecond,
		MaxRetryBackoff: time.Duration(maxRetryBackOff.(float64)) * time.Millisecond,
		DialTimeout:     time.Duration(dialTimeout.(float64)) * time.Millisecond,
		ReadTimeout:     time.Duration(readTimeout.(float64)) * time.Millisecond,
		WriteTimeout:    time.Duration(writeTimeout.(float64)) * time.Millisecond,
	}

	if onConnectLog.(bool) {
		options.OnConnect = func(ctx context.Context, conn *redis.Conn) error {
			// Log connection events
			return nil
		}
	}

	c.client = redis.NewClient(options)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	err = c.Ping(ctx)
	if err != nil {
		return err
	}

	c.initialized = true

	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Init End", nil))
	}

	return nil
}

// Ping checks if the Redis server is reachable
func (c *Client) Ping(ctx context.Context) error {
	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Ping Start", nil))
	}

	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		return NewPingError(err)
	}

	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Ping End", nil))
	}
	return nil
}

// IsInitialized returns whether the client is initialized
func (c *Client) IsInitialized() bool {
	return c.initialized
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	c.wg.Wait()

	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Redis Connection Close Start", nil))
	}

	err := c.client.Close()
	if err == nil {
		if l != nil {
			l.Log(types.NewLogObject(types.DEBUG, "RedisKit", redisMaintenanceType, time.Now(), "Redis Connection Close End", nil))
		}
	}

	return err
}

// Get retrieves a value from Redis by key
func (c *Client) Get(ctx context.Context, key string, val any) error {
	c.wg.Wait()

	err := c.client.Get(ctx, c.generateKey(key)).Scan(val)
	if err != nil {
		return NewReadError(key, err)
	}
	return nil
}

// Set stores a value in Redis with the given key and expiration
func (c *Client) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.wg.Wait()

	err := c.client.Set(ctx, c.generateKey(key), val, expiration).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}

	return nil
}

// SetStruct marshals and stores a struct in Redis
func (c *Client) SetStruct(ctx context.Context, key string, val any, expiration time.Duration) error {
	c.wg.Wait()

	// Marshal the struct to JSON
	marshalled, err := json.Marshal(val)
	if err != nil {
		return NewWriteError(key, val, err)
	}

	return c.Set(ctx, key, marshalled, expiration)
}

// GetStruct retrieves and unmarshals a struct from Redis
func (c *Client) GetStruct(ctx context.Context, key string, val any) error {
	c.wg.Wait()

	var tempArr []byte
	err := c.Get(ctx, key, &tempArr)
	if err != nil {
		return NewReadError(key, err)
	}

	err = json.Unmarshal(tempArr, val)
	if err != nil {
		return NewReadError(key, err)
	}

	return nil
}

// HSet sets field-value pairs in a Redis hash
func (c *Client) HSet(ctx context.Context, key string, expiration time.Duration, val ...any) error {
	c.wg.Wait()

	err := c.client.HSet(ctx, c.generateKey(key), val...).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}

	err = c.client.Expire(ctx, c.generateKey(key), expiration).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}
	return nil
}

// HGet retrieves a field value from a Redis hash
func (c *Client) HGet(ctx context.Context, key string, field string, val any) error {
	c.wg.Wait()

	cmd := c.client.HGet(ctx, c.generateKey(key), field)
	err := cmd.Scan(val)
	if err != nil {
		return NewReadError(fmt.Sprintf("%s:%s", key, field), err)
	}
	return nil
}

// generateKey creates a prefixed key for Redis
func (c *Client) generateKey(key string) string {
	newKey := key
	if c.prefix != "" {
		newKey = c.prefix + "$" + newKey
	}
	return newKey
}
