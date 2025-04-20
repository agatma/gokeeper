package config

import (
	"time"
)

type Config struct {
	Addr             string        `env:"CLI_ADDRESS"`
	JWTPath          string        `env:"CLI_JWT_PATH"`
	PrivateDataPath  string        `env:"CLI_DATA_PATH"`
	ServerTimeout    time.Duration `env:"CLI_SERVER_TIMEOUT"`
	ServerRetries    int           `env:"CLI_SERVER_RETRIES"`
	SenderWorkersNum int           `env:"CLI_SENDER_WORKERS_NUM"`
}

func NewConfig() *Config {
	cfg := &Config{
		JWTPath:          "/tmp/gophkeeper.jwt",
		PrivateDataPath:  "./data.json",
		Addr:             "localhost:8080",
		ServerTimeout:    time.Second * 2,
		ServerRetries:    3,
		SenderWorkersNum: 10,
	}

	return cfg
}
