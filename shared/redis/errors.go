package redis

import "errors"

var (
	ErrPingToRedis       = errors.New("redis failed to connect")
	ErrInvalidPingConfig = errors.New("invalid ping config")
)
