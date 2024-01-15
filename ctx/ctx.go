package ctx

import (
	"context"
	"time"
)

// DetachContext returns a context that keeps all the values of its parent context
// but detaches from the cancellation and error handling.
// Link: https://github.com/golang/tools/blob/master/internal/xcontext/xcontext.go
// Alternative link : https://github.com/elastic/apm-agent-go/blob/main/gocontext.go
func DetachContext(ctx context.Context) context.Context { return detachedContext{ctx} }

type detachedContext struct{ parent context.Context }

func (v detachedContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (v detachedContext) Done() <-chan struct{}             { return nil }
func (v detachedContext) Err() error                        { return nil }
func (v detachedContext) Value(key interface{}) interface{} { return v.parent.Value(key) }
