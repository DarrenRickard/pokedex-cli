package pokecache

import (
	"fmt"
	"sync"
	"time"
)

var mux = &sync.Mutex{}

type CacheEntryData struct {
	Data     []byte
	exists   bool
}

type cacheEntry struct {
	createdAt		time.Time
	Val				[]byte
}

type Cache struct {
	Entries			map[string]cacheEntry
	interval		time.Duration
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
	}
	// cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte ) error {
	// !!! need to check if key exists in cache, if exists, just return nil
	c1 := make(chan CacheEntryData)
	go func() {
		c1 <- c.Get(key)
	}()
	cachedData := <- c1
	if cachedData.exists {
		return nil
	}
	if c.Entries == nil {
		return fmt.Errorf("Error: cacheEntry map not initialized in Cache struct")
	}
	v := newCacheEntry(val)
	mux.Lock()
	c.Entries[key] = v 
	mux.Unlock()
	return nil
}

func (c *Cache) Get(key string) CacheEntryData {
	mux.Lock()
	entry, ok := c.Entries[key] 
	mux.Unlock()
	if !ok {
		return CacheEntryData{}
	}
	cacheData := CacheEntryData {
		Data: entry.Val,
		exists: true,
	}
	return cacheData
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
		for key, entry := range c.Entries {
			if t.Sub(entry.createdAt) > c.interval {
				delete(c.Entries, key) 
			} 
		}
		mux.Unlock()
	}
}
