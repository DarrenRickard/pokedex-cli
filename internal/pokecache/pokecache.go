package pokecache

import (
	"fmt"
	"sync"
	"time"
)


type CacheEntryData struct {
	Data     []byte
	Exists   bool
}

type cacheEntry struct {
	createdAt		time.Time
	Val				[]byte
}

type Cache struct {
	Entries			map[string]cacheEntry
	interval		time.Duration
	mux 			*sync.Mutex
}

func newCacheEntry(v []byte) cacheEntry {
	var c = cacheEntry {
		createdAt: 	time.Now(),
		Val: 		v,	
	}	
	return c
}

func NewCache(i time.Duration) Cache {
	var cache = Cache {
		Entries: make(map[string]cacheEntry),
		interval: i,
		mux: 		&sync.Mutex{},
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte ) error {
	c1 := make(chan CacheEntryData)
	go func() {
		c1 <- c.Get(key)
	}()
	cachedData := <- c1
	if cachedData.Exists {
		return nil
	}
	if c.Entries == nil {
		return fmt.Errorf("Error: cacheEntry map not initialized in Cache struct")
	}
	v := newCacheEntry(val)
	c.mux.Lock()
	c.Entries[key] = v 
	c.mux.Unlock()
	return nil
}

func (c *Cache) Get(key string) CacheEntryData {
	c.mux.Lock()
	entry, ok := c.Entries[key] 
	c.mux.Unlock()
	if !ok {
		return CacheEntryData{}
	}
	cacheData := CacheEntryData {
		Data: entry.Val,
		Exists: true,
	}
	return cacheData
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.mux.Lock()
		for key, entry := range c.Entries {
			if t.Sub(entry.createdAt) > c.interval {
				delete(c.Entries, key)
			}
		}
		c.mux.Unlock()
	}
}

