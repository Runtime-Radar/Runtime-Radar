package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/runtime-radar/runtime-radar/lib/security"
)

const (
	// Redis CA cert file name.
	redisCAFile      = "redis_ca.pem"
	defaultRedisPort = "6379"
	minTLSVersion    = tls.VersionTLS12
)

type Remote struct {
	KeyPrefix string
	Redis     *redis.Client
}

func NewRemote(addr, user, password string, tlsMode, tlsCheckCert bool, keyPrefix string) (*Remote, func() error, error) {
	if !strings.Contains(addr, ":") {
		addr = addr + ":" + defaultRedisPort
	}

	redisOptions := &redis.Options{
		Addr:     addr,
		Username: user,
		Password: password,
	}

	if tlsMode {
		caCert, err := os.ReadFile(redisCAFile)
		if err != nil {
			return nil, nil, fmt.Errorf("can't read redis CA file: %w", err)
		}

		caCertPool, err := security.LoadSystemCABundle(string(caCert))
		if err != nil {
			return nil, nil, fmt.Errorf("can't load redis CA bundle: %w", err)
		}
		redisOptions.TLSConfig = &tls.Config{
			MinVersion:         minTLSVersion,
			InsecureSkipVerify: !tlsCheckCert,
			RootCAs:            caCertPool,
		}

	}

	red := redis.NewClient(redisOptions)

	remote := &Remote{
		keyPrefix,
		red,
	}

	return remote, red.Close, test(red)
}

func test(c *redis.Client) error {
	// Key does not matter
	if err := c.Get(context.Background(), "test").Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}

func (r *Remote) Get(ctx context.Context, key string, val any) (bool, error) {
	key = r.KeyPrefix + key

	item, err := r.Redis.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, unmarshal([]byte(item), val)
}

func (r *Remote) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	key = r.KeyPrefix + key

	item, err := marshal(val)
	if err != nil {
		return err
	}

	return r.Redis.Set(ctx, key, item, expiration).Err()
}

func (r *Remote) Del(ctx context.Context, key string) error {
	key = r.KeyPrefix + key

	if err := r.Redis.Del(ctx, key).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}
