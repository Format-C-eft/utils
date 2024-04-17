package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Format-C-eft/utils/cache"
	"github.com/Format-C-eft/utils/cache/codec"
)

type redisClient struct {
	cache     redis.Cmdable
	prefixKey string
	ttl       time.Duration
	codec     codec.Codec
}

func NewClient(cache redis.Cmdable, prefixKey string, ttl time.Duration, opts ...Options) cache.Store {
	settings := cacheOptions{}

	for _, o := range opts {
		o(&settings)
	}

	return &redisClient{
		cache:     cache,
		prefixKey: prefixKey,
		ttl:       ttl,
		codec:     settings.Codec,
	}
}

func (c *redisClient) Get(ctx context.Context, key string) (value interface{}, err error) {
	data, err := c.cache.Get(ctx, c.getKey(key)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis.Get err: %w", err)
	}

	if c.codec != nil {
		val, errCodec := c.codec.Unmarshal(data)
		if errCodec != nil {
			return nil, fmt.Errorf("redis.Get:codec.Marshal err: %w", errCodec)
		}
		return val, nil
	}

	return data, nil
}

func (c *redisClient) Set(ctx context.Context, key string, value interface{}) error {
	if c.codec != nil {
		val, errMarshal := c.codec.Marshal(value)
		if errMarshal != nil {
			return fmt.Errorf("redis.Set:codec.Marshal err: %w", errMarshal)
		}
		value = val
	}

	if errSet := c.cache.Set(ctx, c.getKey(key), value, c.ttl).Err(); errSet != nil {
		return fmt.Errorf(
			"redis.Set key: %q value: %v valueType: %T codecType: %T err: %w",
			key,
			value,
			value,
			c.codec,
			errSet,
		)
	}

	return nil
}

func (c *redisClient) Remove(ctx context.Context, key string) error {
	err := c.cache.Del(ctx, c.getKey(key)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("redis.Del err: %w", err)
	}

	return nil
}

func (c *redisClient) MultiGet(ctx context.Context, keys []string) (map[string]interface{}, error) {
	result := make(map[string]interface{}, len(keys))
	if len(keys) == 0 {
		return result, nil
	}

	getKeys := make([]string, len(keys))
	for i, key := range keys {
		getKeys[i] = c.getKey(key)
	}

	values, err := c.cache.MGet(ctx, getKeys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return result, nil
		}
		return nil, fmt.Errorf("redis.MultiGet err: %w", err)
	}

	for i, v := range values {
		if v == nil {
			continue
		}

		if c.codec != nil {
			valUnMarshal, errMarshal := c.codec.Unmarshal([]byte(v.(string)))
			if errMarshal != nil {
				return nil, fmt.Errorf("redis.MultiGet:codec.Marshal err: %w", err)
			}
			v = valUnMarshal
		}

		result[keys[i]] = v
	}

	return result, nil
}

func (c *redisClient) MultiSet(ctx context.Context, data map[string]interface{}) error {
	preparedData := make(map[string]interface{}, len(data))

	if c.codec != nil {
		for key, value := range data {
			val, err := c.codec.Marshal(value)
			if err != nil {
				return fmt.Errorf("redis.MultiSet:codec.Marshal err: %w", err)
			}
			preparedData[c.getKey(key)] = val
		}

		if errSet := c.cache.MSet(ctx, preparedData).Err(); errSet != nil {
			return fmt.Errorf("redis.MultiSet err: %w", errSet)
		}

		return nil
	}

	for key, value := range data {
		preparedData[c.getKey(key)] = value
	}

	if errSet := c.cache.MSet(ctx, preparedData).Err(); errSet != nil {
		return fmt.Errorf("redis.MultiSet err: %w", errSet)
	}

	return nil
}

func (c *redisClient) MultiRemove(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	removeKeys := make([]string, len(keys))
	for i, key := range keys {
		removeKeys[i] = c.getKey(key)
	}

	err := c.cache.Del(ctx, removeKeys...).Err()
	if err != nil {
		return fmt.Errorf("redis.MultiRemove err: %w", err)
	}

	return nil
}

func (c *redisClient) GetAllKeys(ctx context.Context) []string {
	result := c.cache.Keys(ctx, c.getKey("*")).Val()
	for i, s := range result {
		result[i] = strings.Replace(s, c.prefixKey+":", "", 1)
	}

	return result
}

func (c *redisClient) getKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefixKey, key)
}
