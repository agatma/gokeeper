package config

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Config struct {
	DSN          string        `env:"DATABASE_DSN" json:"database_dsn"`
	Address      string        `env:"ADDRESS" json:"address"`
	LogLevel     string        `env:"LOG_LEVEL" json:"log_level"`
	JWTSecretKey string        `env:"SECRET_KEY" json:"secret_key"`
	TokenExp     time.Duration `env:"TOKEN_EXP" json:"token_exp"`
}

func NewConfig() *Config {
	defaultSecretKey := make([]byte, 16)
	rand.Read(defaultSecretKey)

	cfg := &Config{
		DSN:          "postgres://postgres:postgres@localhost:6666/gokeeper",
		LogLevel:     "info",
		Address:      ":8080",
		JWTSecretKey: hex.EncodeToString(defaultSecretKey),
		TokenExp:     time.Hour * 24,
	}
	return cfg
}
