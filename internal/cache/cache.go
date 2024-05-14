package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const TTL = time.Second * 2
const GC_INTERVAL = time.Second * 10

type Cache struct {
	mu         sync.RWMutex
	data       map[string]*CacheItem
	ttl        time.Duration
	gcInterval time.Duration
}

type CacheItem struct {
	value     []byte
	expiredAt time.Time
}

func New(ttl ...time.Duration) *Cache {
	actualTTL := TTL

	if len(ttl) > 0 {
		actualTTL = ttl[0]
	}
	cache := &Cache{
		mu:         sync.RWMutex{},
		data:       make(map[string]*CacheItem),
		ttl:        actualTTL,
		gcInterval: GC_INTERVAL, //TODO choose CG interval
	}
	go cache.gc() //start GC
	return cache
}

func (c *Cache) Set(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// TODO some protection from cleanUpCase ?
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.data[key] = &CacheItem{
		value:     bytes,
		expiredAt: time.Now().Add(c.ttl),
	}
	return nil
}

func (c *Cache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if ok {
		if time.Now().After(item.expiredAt) {
			return nil, fmt.Errorf("key not found")
		}
		var value any
		err := json.Unmarshal(item.value, &value)
		if err != nil {
			return nil, err
		}

		item.expiredAt = time.Now().Add(c.ttl) // reset TTL
		return value, nil
	}

	return nil, fmt.Errorf("key not found")
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) gc() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		<-ticker.C
		c.removeExpiredItems()
	}
}

func (c *Cache) removeExpiredItems() {
	//Collect keys of outdated items
	c.mu.RLock()
	var gcList []string
	for key, item := range c.data {
		if time.Now().After(item.expiredAt) {
			gcList = append(gcList, key)
		}
	}
	c.mu.RUnlock()

	//Delete outdated values
	c.mu.Lock()
	for _, key := range gcList {
		delete(c.data, key)
	}
	c.mu.Unlock()
}
