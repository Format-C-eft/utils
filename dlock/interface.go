package dlock

import (
	"context"
)

const LockKeyTpl = "lock:%s"

type Locker interface {
	Lock(ctx context.Context, key string) (teardown func(), acquired bool)
}
