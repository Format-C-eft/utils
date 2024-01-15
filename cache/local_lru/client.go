package local_lru

import (
	"context"
	"fmt"
	"time"

	lruExp "github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/hashicorp/golang-lru/v2/simplelru"

	"github.com/Format-C-eft/utils/cache"
)

type localCache struct {
	cache simplelru.LRUCache[string, any]
}

func New(name string, cap int, ttl time.Duration) (cache.Store, error) {
	if cap <= 0 {
		return nil, fmt.Errorf("can't create cache %s: wrong capacity, it should be positive", name)
	}
	if ttl < 0 {
		return nil, fmt.Errorf("can't create cache %s: wrong TTL, it should be >= 0", name)
	}

	return &localCache{
		cache: lruExp.NewLRU[string, any](cap, nil, ttl),
	}, nil
}

func (c *localCache) Get(_ context.Context, key string) (value interface{}, err error) {
	v, ok := c.cache.Get(key)
	if ok {
		return v, nil
	}

	return nil, nil
}

func (c *localCache) Set(_ context.Context, key string, value interface{}) error {
	c.cache.Add(key, value)

	return nil
}

func (c *localCache) Remove(_ context.Context, key string) error {
	c.cache.Remove(key)

	return nil
}

func (c *localCache) MultiGet(_ context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		if v, ok := c.cache.Get(key); ok {
			result[key] = v
		}
	}

	return result, nil
}

func (c *localCache) MultiSet(_ context.Context, data map[string]interface{}) error {
	for key, value := range data {
		c.cache.Add(key, value)
	}

	return nil
}

func (c *localCache) MultiRemove(_ context.Context, keys []string) error {
	for _, key := range keys {
		c.cache.Remove(key)
	}

	return nil
}

func (c *localCache) GetAllKeys(_ context.Context) []string {
	return c.cache.Keys()
}
