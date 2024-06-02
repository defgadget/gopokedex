package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries  map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.Entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.Entries[key]
	c.mu.Unlock()
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.Entries {
				if time.Since(entry.createdAt) >= c.interval {
					delete(c.Entries, key)
				}
			}
			c.mu.Unlock()
		}
	}

}
