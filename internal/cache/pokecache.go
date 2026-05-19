// Handles caching of recent search results
package pokecache

import (
	"sync"
	"time"
)

// Entry for a cached page result
type cacheEntry struct {
	createdAt  time.Time
	val        []byte
}

// Map of cached search results
type Cache struct {
	memory map[string]cacheEntry
	mutex sync.Mutex
}

// Creates a new cache instance
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		memory: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

// Saves a search result page
func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.memory[key] = cacheEntry{
		createdAt:  time.Now(),
		val: val,
	}
}

// Looks for url contents in cache
func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, exists := cache.memory[key]
	return entry.val, exists
}

// Removes old cached results based on the received interval
func (cache *Cache) reapLoop(interval time.Duration){
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cache.mutex.Lock()
		for key, entry := range cache.memory {	
			if time.Since(entry.createdAt) > interval {
				delete(cache.memory, key)
			}
		}
		cache.mutex.Unlock()
	}
} 
