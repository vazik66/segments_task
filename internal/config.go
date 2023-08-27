package internal

import "fmt"

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
