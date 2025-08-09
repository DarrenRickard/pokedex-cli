package main

import "testing"
import "time"
import "fmt"
import "github.com/darrenrickard/pokedexcli/internal/pokecache"


func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val := cache.Get(c.key)
			if !val.Exists {
				t.Errorf("expected to find key")
				return
			}
			if string(val.Data) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	val := cache.Get("https://example.com")
	if !val.Exists {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	val = cache.Get("https://example.com")
	if val.Exists {
		t.Errorf("expected to not find key")
		return
	}
}


