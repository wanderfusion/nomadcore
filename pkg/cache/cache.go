// Package cache provides a generic caching mechanism.
package cache

import (
	"time"
)

// Cache is a generic cache structure. The cache is a map from string to Value.
type Cache[T any] struct {
	cache map[string]Value[T]
}

// Value is a structure that holds the value and its Time-To-Live (TTL).
type Value[T any] struct {
	TTL      time.Duration // Time duration the value is valid for
	Value    T             // The actual value
	Inserted time.Time     // Time the value was inserted into the cache
}

// Option is a generic structure for optional values.
type Option[T any] struct {
	Some T    // Holds the value if present
	None bool // Indicates if the value is absent
}

// New creates a new Cache and returns its pointer.
func New[T any]() *Cache[T] {
	cache := make(map[string]Value[T])
	return &Cache[T]{
		cache: cache,
	}
}

// SetValue sets a value in the cache with a given TTL.
func (c *Cache[T]) SetValue(key string, value T, ttl time.Duration) {
	c.cache[key] = Value[T]{
		TTL:      ttl,
		Value:    value,
		Inserted: time.Now(),
	}
}

// GetValue retrieves a value from the cache.
// Returns an Option containing the value if found and not expired.
// Returns None if the value is not found or has expired.
func (c *Cache[T]) GetValue(key string) Option[T] {
	cacheValue, ok := c.cache[key]
	if !ok {
		return Option[T]{None: true}
	}

	if time.Since(cacheValue.Inserted) > cacheValue.TTL {
		delete(c.cache, key)
		return Option[T]{None: true}
	}

	return Option[T]{Some: cacheValue.Value}
}
