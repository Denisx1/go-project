package cache

import (
	"context"
	"time"
)

type StatusResult interface {
	Result() (string, error)
	Err() error
}
type IntResult interface {
	Result() (int64, error)
	Uint64() (uint64, error)
	Err() error
}
type StringResult interface {
	Result() (string, error)
	Bytes() ([]byte, error)
	Scan(dest any) error
	Err() error
}

type Cache interface {
	Get(ctx context.Context, key string) StringResult
	Set(ctx context.Context, key string, value any, expiration time.Duration) StatusResult
	Delete(ctx context.Context, keys ...string) IntResult
	Exists(ctx context.Context, key string) IntResult
	Unlock(ctx context.Context, key string) error
	Lock(
		ctx context.Context,
		key string,
		value interface{},
		ttl time.Duration,
	) (bool, error)
}

type Lifecycle interface {
	Ping(ctx context.Context) error
	Close() error
}
type Client interface {
	Cache
	Lifecycle
}
