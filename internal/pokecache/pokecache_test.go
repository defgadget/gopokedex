package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestCacheAddAndRetrieve(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Fatalf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Fatalf("got %s expected %s", val, c.val)
				return
			}

		})
	}
}

func TestCacheReepLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("test", []byte("testdata"))
	if _, ok := cache.Get("test"); !ok {
		t.Fatalf("expected to find key")
	}

	time.Sleep(waitTime)

	cache.Get("test")
	if _, ok := cache.Get("test"); ok {
		t.Fatalf("cache entry found should have been 'reeped'")
	}

}
