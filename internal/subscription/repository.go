package subscription

import (
	"context"
)

type RedisRepository interface {
	IsReady() bool
}

type ExternalDriver interface {
}

type MySQLRepository interface {
	GetProducts(ctx context.Context, limit int) error
}

type PostgresRepository interface {
	CreateUser(user User) (User, error)
}
