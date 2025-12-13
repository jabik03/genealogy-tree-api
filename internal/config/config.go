package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	ApiConf  HttpConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type HttpConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{
		DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
		},
		HttpConfig{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func (db *DatabaseConfig) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name)
}
