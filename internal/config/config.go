package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	AppEnv              string
	HTTPPort            string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	ShutdownTimeout     time.Duration
	DatabaseURL         string
	DBMaxConns          int32
	DBMinConns          int32
	DBMaxConnLifetime   time.Duration
}

// Load reads configuration from the environment and applies sensible defaults.
func Load() (*Config, error) {
	maxConns, err := getInt32Env("DB_MAX_CONNS", 10)
	if err != nil {
		return nil, err
	}

	minConns, err := getInt32Env("DB_MIN_CONNS", 2)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		ReadTimeout:       getDurationEnv("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:       getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:   getDurationEnv("HTTP_SHUTDOWN_TIMEOUT", 15*time.Second),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todos?sslmode=disable"),
		DBMaxConns:        maxConns,
		DBMinConns:        minConns,
		DBMaxConnLifetime: getDurationEnv("DB_MAX_CONN_LIFETIME", time.Hour),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.HTTPPort == "" {
		return fmt.Errorf("HTTP_PORT must not be empty")
	}

	port, err := strconv.Atoi(c.HTTPPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("HTTP_PORT must be a valid port number (1-65535), got %q", c.HTTPPort)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL must not be empty")
	}

	if c.DBMaxConns < 1 {
		return fmt.Errorf("DB_MAX_CONNS must be at least 1")
	}

	if c.DBMinConns < 0 {
		return fmt.Errorf("DB_MIN_CONNS must be zero or greater")
	}

	if c.DBMinConns > c.DBMaxConns {
		return fmt.Errorf("DB_MIN_CONNS must not exceed DB_MAX_CONNS")
	}

	return nil
}

// Addr returns the HTTP listen address.
func (c *Config) Addr() string {
	return ":" + c.HTTPPort
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return d
}

func getInt32Env(key string, fallback int32) (int32, error) {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer, got %q", key, value)
	}

	return int32(parsed), nil
}
