package cache

import (
	"context"
	"sync"
	"time"
)

type Item struct {
	Value      interface{}
	Expiration int64
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type CacheManager struct {
	mu          sync.RWMutex
	items       map[string]Item
	cleanupFreq time.Duration
}

func NewCacheManager(cleanupFreq time.Duration) *CacheManager {
	return &CacheManager{
		items:       make(map[string]Item),
		cleanupFreq: cleanupFreq,
	}
}

// StartJanitor launches a background goroutine to clean expired items automatically
func (cm *CacheManager) StartJanitor(ctx context.Context) {
	ticker := time.NewTicker(cm.cleanupFreq)
	go func() {
		for {
			select {
			case <-ticker.C:
				cm.evictExpired()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (cm *CacheManager) evictExpired() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	now := time.Now().UnixNano()
	for key, item := range cm.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(cm.items, key)
		}
	}
}

// Flush clears all cached items manually during system maintenance
func (cm *CacheManager) Flush() {
	cm.mu.Lock()
	cm.items = make(map[string]Item)
	cm.mu.Unlock()
}

func (c *CacheManager) Set(key string, value interface{}, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	c.mu.Lock()
	c.items[key] = Item{
		Value:      value,
		Expiration: exp,
	}
	c.mu.Unlock()
}

func (c *CacheManager) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()
	if !found || item.Expired() {
		return nil, false
	}
	return item.Value, true
}

func (c *CacheManager) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
