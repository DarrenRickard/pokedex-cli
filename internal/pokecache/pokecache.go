package pokecache

import (
	"fmt"
	"time"
)

type cacheEntry struct {
	createdAt		time.Time
	val				[]byte
}

type Cache struct {
	entries			map[string]cacheEntry
}

func NewCache(v []byte) cacheEntry {
	var c = cacheEntry {
		createdAt: 	time.Now(),
		val: 		v,	
	}	
	return c
}

func (c *Cache) Add(key string, val []byte ) error {
	if c.entries == nil {
		return fmt.Errorf("Error: cacheEntry map not initialized in Cache struct")
	}
	v := NewCache(val)
	c.entries[key] = v 
	return nil
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key] 
	if !ok {
		return nil, false	
	}
	return entry.val, true
}
