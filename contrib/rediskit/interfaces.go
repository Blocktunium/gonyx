package rediskit

import (
	"context"
	"time"
)

// IRedisClient interface defines the operations available on a Redis client
type IRedisClient interface {
	Ping(ctx context.Context) error
	Init(name string, configPrefix string, cachePrefix string) error
	IsInitialized() bool
	Close() error
	Get(ctx context.Context, key string, val any) error
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	SetStruct(ctx context.Context, key string, val any, expiration time.Duration) error
	GetStruct(ctx context.Context, key string, val any) error
	HSet(ctx context.Context, key string, expiration time.Duration, val ...any) error
	HGet(ctx context.Context, key string, field string, val any) error
}
