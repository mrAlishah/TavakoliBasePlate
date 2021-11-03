package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var stringType = "type"
var stringConnection = "connection"

type Config struct {
	HOST     string `yaml:"HOST"`
	PORT     int    `yaml:"PORT"`
	NAME     string `yaml:"NAME"`
	USER     string `yaml:"USER"`
	PASSWORD string `yaml:"PASSWORD"`
	SSLMODE  string `yaml:"SSLMODE"`
	DEBUG    bool   `yaml:"DEBUG"`
}

type Connections interface {
	OpenSQLX() (*sqlx.DB, error)
	OpenGORM() (*gorm.DB, error)
}

type databaseConfig struct {
	config Config
	domain string
}

func (dc *databaseConfig) connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s  sslmode=%s",
		dc.config.HOST, dc.config.PORT, dc.config.USER, dc.config.PASSWORD, dc.config.NAME, dc.config.SSLMODE)
}

// CreateConnection Creates a connection to Postgres database using config and domain name
func CreateConnection(config Config, domain string) Connections {
	return &databaseConfig{
		config: config,
		domain: domain,
	}
}

// OpenSQLX creates a new *sqlx.DB Database instance
func (dc *databaseConfig) OpenSQLX() (*sqlx.DB, error) {
	// log fields for logrus, no need to write this multiple times
	logFields := logrus.Fields{
		"platform": "postgres",
		"domain":   dc.domain,
	}
	logConnectionStringFields := logrus.Fields{
		stringType:       "main",
		stringConnection: dc.connectionString(),
	}

	logrus.WithFields(logFields).Info("Connecting to Postgres SQL DB [RAW]")

	logrus.WithFields(logFields).Info("Opening Connection to Main")
	dbMaster, err := sqlx.Open("postgres", dc.connectionString())
	if err != nil {
		logrus.WithFields(logConnectionStringFields).Fatal(err)
		return nil, err
	}
	err = dbMaster.Ping()
	if err != nil {
		logrus.WithFields(logConnectionStringFields).Fatal(err)
		return nil, err
	}

	return dbMaster, nil
}

// OpenGORM creates a new *gorm.DB Database instance
func (dc *databaseConfig) OpenGORM() (*gorm.DB, error) {
	logFields := logrus.Fields{
		"platform": "postgres",
		"domain":   dc.domain,
	}
	logrus.WithFields(logFields).Info("Connecting to Postgres SQL DB [GORM]")

	config := &gorm.Config{}
	if dc.config.DEBUG {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	database, err := gorm.Open(postgres.Open(dc.connectionString()), config)

	if err != nil {
		logrus.WithFields(logFields).Fatal(err)
		return nil, err
	}
	return database, err
}
