package cfg

import (
	"fmt"
	"github.com/caarlos0/env/v7"
	"time"
)

type RabbitMQ struct {
	DSN   string `env:"RABBITMQ_DSN,notEmpty"`
	Queue string `env:"RABBITMQ_QUEUE,notEmpty"`
}

type Config struct {
	LogLevel string `env:"LOG_LEVEL,notEmpty"`

	HttpAPI struct {
		Addr                  string `env:"ADDR,notEmpty"`
		ServerShutdownTimeout time.Duration
	}

	Postgres struct {
		Host              string `env:"POSTGRES_HOST,notEmpty"`
		Port              string `env:"POSTGRES_PORT,notEmpty"`
		User              string `env:"POSTGRES_USER,notEmpty"`
		Password          string `env:"POSTGRES_PASSWORD,notEmpty"`
		Database          string `env:"POSTGRES_DB,notEmpty"`
		ConnectionTimeout time.Duration
	}

	RabbitMQ RabbitMQ
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s connect_timeout=%d sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.Database, c.Postgres.ConnectionTimeout,
	)
}

func Read() (*Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.HttpAPI.Addr = fmt.Sprintf(":%s", config.HttpAPI.Addr)

	return setTimeouts(&config), nil
}

func setTimeouts(c *Config) *Config {
	c.HttpAPI.ServerShutdownTimeout = 10 * time.Second
	c.Postgres.ConnectionTimeout = 5 * time.Second

	return c
}
