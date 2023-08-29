package internal

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string         `env:"PORT"`
	Db   DatabaseConfig `env-prefix:"DB_"`
}

type DatabaseConfig struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
}

func (c *DatabaseConfig) BuildDsn(databaseType string) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s", databaseType, c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		err = cleanenv.ReadEnv(cfg)
	}

	if err != nil {
		return nil, fmt.Errorf("Could not get env vars from .env file nor from environment")
	}

	return cfg, nil
}
