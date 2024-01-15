package cache

import "context"

type Store interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
	Remove(ctx context.Context, key string) error
	MultiRemove(ctx context.Context, keys []string) error
	MultiSet(ctx context.Context, data map[string]interface{}) error
	MultiGet(ctx context.Context, keys []string) (map[string]interface{}, error)
	GetAllKeys(ctx context.Context) []string
}
