package cachekit

import (
	"context"
	"time"
)

type Cache interface {
	// Save stores a value against a specified key in the store.
	// If a duration is provided, the key-value pair will expire after that duration.
	Save(ctx context.Context, key, value string, ttl time.Duration) error

	// SaveKeepTTL stores a value against a specified key in the store, keeping the existing TTL.
	// Use when updating a key without changing its TTL.
	SaveKeepTTL(ctx context.Context, key, value string) error

	// RemoveKeys removes one or more keys from the store.
	RemoveKeys(ctx context.Context, keys ...string) error

	// Get retrieves the value associated with the given key from the store.
	Get(ctx context.Context, key string) (string, error)

	// IsKeyActive checks if the specified key exists in the store and is active.
	// Returns true if the key exists, false otherwise.
	IsKeyActive(ctx context.Context, key string) (bool, error)

	// Implements Set interface
	Set

	// Ping pings the store to check if it's reachable.
	Ping(ctx context.Context) error

	// Close closes the connection to the store, releasing any associated resources.
	Close()
}

type Set interface {
	// AddSet adds a value to the specified set in the store.
	// Returns the number of values added.
	// Duplicate values are not added again and count does not increase.
	AddSet(ctx context.Context, set string, value ...string) (int64, error)

	// RemoveSetValues removes one or more values from the specified set in the store.
	// Returns the number of values removed.
	RemoveSetValues(ctx context.Context, set string, values ...string) (int64, error)

	// RemoveSet removes the specified set from the store.
	RemoveSet(ctx context.Context, set string) error

	// ContainsSet checks if the specified value exists in the given set in the store.
	// Returns 0 if the value doesn't exist, 1 otherwise.
	ContainsSet(ctx context.Context, set string, value string) (bool, error)
}
