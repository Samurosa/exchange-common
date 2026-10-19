package redis

import "time"

type Config struct {
	Addr            string        `yaml:"addr"`
	Password        string        `yaml:"password"`
	DB              int           `yaml:"db"`
	PingMaxAttempts int           `yaml:"ping_max_attempts"`
	PingRetryDelay  time.Duration `yaml:"ping_retry_delay"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
	Timeout         time.Duration `yaml:"timeout"`
}
