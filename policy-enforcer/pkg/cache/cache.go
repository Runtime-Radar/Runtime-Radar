package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, val any) (bool, error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}
