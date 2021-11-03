package mysql

import (
	"GolangTraining/platform/mysql"
	"context"
)

type Repository struct {
	database *mysql.Connection
}

func CreateRepository(db *mysql.Connection) (*Repository, error) {
	return &Repository{
		database: db,
	}, nil
}

func (m *Repository) HealthCheck() error {
	err := m.database.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) GetProducts(ctx context.Context, limit int) error {
	return nil
}
