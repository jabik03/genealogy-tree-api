package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	ApiConf  HttpConfig
	JWT      JWTConfig
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

type JWTConfig struct {
	SecretKey string
}

func NewConfig() *Config {
	secret := getEnv("JWT_SECRET_KEY", "")
	if secret == "" {
		panic("JWT_SECRET_KEY is required. Generate one with: openssl rand -base64 32")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
		},
		ApiConf: HttpConfig{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnv("HTTP_PORT", "8080"),
		},
		JWT: JWTConfig{
			SecretKey: secret,
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
