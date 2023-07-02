package config

import (
	"errors"
	"strings"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Port int `koanf:"port"`

	DatabaseConfig DatabaseConfig `koanf:"database"`
}

type DatabaseConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
	SSLMode  string `koanf:"sslmode"`
}

var (
	k = koanf.New(".")
)

func LoadConfig() (*Config, error) {
	// load default values using the confmap provider
	k.Load(confmap.Provider(map[string]any{
		"port": 8080,
	}, "."), nil)

	// load configured values from environment variables
	k.Load(env.Provider("BOOKKEEPER_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "BOOKKEEPER_")), "_", ".")
	}), nil)

	var config Config
	err := k.Unmarshal("", &config)
	if err != nil {
		return nil, err
	}

	err = validate(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func validate(config *Config) error {
	if config.DatabaseConfig.Host == "" {
		return errors.New("database host missing")
	}
	if config.DatabaseConfig.Port == 0 {
		return errors.New("database port missing")
	}
	if config.DatabaseConfig.User == "" {
		return errors.New("database username missing")
	}
	if config.DatabaseConfig.Password == "" {
		return errors.New("database password missing")
	}
	if config.DatabaseConfig.Name == "" {
		return errors.New("database name missing")
	}

	return nil
}
