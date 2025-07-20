package utils

import (
	"sync"
	"time"

	"hatena-bookmark-mcp/internal/types"
)

// CacheEntry represents a cached item with expiration
type CacheEntry struct {
	Data      *types.GetHatenaBookmarksResponse
	ExpiresAt time.Time
}

// Cache provides simple in-memory caching functionality
type Cache struct {
	entries map[string]*CacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
}

// NewCache creates a new cache instance
func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanup()
	
	return cache
}

// Get retrieves a cached item by key
func (c *Cache) Get(key string) (*types.GetHatenaBookmarksResponse, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}
	
	// Check if entry has expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	return entry.Data, true
}

// Set stores an item in the cache
func (c *Cache) Set(key string, data *types.GetHatenaBookmarksResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.entries[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.entries, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.entries = make(map[string]*CacheEntry)
}

// Size returns the number of items in the cache
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return len(c.entries)
}

// cleanup removes expired entries from the cache
func (c *Cache) cleanup() {
	ticker := time.NewTicker(time.Minute) // Clean up every minute
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()
	}
}

// GenerateCacheKey creates a cache key from bookmark parameters
func GenerateCacheKey(params types.GetHatenaBookmarksParams) string {
	// Create a deterministic key from the parameters
	key := params.Username
	
	if params.Tag != "" {
		key += "_tag:" + params.Tag
	}
	
	if params.Date != "" {
		key += "_date:" + params.Date
	}
	
	if params.URL != "" {
		key += "_url:" + params.URL
	}
	
	if params.Page > 0 {
		key += "_page:" + string(rune(params.Page))
	}
	
	return key
}