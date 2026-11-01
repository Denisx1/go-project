package redisclient

import (
	"context"
	"fmt"
	"go_project/infrastructure/cache"
	"time"

	"github.com/bytedance/sonic"
	redis "github.com/redis/go-redis/v9"
)

type client struct {
	redis       *redis.Client
	opTimeout   time.Duration
	pingTimeout time.Duration
}

func NewClient(ctx context.Context, cfg Config) (cache.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:         cfg.SetAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := c.Ping(pingCtx).Err(); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &client{c, cfg.ReadTimeout, cfg.DialTimeout}, nil
}

func (p *client) withTimeout(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	if p.opTimeout <= 0 {
		return ctx, func() {}
	}
	if deadline, ok := ctx.Deadline(); ok {
		if time.Until(deadline) < p.opTimeout {
			return ctx, func() {}
		}
	}
	return context.WithTimeout(ctx, p.opTimeout)
}

func (c *client) Close() error {
	if err := c.redis.Close(); err != nil {
		return fmt.Errorf("redis error: %w", err)
	}
	return nil
}
func (a *client) Ping(ctx context.Context) error {
	reqCtx, cancel := context.WithTimeout(ctx, a.pingTimeout)
	defer cancel()
	cmd := a.redis.Ping(reqCtx)
	return statusCMD{status: cmd}.Err()
}

func (c *client) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) cache.StatusResult {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	payload := value
	switch v := value.(type) {
	case string, []byte:
		payload = v
	default:
		bytes, err := sonic.Marshal(value)
		if err != nil {
			cmd := redis.NewStatusCmd(reqCtx)
			cmd.SetErr(err)
			return statusCMD{status: cmd}
		}
		payload = bytes
	}
	cmd := c.redis.Set(reqCtx, key, payload, expiration)

	return statusCMD{
		status: cmd,
	}
}

func (c *client) Get(
	ctx context.Context,
	key string,
) cache.StringResult {
	reqCtx, cancel := context.WithTimeout(ctx, c.opTimeout)
	defer cancel()
	cmd := c.redis.Get(reqCtx, key)

	return stringCMD{
		cmd: cmd,
	}
}
func (c *client) Exists(
	ctx context.Context,
	key string,
) cache.IntResult {
	reqCtx, cancel := context.WithTimeout(ctx, c.opTimeout)
	defer cancel()
	cmd := c.redis.Exists(reqCtx, key)
	return intCMD{
		cmd: cmd,
	}
}

func (c *client) Delete(
	ctx context.Context,
	keys ...string,
) cache.IntResult {
	reqCtx, cancel := context.WithTimeout(ctx, c.opTimeout)
	defer cancel()
	cmd := c.redis.Del(reqCtx, keys...)
	return intCMD{
		cmd: cmd,
	}
}

func (r *client) Lock(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) (bool, error) {

	acquired, err := r.redis.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, err
	}

	return acquired, nil
}

// 3. Снятие блокировки (если создание в БД завершилось ошибкой)
func (r *client) Unlock(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}
