package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	LogLevel string `yaml:"log_level"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
