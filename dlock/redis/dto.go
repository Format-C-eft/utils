package redis

import (
	"time"
)

type Config struct {
	TTL                time.Duration
	AcquireLockTimeout time.Duration
	RetryPeriod        time.Duration
}
