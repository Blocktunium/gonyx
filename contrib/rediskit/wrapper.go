package rediskit

import (
	"context"
	"time"
)

// Set stores a simple value in Redis with the specified client instance, key, and expiration
func Set(ctx context.Context, clientName string, key string, val any, expiration time.Duration) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, val, expiration)
}

// Get retrieves a simple value from Redis with the specified client instance and key
func Get(ctx context.Context, clientName string, key string, val any) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.Get(ctx, key, val)
}

// SetStruct marshals and stores a struct in Redis with the specified client instance, key, and expiration
func SetStruct(ctx context.Context, clientName string, key string, val any, expiration time.Duration) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.SetStruct(ctx, key, val, expiration)
}

// GetStruct retrieves and unmarshals a struct from Redis with the specified client instance and key
func GetStruct(ctx context.Context, clientName string, key string, val any) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.GetStruct(ctx, key, val)
}

// HSet sets field-value pairs in a Redis hash
func HSet(ctx context.Context, clientName string, key string, expiration time.Duration, val ...any) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.HSet(ctx, key, expiration, val...)
}

// HGet retrieves a field value from a Redis hash
func HGet(ctx context.Context, clientName string, key string, field string, val any) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.HGet(ctx, key, field, val)
}

// Ping checks if the Redis server for the specified client instance is reachable
func Ping(ctx context.Context, clientName string) error {
	client, err := GetManager().GetClient(clientName)
	if err != nil {
		return err
	}

	return client.Ping(ctx)
}
