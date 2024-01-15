package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"

	"github.com/Format-C-eft/utils/dlock"
	"github.com/Format-C-eft/utils/logger"
)

type dLockRedis struct {
	client *redislock.Client
	cfg    Config
}

func New(
	client redis.Cmdable,
	cfg Config,
) dlock.Locker {
	return &dLockRedis{
		client: redislock.New(client),
		cfg:    cfg,
	}
}

func (l *dLockRedis) Lock(ctx context.Context, key string) (func(), bool) {
	lockKey := fmt.Sprintf(dlock.LockKeyTpl, key)
	lock := l.tryAcquireLock(ctx, lockKey)

	if lock == nil {
		return func() {}, false
	}

	return func() {
		l.releaseLock(ctx, lockKey, lock)
	}, true
}

func (l *dLockRedis) acquireLock(ctx context.Context, key string) *redislock.Lock {
	lock, err := l.client.Obtain(ctx, key, l.cfg.TTL, nil)
	if err != nil && !errors.Is(err, redislock.ErrNotObtained) {
		logger.Error(ctx, fmt.Errorf("RedisLock: add failed with key = %s err: %w", key, err))
	}

	return lock
}

func (l *dLockRedis) releaseLock(ctx context.Context, key string, lock *redislock.Lock) {
	if err := lock.Release(ctx); err != nil {
		logger.Debug(ctx, fmt.Errorf("RedisLock: delete failed with key = %s err: %w", key, err))
	}
}

func (l *dLockRedis) tryAcquireLock(ctx context.Context, key string) *redislock.Lock {
	lCtx, cancel := context.WithTimeout(ctx, l.cfg.AcquireLockTimeout)
	defer cancel()

	for {
		select {
		case <-lCtx.Done():
			return nil

		default:
			if lock := l.acquireLock(ctx, key); lock != nil {
				return lock
			}
			time.Sleep(l.cfg.RetryPeriod)
		}
	}
}
