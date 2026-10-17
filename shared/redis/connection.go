package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ErrPingToRedis = errors.New("ping to redis error")

func NewClient(ctx context.Context, cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:                  cfg.Addr,
		Password:              cfg.Password,
		DB:                    cfg.DB,
		MaxRetries:            cfg.MaxRetries,
		DialTimeout:           cfg.DialTimeout,
		ReadTimeout:           cfg.Timeout,
		WriteTimeout:          cfg.Timeout,
		ContextTimeoutEnabled: true,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		pingErr := fmt.Errorf("%w: %w", ErrPingToRedis, err)
		if closeErr := client.Close(); closeErr != nil {
			return nil, errors.Join(pingErr, fmt.Errorf("close redis client: %w", closeErr))
		}
		return nil, pingErr
	}

	return client, nil
}
