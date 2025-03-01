package cache

import (
	"context"
	"github.com/Blocktunium/gonyx/internal/cache"
	"time"
)

// SetIntoCache - set simple type value into cache by key
func SetIntoCache(ctx context.Context, cacheInstanceName string, key string, val any, expiration time.Duration) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	err = cacheInstance.Set(ctx, key, val, expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetFromCache - get simple type value from the cache by key
func GetFromCache(ctx context.Context, cacheInstanceName string, key string, val any) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	err = cacheInstance.Get(ctx, key, val)
	if err != nil {
		return err
	}

	return nil
}

// SetHashmapIntoCache - set map/struct/array value into cache by key
func SetHashmapIntoCache(ctx context.Context, cacheInstanceName string, key string, val any, expiration time.Duration) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	err = cacheInstance.SetStruct(ctx, key, val, expiration)
	if err != nil {
		return err
	}

	return nil
}

// GetHashmapFromCache - get map/struct/array value from the cache by key
func GetHashmapFromCache(ctx context.Context, cacheInstanceName string, key string, val any) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	err = cacheInstance.GetStruct(ctx, key, val)
	if err != nil {
		return err
	}

	return nil
}

// HSetIntoCache - set struct value into the cache by key
func HSetIntoCache(ctx context.Context, cacheInstanceName string, key string, expiration time.Duration, val ...any) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	return cacheInstance.HSet(ctx, key, expiration, val...)
}

func HGetFromCache(ctx context.Context, cacheInstanceName string, key string, field string, val any) error {
	cacheInstance, err := cache.GetManager().GetCache(cacheInstanceName)
	if err != nil {
		return err
	}

	return cacheInstance.HGet(ctx, key, field, val)
}

// Release - release the cache connection
func Release() error {
	return cache.GetManager().Release()
}
