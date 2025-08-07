package pokecache

import (
	"fmt"
	"sync"
	"time"
)

var mux = &sync.Mutex{}

type cacheEntry struct {
	createdAt		time.Time
	val				[]byte
}

type Cache struct {
	entries			map[string]cacheEntry
	interval		time.Duration
}

func newCacheEntry(v []byte) cacheEntry {
	var c = cacheEntry {
		createdAt: 	time.Now(),
		val: 		v,	
	}	
	return c
}

func NewCache(i time.Duration) Cache {
	var cache = Cache {
		entries: make(map[string]cacheEntry),
		interval: i,
	}
	return cache
}

func (c *Cache) Add(key string, val []byte ) error {
	if c.entries == nil {
		return fmt.Errorf("Error: cacheEntry map not initialized in Cache struct")
	}
	v := newCacheEntry(val)
	mux.Lock()
	c.entries[key] = v 
	mux.Unlock()
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	mux.Lock()
	entry, ok := c.entries[key] 
	mux.Unlock()
	if !ok {
		return nil, false	
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	// every Cache.interval, call reapLoop()
	// remove all cacheEntry created more than Cache.interval time ago
	// if time.Now().Sub(newCacheEntry.createdAt) > 5 * time.Second {delete the newCacheEntry}
	// look into using time.Ticker
	// should use NewTicker(d duration) *Ticker
	// returns a new Ticker containing a channel that will send the current time on the channel after each tick
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	t := <-ticker.C 
	for {
		mux.Lock()
		for key, entry := range c.entries {
			if t.Sub(entry.createdAt) > c.interval {
				delete(c.entries, key) 
			} 
		}
		mux.Unlock()
	}
}
