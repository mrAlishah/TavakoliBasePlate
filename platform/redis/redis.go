package redis

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host                 string        `yaml:"HOST"`
	Pass                 string        `yaml:"PASS"`
	Db                   int           `yaml:"DB"`
	TimeoutDuration      time.Duration `yaml:"TIMEOUT_DURATION"`
	Key                  string        `yaml:"KEY"`
}

type Connection struct {
	Client *redis.Client
	Ready  bool
	Cfg    *Config
}

func CreateConnection(cfg *Config) (*Connection, error) {
	logFields := logrus.Fields{
		"platform": "redis",
		"domain":   cfg.Host,
	}
	logrus.WithFields(logFields).Info("Connecting to Redis DB")
	options := &redis.Options{
		Addr:            cfg.Host,
		Password:        cfg.Pass,
		DB:              cfg.Db,
		MaxRetries:      3,
		MaxRetryBackoff: 10 * time.Second,
		ReadTimeout:     cfg.TimeoutDuration * time.Millisecond,
		DialTimeout:     cfg.TimeoutDuration * time.Millisecond,
	}

	client := redis.NewClient(options)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Connection{
		Client: client,
		Ready:  true,
		Cfg:    cfg,
	}, nil

}

// Close redis connection
func (r *Connection) Close() {
	r.Ready = false
	err := r.Client.Close()
	if err != nil {
		logrus.Errorf("Redis could not disconnect : %v", err)
	}
}

func (r *Connection) GetInfo() (string, error) {
	return r.Client.Info().Result()
}
