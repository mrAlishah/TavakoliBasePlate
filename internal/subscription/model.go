package subscription

import (
	"gorm.io/gorm"
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
	Pass  string `form:"pass" binding:"required"`
}

type User struct {
	gorm.Model
	ID        uint64 `gorm:"primary_key"`
	Email     string `gorm:"not null"`
	UserName  string
	FullName  string
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	DeletedAt gorm.DeletedAt
}
