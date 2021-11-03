package subscription

import (
	"GolangTraining/internal/logger"
	"GolangTraining/internal/metrics"
	"context"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Test(req Request) (string, error)
	Create(ctx context.Context, req Request) (User, error)
}

type service struct {
	validate   *validator.Validate
	mysql      MySQLRepository
	postgres   PostgresRepository
	redis      RedisRepository
	logger     *logger.StandardLogger
	prometheus *metrics.Prometheus
	config     *Config
}

func CreateService(
	config *Config,
	logger *logger.StandardLogger,
	mysql MySQLRepository,
	postgres PostgresRepository,
	redis RedisRepository,
	prometheus *metrics.Prometheus,
	validator *validator.Validate) Service {
	return &service{
		validate:   validator,
		redis:      redis,
		mysql:      mysql,
		postgres:   postgres,
		logger:     logger,
		prometheus: prometheus,
		config:     config,
	}
}
