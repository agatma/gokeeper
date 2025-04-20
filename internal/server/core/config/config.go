package config

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Config struct {
	DSN          string        `env:"DATABASE_DSN"`
	Address      string        `env:"ADDRESS"`
	LogLevel     string        `env:"LOG_LEVEL"`
	JWTSecretKey string        `env:"SECRET_KEY"`
	TokenExp     time.Duration `env:"TOKEN_EXP"`
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
