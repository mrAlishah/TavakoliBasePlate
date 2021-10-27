package mysql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DriverName string            `yaml:"DRIVER_NAME"`
	DataSource string            `yaml:"DATA_SOURCE"`
	Tables     map[string]string `yaml:"TABLES"`
}

type Connection struct {
	Db     *sql.DB
	Config *Config
}

func NewConnection(ctx context.Context, config *Config) (*Connection, error) {
	logFields := logrus.Fields{
		"platform": config.DriverName,
		"domain":   "localhost:3306",
	}
	logrus.WithFields(logFields).Info("Connecting to mySQL DB")

	db, err := sql.Open(config.DriverName, config.DataSource)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Connection{
		Db:     db,
		Config: config,
	}, nil
}

func (m *Connection) Close() error {
	return m.Db.Close()
}
