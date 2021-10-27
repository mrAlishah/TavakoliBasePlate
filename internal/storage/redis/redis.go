package redis

import (
	"GolangTraining/platform/redis"
	"github.com/pkg/errors"
)

type Repository struct {
	database *redis.Connection
}

func CreateRepository(db *redis.Connection) (*Repository, error) {
	if db != nil {
		return &Repository{
			database: db,
		}, nil
	} else {
		return &Repository{
			database: db,
		}, errors.New("redis initialization error")
	}

}

func (r *Repository) IsReady() bool {
	if r.database == nil {
		return false
	}
	return r.database.Ready
}
