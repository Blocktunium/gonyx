package cache

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
	cacheMaintenanceType = types.NewLogType("CACHE_MAINTENANCE")
)

// Mark: RedisClientCache

// RedisClientCache object
type RedisClientCache struct {
	name        string
	prefix      string
	initialized bool
	client      *redis.Client
	wg          sync.WaitGroup
	lock        sync.Mutex
	lockEnable  bool
}

// MARK: Public functions

// Init - Constructor: It reads the redis client configurations and initialize the connection
func (ins *RedisClientCache) Init(name string, configPrefix string, cachePrefix string) error {
	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Init Start", nil))
	}

	ins.wg.Add(1)
	defer ins.wg.Done()

	ins.name = name
	ins.prefix = cachePrefix
	ins.initialized = false

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
	ins.lockEnable = lockEnable.(bool)

	// TODO: read Others config

	config1 := &redis.Options{
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
		config1.OnConnect = func(ctx context.Context, conn *redis.Conn) error {
			// Log here
			return nil
		}
	}

	ins.client = redis.NewClient(config1)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	err = ins.Ping(ctx)
	if err != nil {
		return err
	}

	ins.initialized = true

	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Init End", nil))
	}

	return nil
}

// Ping - ping redis server
func (ins *RedisClientCache) Ping(ctx context.Context) error {
	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Ping Start", nil))
	}

	//cache.wg.Wait()
	_, err := ins.client.Ping(ctx).Result()
	if err != nil {
		return NewPingError(err)
	}

	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Ping End", nil))
	}
	return nil
}

// IsInitialized receiver - that return boolean value
func (ins *RedisClientCache) IsInitialized() bool {
	return ins.initialized
}

// Close - It closes the connection.
func (ins *RedisClientCache) Close() error {
	ins.wg.Wait()

	l, _ := logger.GetManager().GetLogger()
	if l != nil {
		l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Cache Connection Close Start", nil))
	}

	err := ins.client.Close()
	if err == nil {
		if l != nil {
			l.Log(types.NewLogObject(types.DEBUG, "Cache.Redis", cacheMaintenanceType, time.Now(), "Cache Connection Close End", nil))
		}
	}

	return err
}

// Get - get by key receiver
func (ins *RedisClientCache) Get(ctx context.Context, key string, val any) error {
	ins.wg.Wait()

	err := ins.client.Get(ctx, ins.generateKey(key)).Scan(val)
	if err != nil {
		return NewReadError(key, err)
	}
	return nil
}

// Set - set by key and expiration receiver
func (ins *RedisClientCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	ins.wg.Wait()

	err := ins.client.Set(ctx, ins.generateKey(key), val, expiration).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}

	return nil
}

// SetStruct - set the struct value by key
func (ins *RedisClientCache) SetStruct(ctx context.Context, key string, val any, expiration time.Duration) error {
	ins.wg.Wait()

	// first marshal it
	marshalled, err := json.Marshal(val)
	if err != nil {
		return NewWriteError(key, val, err)
	}

	return ins.Set(ctx, key, marshalled, expiration)

	//
	//
	//err := ins.client.HSet(ctx, ins.generateKey(key), val).Err()
	//if err != nil {
	//	return NewWriteError(key, val, err)
	//}
	//
	//err = ins.client.Expire(ctx, ins.generateKey(key), expiration).Err()
	//if err != nil {
	//	return NewWriteError(key, val, err)
	//}

	//return nil
}

// GetStruct - get the struct value by key
func (ins *RedisClientCache) GetStruct(ctx context.Context, key string, val any) error {
	ins.wg.Wait()

	var tempArr []byte
	err := ins.Get(ctx, key, &tempArr)
	if err != nil {
		return NewReadError(key, err)
	}

	err = json.Unmarshal(tempArr, val)
	if err != nil {
		return NewReadError(key, err)
	}

	return nil
}

// MARK: Private Receivers
func (ins *RedisClientCache) generateKey(key string) string {
	newKey := key
	if ins.prefix != "" {
		newKey = ins.prefix + "$" + newKey
	}
	return newKey
}

func (ins *RedisClientCache) HSet(ctx context.Context, key string, expiration time.Duration, val ...any) error {
	ins.wg.Wait()

	err := ins.client.HSet(ctx, key, val...).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}

	err = ins.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return NewWriteError(key, val, err)
	}
	return nil
}

func (ins *RedisClientCache) HGet(ctx context.Context, key string, field string, val any) error {
	ins.wg.Wait()

	cmd := ins.client.HGet(ctx, key, field)
	err := cmd.Scan(val)
	if err != nil {
		return NewReadError(fmt.Sprintf("%s:%s", key, field), err)
	}
	return nil
}
