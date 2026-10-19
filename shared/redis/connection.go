package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewClient(ctx context.Context, logger *zap.Logger, cfg Config) (*redis.Client, error) {
	log := logger.Named("redis-connection")

	if cfg.PingMaxAttempts <= 0 {
		return nil, fmt.Errorf(
			"invalid ping max attempts: %d: %w",
			cfg.PingMaxAttempts,
			ErrInvalidPingConfig,
		)
	}

	client := redis.NewClient(&redis.Options{
		Addr:                  cfg.Addr,
		Password:              cfg.Password,
		DB:                    cfg.DB,
		DialTimeout:           cfg.DialTimeout,
		ReadTimeout:           cfg.Timeout,
		WriteTimeout:          cfg.Timeout,
		ContextTimeoutEnabled: true,
	})

	log.Info(
		"redis retry configuration",
		zap.Int("ping_max_attempts", cfg.PingMaxAttempts),
		zap.Duration("ping_retry_delay", cfg.PingRetryDelay),
	)

	var lastErr error
	attempts := 0

	for attempt := 1; attempt <= cfg.PingMaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			lastErr = err
			break
		}

		attempts = attempt
		lastErr = client.Ping(ctx).Err()

		if lastErr == nil {
			log.Info(
				"redis connected",
				zap.Int("attempt", attempt),
				zap.String("address", cfg.Addr),
			)

			return client, nil
		}

		log.Warn(
			"redis ping failed",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", cfg.PingMaxAttempts),
			zap.Error(lastErr),
		)

		if attempt == cfg.PingMaxAttempts {
			break
		}

		timer := time.NewTimer(cfg.PingRetryDelay)

		select {
		case <-ctx.Done():
			lastErr = ctx.Err()
		case <-timer.C:
		}

		timer.Stop()

		if ctx.Err() != nil {
			lastErr = ctx.Err()
			break
		}
	}

	pingErr := fmt.Errorf(
		"%w after %d attempts: %w",
		ErrPingToRedis,
		attempts,
		lastErr,
	)

	if closeErr := client.Close(); closeErr != nil {
		return nil, errors.Join(
			pingErr,
			fmt.Errorf("close redis client: %w", closeErr),
		)
	}

	return nil, pingErr
}
