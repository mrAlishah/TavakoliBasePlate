package subscription

import (
	"time"
)

type Config struct {
	Service        string        `yaml:"SERVICE"`
	DefaultTTL     int64         `yaml:"DEFAULT_TTL"`
	CacheDisabled  bool          `yaml:"CACHE_DISABLED"`
	SleepTime      time.Duration `yaml:"SLEEP_TIME"`
	FetchBatchSize int           `yaml:"FETCH_BATCH_SIZE"`
}

type Request struct {
	Email string `form:"email" binding:"required"`
	Pass string `form:"pass" binding:"required"`
}

