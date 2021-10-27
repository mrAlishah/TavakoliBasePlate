package main

import (
	"GolangTraining/internal/logger"
	"GolangTraining/internal/subscription"
	"GolangTraining/platform/mysql"
	"GolangTraining/platform/redis"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

type ServerConfig struct {
	Enabled        bool             `yaml:"ENABLED"`
	User           string           `yaml:"USER"`
	Pass           string           `yaml:"PASS"`
	AuthEnabled    bool             `yaml:"AUTH_ENABLED"`
	IngestNode     bool             `yaml:"INGEST_NODE"`
	TestThirdParty bool             `yaml:"TEST_THIRD_PARTY"`
	Port          int32 `yaml:"PORT"`
}

type MainConfig struct {
	Logger  logger.Config       `yaml:"LOGGER"`
	Service subscription.Config `yaml:"SERVICE"`
	MySQL   mysql.Config        `yaml:"MYSQL"`
	Redis   redis.Config        `yaml:"REDIS"`
	Server  ServerConfig        `yaml:"SERVER"`
}

// LoadConfig loads configs form provided yaml file or overrides it with env variables
func LoadConfig(filePath string) (*MainConfig, error) {
	cfg := MainConfig{}
	if filePath != "" {
		err := readFile(&cfg, filePath)
		if err != nil {
			return nil, err
		}
	}
	err := readEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFile(cfg *MainConfig, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func readEnv(cfg *MainConfig) error {
	return envconfig.Process("", cfg)
}
